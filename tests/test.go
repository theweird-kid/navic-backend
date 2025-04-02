package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"navic-backend/internal/database"
	"navic-backend/internal/handlers"
	"navic-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func TestMain(m *testing.M) {
	// Setup
	var err error
	client, err = database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Run tests
	code := m.Run()

	// Teardown
	client.Disconnect(context.Background())

	os.Exit(code)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/devices", handlers.AddDevice)
		api.PUT("/devices/:deviceId", handlers.UpdateDevice)
		api.DELETE("/devices/:deviceId", handlers.DeleteDevice)
		api.GET("/devices", handlers.GetDevices)
		api.GET("/devices/:deviceId/history", handlers.GetDeviceHistory)
	}
	return r
}

func setupTestDatabase() {
	collection := client.Database("navicdb").Collection("devices")
	collection.DeleteMany(context.Background(), bson.M{})
}

func TestAddDevice(t *testing.T) {
	setupTestDatabase()
	r := setupRouter()

	device := models.Device{
		Name:         "Test Device",
		DeviceID:     "TEST-001",
		Type:         "Test",
		Status:       "Active",
		LastUpdated:  time.Now(),
		Location:     models.Location{Lat: 0, Lng: 0},
		BatteryLevel: 100,
		History:      []models.HistoryEntry{},
	}

	jsonValue, _ := json.Marshal(device)
	req, _ := http.NewRequest("POST", "/api/devices", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateDevice(t *testing.T) {
	setupTestDatabase()
	r := setupRouter()

	// Insert a test device
	collection := client.Database("navicdb").Collection("devices")
	device := models.Device{
		ID:           primitive.NewObjectID(),
		Name:         "Test Device",
		DeviceID:     "TEST-001",
		Type:         "Test",
		Status:       "Active",
		LastUpdated:  time.Now(),
		Location:     models.Location{Lat: 0, Lng: 0},
		BatteryLevel: 100,
		History:      []models.HistoryEntry{},
	}
	collection.InsertOne(context.Background(), device)

	updatedDevice := models.Device{
		Name:         "Updated Test Device",
		Type:         "Updated Test",
		Status:       "Inactive",
		BatteryLevel: 50,
	}

	jsonValue, _ := json.Marshal(updatedDevice)
	req, _ := http.NewRequest("PUT", "/api/devices/TEST-001", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteDevice(t *testing.T) {
	setupTestDatabase()
	r := setupRouter()

	// Insert a test device
	collection := client.Database("navicdb").Collection("devices")
	device := models.Device{
		ID:           primitive.NewObjectID(),
		Name:         "Test Device",
		DeviceID:     "TEST-001",
		Type:         "Test",
		Status:       "Active",
		LastUpdated:  time.Now(),
		Location:     models.Location{Lat: 0, Lng: 0},
		BatteryLevel: 100,
		History:      []models.HistoryEntry{},
	}
	collection.InsertOne(context.Background(), device)

	req, _ := http.NewRequest("DELETE", "/api/devices/TEST-001", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDevices(t *testing.T) {
	setupTestDatabase()
	r := setupRouter()

	// Insert a test device
	collection := client.Database("navicdb").Collection("devices")
	device := models.Device{
		ID:           primitive.NewObjectID(),
		Name:         "Test Device",
		DeviceID:     "TEST-001",
		Type:         "Test",
		Status:       "Active",
		LastUpdated:  time.Now(),
		Location:     models.Location{Lat: 0, Lng: 0},
		BatteryLevel: 100,
		History:      []models.HistoryEntry{},
	}
	collection.InsertOne(context.Background(), device)

	req, _ := http.NewRequest("GET", "/api/devices", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var devices []models.Device
	json.Unmarshal(w.Body.Bytes(), &devices)
	assert.NotEmpty(t, devices)
}

func TestGetDeviceHistory(t *testing.T) {
	setupTestDatabase()
	r := setupRouter()

	// Insert a test device
	collection := client.Database("navicdb").Collection("devices")
	device := models.Device{
		ID:           primitive.NewObjectID(),
		Name:         "Test Device",
		DeviceID:     "TEST-001",
		Type:         "Test",
		Status:       "Active",
		LastUpdated:  time.Now(),
		Location:     models.Location{Lat: 0, Lng: 0},
		BatteryLevel: 100,
		History: []models.HistoryEntry{
			{Timestamp: time.Now(), Location: models.Location{Lat: 0, Lng: 0}},
		},
	}
	collection.InsertOne(context.Background(), device)

	req, _ := http.NewRequest("GET", "/api/devices/TEST-001/history", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var history []models.HistoryEntry
	json.Unmarshal(w.Body.Bytes(), &history)
	assert.NotEmpty(t, history)
}
