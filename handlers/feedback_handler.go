package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	downloadreport "ai_interview/DownlaodReport"

	"github.com/gin-gonic/gin"
)

func FeedbackHandler(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.String(http.StatusBadRequest, "session_id is required")
		return
	}

	conversationHistory, ok := ConvHistory[sessionID]
	if !ok {
		c.String(http.StatusNotFound, "History is empty")
		return
	}

	filePath, err := downloadreport.CreateFeedbackReportFile(conversationHistory)
	if err != nil {
		log.Printf("Error creating feedback report file: %s", err)
		c.String(http.StatusInternalServerError, "Failed to create feedback report")
		return
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove temporary file %s: %s", filePath, err)
		}
	}()

	fileName := filepath.Base(filePath)
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Writer.Header().Set("Content-Type", "text/plain")

	http.ServeFile(c.Writer, c.Request, filePath)
	log.Printf("Successfully served and initiated download of feedback report for session %s", sessionID)
}
