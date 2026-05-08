package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/digitalmaxing/devices-api/internal/models"
)

// mockDeviceService is a mock for the service layer
type mockDeviceService struct {
	mock.Mock
}

func (m *mockDeviceService) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	args := m.Called(ctx, device)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *mockDeviceService) GetDevice(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *mockDeviceService) ListDevices(ctx context.Context, brand, state string) ([]models.Device, error) {
	args := m.Called(ctx, brand, state)
	return args.Get(0).([]models.Device), args.Error(1)
}

func (m *mockDeviceService) UpdateDevice(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Device, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *mockDeviceService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *mockDeviceService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockSvc := new(mockDeviceService)
	handler := NewDeviceHandler(mockSvc)

	api := router.Group("/devices")
	{
		api.POST("", handler.CreateDevice)
		api.GET("", handler.ListDevices)
		api.GET("/:id", handler.GetDevice)
		api.PATCH("/:id", handler.UpdateDevice)
		api.DELETE("/:id", handler.DeleteDevice)
	}
	return router, mockSvc
}

func TestCreateDevice(t *testing.T) {
	router, mockSvc := setupTestRouter()

	device := &models.Device{
		ID:    uuid.New(),
		Name:  "Test Device",
		Brand: "Test Brand",
		State: models.StateAvailable,
	}

	mockSvc.On("CreateDevice", mock.Anything, mock.Anything).Return(device, nil)

	body, _ := json.Marshal(map[string]string{
		"name":  "Test Device",
		"brand": "Test Brand",
		"state": "available",
	})
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetDevice(t *testing.T) {
	router, mockSvc := setupTestRouter()
	id := uuid.New()

	device := &models.Device{ID: id, Name: "Test", Brand: "Test", State: models.StateAvailable}
	mockSvc.On("GetDevice", mock.Anything, id).Return(device, nil)

	req, _ := http.NewRequest("GET", "/devices/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestListDevices(t *testing.T) {
	router, mockSvc := setupTestRouter()

	devices := []models.Device{{Name: "Test", Brand: "Test", State: models.StateAvailable}}
	mockSvc.On("ListDevices", mock.Anything, "", "").Return(devices, nil)

	req, _ := http.NewRequest("GET", "/devices", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateDevice(t *testing.T) {
	router, mockSvc := setupTestRouter()
	id := uuid.New()

	updated := &models.Device{ID: id, Name: "Updated", Brand: "Test", State: models.StateAvailable}
	mockSvc.On("UpdateDevice", mock.Anything, id, mock.Anything).Return(updated, nil)

	body, _ := json.Marshal(map[string]string{"name": "Updated"})
	req, _ := http.NewRequest("PATCH", "/devices/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteDevice(t *testing.T) {
	router, mockSvc := setupTestRouter()
	id := uuid.New()

	mockSvc.On("DeleteDevice", mock.Anything, id).Return(nil)

	req, _ := http.NewRequest("DELETE", "/devices/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

// === Additional Sanity + Coverage Tests ===

func TestCreateDevice_InvalidJSON(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDevice_InvalidUUID(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/devices/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDevice_ValidationError(t *testing.T) {
	router, mockSvc := setupTestRouter()
	id := uuid.New()

	mockSvc.On("UpdateDevice", mock.Anything, id, mock.Anything).Return(nil, errors.New("name cannot be updated for in-use device"))

	body, _ := json.Marshal(map[string]string{"name": "New Name"})
	req, _ := http.NewRequest("PATCH", "/devices/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}