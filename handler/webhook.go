package handler

import (
	"context"
	"fmt"
	"log"
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
	var payload []WebhookPayload

	// Parse dữ liệu JSON
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Lưu toàn bộ dữ liệu vào MongoDB
	insertedCount := 0
	for _, data := range payload {
		_, err := database.DB.InsertOne(context.TODO(), data)
		if err != nil {
			log.Println("❌ Lỗi khi lưu dữ liệu:", err)
			continue
		}
		insertedCount++
	}

	// Trả về kết quả
	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"inserted_count": insertedCount,
		"message":        fmt.Sprintf("Đã nhận và lưu %d dữ liệu", insertedCount),
	})
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
