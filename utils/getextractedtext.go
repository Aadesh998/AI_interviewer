package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var textBuilder strings.Builder
	numPages := r.NumPage()

	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			log.Printf("WARN: Page %d is null in PDF: %s", i, path)
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			log.Printf("WARN: Error reading text from page %d: %v", i, err)
			continue
		}
		textBuilder.WriteString(content)
	}

	return textBuilder.String(), nil
}
