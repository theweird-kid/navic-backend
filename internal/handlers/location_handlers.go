package handlers

import (
	"context"
	"navic-backend/internal/database"
	"navic-backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateDeviceLocation(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.MongoDB.Collection("devices")
	filter := bson.M{"deviceId": deviceId}
	update := bson.M{
		"$set": bson.M{
			"location":    location,
			"lastUpdated": time.Now(),
		},
		"$push": bson.M{
			"history": models.HistoryEntry{
				Timestamp: time.Now(),
				Location:  location,
			},
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully"})
}

func ClearDeviceLocation(c *gin.Context) {
	deviceId := c.Param("deviceId")

	collection := database.MongoDB.Collection("devices")
	filter := bson.M{"deviceId": deviceId}
	update := bson.M{
		"$set": bson.M{
			"history": []models.HistoryEntry{},
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device location history cleared successfully"})
}
