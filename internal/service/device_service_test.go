package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/digitalmaxing/devices-api/internal/models"
)

// mockDeviceRepository is a testify mock for the repository interface.
type mockDeviceRepository struct {
	mock.Mock
}

func (m *mockDeviceRepository) Create(ctx context.Context, device *models.Device) error {
	args := m.Called(ctx, device)
	return args.Error(0)
}

func (m *mockDeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *mockDeviceRepository) List(ctx context.Context, brand, state string) ([]models.Device, error) {
	args := m.Called(ctx, brand, state)
	return args.Get(0).([]models.Device), args.Error(1)
}

func (m *mockDeviceRepository) Update(ctx context.Context, device *models.Device) error {
	args := m.Called(ctx, device)
	return args.Error(0)
}

func (m *mockDeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestDeviceService_CreateDevice(t *testing.T) {
	mockRepo := new(mockDeviceRepository)
	svc := NewDeviceService(mockRepo)

	device := &models.Device{Name: "iPhone 15", Brand: "Apple"}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Device")).Return(nil)

	created, err := svc.CreateDevice(context.Background(), device)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, models.StateAvailable, created.State)
	mockRepo.AssertExpectations(t)
}

func TestDeviceService_DeleteDevice_InUse(t *testing.T) {
	mockRepo := new(mockDeviceRepository)
	svc := NewDeviceService(mockRepo)

	id := uuid.New()
	inUseDevice := &models.Device{ID: id, State: models.StateInUse}

	mockRepo.On("GetByID", mock.Anything, id).Return(inUseDevice, nil)

	err := svc.DeleteDevice(context.Background(), id)

	assert.Error(t, err)
	assert.Equal(t, "in-use devices cannot be deleted", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeviceService_UpdateDevice_NameProtectedWhenInUse(t *testing.T) {
	mockRepo := new(mockDeviceRepository)
	svc := NewDeviceService(mockRepo)

	id := uuid.New()
	inUseDevice := &models.Device{ID: id, Name: "Old", Brand: "OldBrand", State: models.StateInUse}

	mockRepo.On("GetByID", mock.Anything, id).Return(inUseDevice, nil)

	updates := map[string]interface{}{"name": "New Name"}

	_, err := svc.UpdateDevice(context.Background(), id, updates)

	assert.Error(t, err)
	assert.Equal(t, "name cannot be updated for in-use device", err.Error())
	mockRepo.AssertNotCalled(t, "Update") // Should not reach update
}