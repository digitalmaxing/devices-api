package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/digitalmaxing/devices-api/internal/models"
)

// DeviceRepository defines the interface for device data access operations.
// This allows for easy mocking in tests and swapping implementations (e.g., different DBs).
type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error)
	List(ctx context.Context, brand string, state string) ([]models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// postgresDeviceRepository is the PostgreSQL implementation using GORM.
type postgresDeviceRepository struct {
	db *gorm.DB
}

// NewPostgresDeviceRepository creates a new repository instance.
func NewPostgresDeviceRepository(db *gorm.DB) DeviceRepository {
	return &postgresDeviceRepository{db: db}
}

// Create persists a new device to the database.
func (r *postgresDeviceRepository) Create(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// GetByID retrieves a single device by its UUID.
func (r *postgresDeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).First(&device, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("device not found")
	}
	return &device, err
}

// List retrieves devices, optionally filtered by brand and/or state.
// Empty strings mean no filter for that field.
func (r *postgresDeviceRepository) List(ctx context.Context, brand, state string) ([]models.Device, error) {
	var devices []models.Device
	query := r.db.WithContext(ctx).Model(&models.Device{})

	if brand != "" {
		query = query.Where("brand = ?", brand)
	}
	if state != "" {
		query = query.Where("state = ?", state)
	}

	err := query.Find(&devices).Error
	return devices, err
}

// Update persists changes to an existing device (full or selective via caller).
func (r *postgresDeviceRepository) Update(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

// Delete removes a device by ID (caller must ensure business rules like not in-use).
func (r *postgresDeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Device{}, "id = ?", id).Error
}