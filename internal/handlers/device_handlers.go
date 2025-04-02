package handlers

import (
	"context"
	"navic-backend/internal/database"
	message_queue "navic-backend/internal/message-queue"
	"navic-backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddDevice(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device.ID = primitive.NewObjectID()
	device.LastUpdated = time.Now()
	device.History = []models.HistoryEntry{} // Initialize history as an empty array

	collection := database.MongoDB.Collection("devices")
	_, err := collection.InsertOne(context.Background(), device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = message_queue.CreateQueue(device.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func UpdateDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device.Location = models.Location{}
	device.LastUpdated = time.Now()

	collection := database.MongoDB.Collection("devices")
	filter := bson.M{"deviceId": deviceId}
	update := bson.M{
		"$set": bson.M{
			"name":         device.Name,
			"type":         device.Type,
			"status":       device.Status,
			"batteryLevel": device.BatteryLevel,
			"lastUpdated":  device.LastUpdated,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully"})
}

func DeleteDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")

	collection := database.MongoDB.Collection("devices")
	filter := bson.M{"deviceId": deviceId}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = message_queue.DeleteQueue(deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func GetDevices(c *gin.Context) {
	collection := database.MongoDB.Collection("devices")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var devices []models.Device
	if err := cursor.All(context.Background(), &devices); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, devices)
}

func GetDeviceHistory(c *gin.Context) {
	deviceId := c.Param("deviceId")

	collection := database.MongoDB.Collection("devices")
	filter := bson.M{"deviceId": deviceId}

	var device models.Device
	err := collection.FindOne(context.Background(), filter).Decode(&device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device.History)
}

func GetDeviceByID(c *gin.Context) {
	deviceId := c.Param("deviceId")

	collection := database.MongoDB.Collection("devices")
	filter := bson.M{"deviceId": deviceId}

	var device models.Device
	err := collection.FindOne(context.Background(), filter).Decode(&device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}
