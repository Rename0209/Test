package handler

import (
	"context"
	"net/http"
	"test/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Định nghĩa cấu trúc dữ liệu nhận từ webhook
type WebhookPayload struct {
	NotificationToken string `json:"notification_messages_token"`
	RecipientID       string `json:"recipient_id"`
	Reoptin           string `json:"notification_messages_reoptin"`
	TopicTitle        string `json:"topic_title"`
	CreationTime      int64  `json:"creation_timestamp"`
	TokenExpiry       int64  `json:"token_expiry_timestamp"`
	TokenStatus       string `json:"user_token_status"`
	TimeZone          string `json:"notification_messages_timezone"`
	NextEligibleTime  int64  `json:"next_eligible_time"`
}

// Hàm xử lý webhook
func WebhookHandler(c *gin.Context) {
	var payload WebhookPayload

	// Parse JSON từ request
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Lưu vào MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.DB.InsertOne(ctx, bson.M{
		"notification_messages_token":    payload.NotificationToken,
		"recipient_id":                   payload.RecipientID,
		"notification_messages_reoptin":  payload.Reoptin,
		"topic_title":                    payload.TopicTitle,
		"creation_timestamp":             payload.CreationTime,
		"token_expiry_timestamp":         payload.TokenExpiry,
		"user_token_status":              payload.TokenStatus,
		"next_eligible_time":             payload.NextEligibleTime,
		"notification_messages_timezone": payload.TimeZone,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Received and stored successfully"})
}

// Lấy danh sách dữ liệu từ MongoDB
func GetDataHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"notification_messages_token":   1,
		"recipient_id":                  1,
		"notification_messages_reoptin": 1,
		"topic_title":                   1,
		"creation_timestamp":            1,
		"token_expiry_timestamp":        1,
		"next_eligible_time":            1,
		"user_token_status":             1,
		"_id":                           0, // Ẩn ObjectID
	}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch data"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing data"})
		return
	}

	c.JSON(http.StatusOK, results)
}
