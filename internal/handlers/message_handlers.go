package handlers

import (
	message_queue "navic-backend/internal/message-queue"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendMessageToDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var message struct {
		Content string `json:"message"`
	}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := message_queue.PublishMessage(deviceId, message.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
