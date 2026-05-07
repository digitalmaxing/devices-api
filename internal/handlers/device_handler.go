package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/digitalmaxing/devices-api/internal/models"
	"github.com/digitalmaxing/devices-api/internal/service"
)

// DeviceHandler handles HTTP requests for device operations.
type DeviceHandler struct {
	service *service.DeviceService
}

// NewDeviceHandler creates a new handler with injected service.
func NewDeviceHandler(svc *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: svc}
}

// respondError is a helper to return consistent JSON error responses.
func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// CreateDevice handles POST /devices
// @Summary Create a new device
// @Description Creates a new device resource. State defaults to 'available' if not provided.
// @Tags devices
// @Accept json
// @Produce json
// @Param device body models.Device true "Device creation payload (name, brand required; state optional)"
// @Success 201 {object} models.Device
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	created, err := h.service.CreateDevice(c.Request.Context(), &device)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create device")
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetDevice handles GET /devices/:id
// @Summary Get a device by ID
// @Description Retrieves a single device by its UUID
// @Tags devices
// @Produce json
// @Param id path string true "Device UUID"
// @Success 200 {object} models.Device
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Device not found"
// @Router /devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid device ID format")
		return
	}

	device, err := h.service.GetDevice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("device not found")) || err.Error() == "device not found" {
			respondError(c, http.StatusNotFound, "device not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to retrieve device")
		return
	}

	c.JSON(http.StatusOK, device)
}

// ListDevices handles GET /devices with optional filters
// @Summary List all devices (with optional filters)
// @Description Returns list of devices. Supports filtering by brand and/or state via query params.
// @Tags devices
// @Produce json
// @Param brand query string false "Filter by brand"
// @Param state query string false "Filter by state (available, in-use, inactive)"
// @Success 200 {array} models.Device
// @Failure 500 {object} map[string]string
// @Router /devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	brand := c.Query("brand")
	state := c.Query("state")

	devices, err := h.service.ListDevices(c.Request.Context(), brand, state)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to list devices")
		return
	}

	c.JSON(http.StatusOK, devices)
}

// UpdateDevice handles PATCH /devices/:id for partial updates
// @Summary Partially update a device
// @Description Updates allowed fields. Enforces: no creation time change, no name/brand change if in-use.
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device UUID"
// @Param updates body map[string]interface{} true "Partial fields to update (e.g. {\"name\": \"New Name\", \"state\": \"inactive\"})"
// @Success 200 {object} models.Device
// @Failure 400 {object} map[string]string "Invalid input or validation error"
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [patch]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid device ID format")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	updated, err := h.service.UpdateDevice(c.Request.Context(), id, updates)
	if err != nil {
		if err.Error() == "device not found" {
			respondError(c, http.StatusNotFound, "device not found")
			return
		}
		// Domain validation errors -> 400 or 409
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteDevice handles DELETE /devices/:id
// @Summary Delete a device
// @Description Deletes a device. Fails for in-use devices per business rules.
// @Tags devices
// @Param id path string true "Device UUID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid ID or business rule violation"
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid device ID format")
		return
	}

	err = h.service.DeleteDevice(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "device not found" {
			respondError(c, http.StatusNotFound, "device not found")
			return
		}
		if err.Error() == "in-use devices cannot be deleted" {
			respondError(c, http.StatusConflict, "in-use devices cannot be deleted")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to delete device")
		return
	}

	c.Status(http.StatusNoContent)
}