package models

import (
	"time"

	"github.com/google/uuid"
)

// DeviceState represents the possible states of a device.
type DeviceState string

const (
	StateAvailable DeviceState = "available"
	StateInUse     DeviceState = "in-use"
	StateInactive  DeviceState = "inactive"
)

// Device represents the core domain entity for a device resource.
type Device struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string      `gorm:"size:255;not null" json:"name" binding:"required,min=1,max=255"`
	Brand     string      `gorm:"size:255;not null" json:"brand" binding:"required,min=1,max=255"`
	State     DeviceState `gorm:"size:20;not null;default:'available'" json:"state" binding:"required,oneof=available in-use inactive"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
}

// IsInUse returns true if the device is currently in the "in-use" state.
func (d *Device) IsInUse() bool {
	return d.State == StateInUse
}

// ValidateStateTransition checks if a state change is valid (basic example, can be extended).
func (d *Device) ValidateStateTransition(newState DeviceState) bool {
	// Example: from inactive can go to available, etc. For simplicity, allow most transitions.
	// In real app, more business rules could apply here.
	return true
}