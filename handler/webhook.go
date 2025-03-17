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
)

// Định nghĩa cấu trúc dữ liệu nhận từ webhook
type WebhookPayload struct {
	Notification_Token string `bson:"notification_messages_token" json:"notification_messages_token"`
	Recipient_ID       string `bson:"recipient_id" json:"recipient_id"`
	Reoptin            string `bson:"notification_messages_reoptin" json:"notification_messages_reoptin"`
	TopicTitle         string `bson:"topic_title" json:"topic_title"`
	CreationTime       int64  `bson:"creation_timestamp" json:"creation_timestamp"`
	TokenExpiry        int64  `bson:"token_expiry_timestamp" json:"token_expiry_timestamp"`
	TokenStatus        string `bson:"user_token_status" json:"user_token_status"`
	TimeZone           string `bson:"notification_messages_timezone" json:"notification_messages_timezone"`
	NextEligibleTime   int64  `bson:"next_eligible_time" json:"next_eligible_time"`
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

func GetDataHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Find(ctx, bson.M{}) // Không cần projection để kiểm tra dữ liệu
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

	// Nếu dữ liệu có vấn đề, kiểm tra projection hoặc database schema
	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No data found"})
		return
	}

	c.JSON(http.StatusOK, results)
}
