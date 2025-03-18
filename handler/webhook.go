package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	// Đọc dữ liệu raw từ request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu"})
		return
	}

	// Kiểm tra dữ liệu JSON
	if json.Unmarshal(body, &payload) != nil {
		// Nếu không phải mảng, thử parse object đơn lẻ
		var singlePayload WebhookPayload
		if json.Unmarshal(body, &singlePayload) != nil {
			log.Println("❌ Dữ liệu không hợp lệ:", string(body))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}
		payload = append(payload, singlePayload)
	}

	// Nếu payload rỗng
	if len(payload) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload không có dữ liệu"})
		return
	}

	// Chuyển đổi dữ liệu sang []interface{} để InsertMany
	var insertData []interface{}
	for _, data := range payload {
		insertData = append(insertData, data)
	}

	// Lưu toàn bộ dữ liệu vào MongoDB với InsertMany
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := database.DB.InsertMany(ctx, insertData)
	if err != nil {
		log.Println("❌ Lỗi khi lưu dữ liệu:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu dữ liệu"})
		return
	}

	// Trả về kết quả
	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"inserted_count": len(res.InsertedIDs),
		"message":        fmt.Sprintf("Đã nhận và lưu %d dữ liệu", len(res.InsertedIDs)),
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

// Hàm tìm kiếm dữ liệu theo title
func FindByTitle(c *gin.Context) {
	topicTitle := c.Query("topic_title")
	if topicTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic_title là bắt buộc"})
		return
	}

	findDataByField(c, "topic_title", topicTitle)
}

// Hàm tìm kiếm dữ liệu theo Recipient_ID
func FindByRecipientID(c *gin.Context) {
	recipientID := c.Query("recipient_id")
	if recipientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipient_id là bắt buộc"})
		return
	}

	findDataByField(c, "recipient_id", recipientID)
}

// Hàm chung tìm kiếm theo một trường cụ thể
func findDataByField(c *gin.Context, fieldName, fieldValue string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{fieldName: fieldValue}
	//projection := options.Find().SetProjection(bson.M{"recipient_id": 1})

	cursor, err := database.DB.Find(ctx, filter /*, projection*/)
	if err != nil {
		log.Println("❌ Lỗi khi tìm dữ liệu:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn dữ liệu"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi parse dữ liệu"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Không tìm thấy dữ liệu"})
		return
	}

	c.JSON(http.StatusOK, results)
}
