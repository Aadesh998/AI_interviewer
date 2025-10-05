package handlers

import (
	airesponse "ai_interview/AiResponse"
	"ai_interview/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateSession(c *gin.Context) {

	file, err := c.FormFile("resume")
	if err != nil {
		log.Printf("ERROR: Failed to get file from form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to get uploaded file",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		log.Printf("WARN: Unsupported file type uploaded: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Unsupported file type. Please upload a PDF.",
		})
		return
	}

	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		log.Printf("ERROR: Failed to create temp file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create temp file",
		})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to open uploaded file",
		})
		return
	}
	defer src.Close()

	if _, err := io.Copy(tempFile, src); err != nil {
		log.Printf("ERROR: Failed to copy uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to process the. file",
		})
		return
	}

	extractedText, err := utils.ExtractTextFromPDF(tempFile.Name())
	if err != nil {
		log.Printf("ERROR: Failed to extract text: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to extract text from PDF",
		})
		return
	}

	log.Printf("INFO: Extracted %d characters from %s", len(extractedText), file.Filename)

	session_id := utils.GetUUID()

	greeting := airesponse.ResumeGreeter(extractedText)

	domain := "Excel Interview"
	questions, err := airesponse.GenerateInterviewQuestions(domain)
	if err != nil {
		log.Printf("ERROR: Failed to generate interview questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate interview questions",
		})
		return
	}

	InitializeInterviewState(session_id, greeting, questions)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"session_id": session_id,
		"message":    "session created Successfully for interview.",
	})
}
