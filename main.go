package main

import (
	"ai_interview/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.StaticFile("/", "./index.html")
	router.StaticFile("/dashboard", "./dashboard.html")

	router.POST("/create-session", handlers.CreateSession)
	router.GET("/ws", handlers.WSHandler)
	router.POST("/audio", handlers.AudioHandler)
	router.GET("/download-feedback", handlers.FeedbackHandler)

	router.Run(":8000")
}
