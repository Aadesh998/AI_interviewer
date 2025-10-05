package handlers

import (
	"ai_interview/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func AudioHandler(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "NO_FILE_PROVIDED",
			"message": "No file provided or failed to read uploaded file",
		})
		log.Printf("[ERROR] Failed to get form file: %v", err)
		return
	}

	tempfile, err := os.CreateTemp("", "upload-*"+filepath.Ext(file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "TEMP_FILE_ERROR",
			"message": "Failed to create temp file",
		})
		log.Printf("[ERROR] Failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tempfile.Name())
	defer tempfile.Close()

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "FILE_OPEN_ERROR",
			"message": "Failed to open uploaded file",
		})
		log.Printf("[ERROR] Failed to open uploaded file: %v", err)
		return
	}
	defer src.Close()

	_, err = io.Copy(tempfile, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "FILE_SAVE_ERROR",
			"message": "Failed to save uploaded file",
		})
		log.Printf("[ERROR] Failed to copy file contents: %v", err)
		return
	}

	log.Printf("[INFO] File uploaded: %s", file.Filename)
	log.Printf("[INFO] File saved temporarily at: %s", tempfile.Name())

	text, err := extractTextFromAudio(tempfile.Name())
	fmt.Println(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "TRANSCRIPTION_ERROR",
			"message": "Audio transcription failed: " + err.Error(),
		})
		log.Printf("[ERROR] Transcription failed: %v", err)
		return
	}

	log.Println("[SUCCESS] Transcription completed successfully")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transcription successful",
		"text":    text,
	})
}

func extractTextFromAudio(audiofile string) (string, error) {
	log.Printf("[INFO] Reading audio file: %s", audiofile)

	audio, err := os.ReadFile(audiofile)
	if err != nil {
		log.Printf("[ERROR] Cannot read audio file: %v", err)
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.deepgram.com/v1/listen", bytes.NewReader(audio))
	if err != nil {
		log.Printf("[ERROR] Cannot create Deepgram request: %v", err)
		return "", err
	}

	token := "DEEPGRAM_API_KEY"
	if token == "" {
		log.Println("[ERROR] DEEPGRAM_TOKEN is not set")
		return "", err
	}

	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/octet-stream")

	log.Println("[INFO] Sending request to Deepgram API")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] Request to Deepgram failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] Deepgram API returned status %d: %s", resp.StatusCode, string(body))
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read Deepgram response: %v", err)
		return "", err
	}

	var dg models.DeepgramResponse
	if err := json.Unmarshal(body, &dg); err != nil {
		log.Printf("[ERROR] Failed to parse Deepgram JSON: %v", err)
		return "", err
	}

	log.Println("[INFO] Transcript received from Deepgram API")
	return dg.Results.Channels[0].Alternatives[0].Transcript, nil
}
