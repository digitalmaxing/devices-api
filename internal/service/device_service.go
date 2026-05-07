package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/digitalmaxing/devices-api/internal/models"
	"github.com/digitalmaxing/devices-api/internal/repository"
)

// DeviceService provides business logic and enforces domain validations
// on top of the repository layer.
type DeviceService struct {
	repo repository.DeviceRepository
}

// NewDeviceService creates a new service with the given repository.
func NewDeviceService(repo repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

// CreateDevice handles creation with defaults (UUID, initial state).
func (s *DeviceService) CreateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}
	if device.State == "" {
		device.State = models.StateAvailable
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}
	return device, nil
}

// GetDevice retrieves a device by ID.
func (s *DeviceService) GetDevice(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	return s.repo.GetByID(ctx, id)
}

// ListDevices returns devices filtered optionally by brand and state.
func (s *DeviceService) ListDevices(ctx context.Context, brand, state string) ([]models.Device, error) {
	return s.repo.List(ctx, brand, state)
}

// UpdateDevice performs partial update with domain validations:
// - Creation time cannot be updated
// - Name/Brand cannot change if device is in-use
// - State changes are allowed (extendable)
func (s *DeviceService) UpdateDevice(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Device, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Enforce: creation time immutable
	if _, ok := updates["created_at"]; ok || updates["CreatedAt"] != nil {
		return nil, errors.New("creation time cannot be updated")
	}

	// Enforce: name and brand immutable for in-use devices
	if existing.IsInUse() {
		if _, ok := updates["name"]; ok {
			return nil, errors.New("name cannot be updated for in-use device")
		}
		if _, ok := updates["brand"]; ok {
			return nil, errors.New("brand cannot be updated for in-use device")
		}
	}

	// Apply allowed updates
	for key, value := range updates {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				existing.Name = v
			}
		case "brand":
			if v, ok := value.(string); ok && v != "" {
				existing.Brand = v
			}
		case "state":
			if v, ok := value.(string); ok {
				newState := models.DeviceState(v)
				// Could add transition validation here
				existing.State = newState
			}
		}
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}
	return existing, nil
}

// DeleteDevice enforces business rule: in-use devices cannot be deleted.
func (s *DeviceService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existing.IsInUse() {
		return errors.New("in-use devices cannot be deleted")
	}

	return s.repo.Delete(ctx, id)
}