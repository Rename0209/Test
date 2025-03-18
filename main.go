package main

import (
	"fmt"
	"os"
	"test/database"
	"test/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectMongoDB()

	r := gin.Default()
	r.LoadHTMLGlob("static/*.html")
	r.Static("/static", "./static")

	r.GET("/find/title", handler.FindByTitle)
	r.GET("/find/recipient", handler.FindByRecipientID)

	r.POST("/webhook", handler.WebhookHandler)
	r.GET("/data", handler.GetDataHandler)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // Trả về No Content (Không có lỗi)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Mặc định dùng 8080 nếu không có biến môi trường
	}
	fmt.Printf("Server đang chạy tại http://localhost:%s\n", port)
	r.Run(fmt.Sprintf(":%s", port))
}
