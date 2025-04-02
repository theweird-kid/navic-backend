package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the Details of a User for Login and SignUp
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
}

// Location represents a geographical location with latitude and longitude.
type Location struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

// HistoryEntry represents a single entry in the location history.
type HistoryEntry struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Location  Location  `bson:"location" json:"location"`
}

// Device represents a tracking device with its metadata, current status, and location history.
type Device struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	DeviceID     string             `bson:"deviceId" json:"deviceId"`
	Type         string             `bson:"type" json:"type"`
	Status       string             `bson:"status" json:"status"`
	LastUpdated  time.Time          `bson:"lastUpdated" json:"lastUpdated"`
	Location     Location           `bson:"location" json:"location"`
	BatteryLevel int                `bson:"batteryLevel" json:"batteryLevel"`
	History      []HistoryEntry     `bson:"history" json:"history"`
}
