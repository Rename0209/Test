package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Collection

func ConnectMongoDB() {

	// Load file .env
	_ = godotenv.Load(".env")

	// Cấu hình URI MongoDB
	uri := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(uri)

	// Kết nối MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Không thể kết nối MongoDB: %v", err)
	}

	// Kiểm tra kết nối
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	fmt.Println("Kết nối MongoDB thành công!")

	// Chọn database và collection
	DB = client.Database("webhookDB").Collection("events")
}
