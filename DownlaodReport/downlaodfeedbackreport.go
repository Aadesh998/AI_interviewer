package downloadreport

import (
	airesponse "ai_interview/AiResponse"
	"fmt"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func CreateFeedbackReportFile(conversationHistory []string) (string, error) {

	prompt := fmt.Sprintf(`
You are an expert feedback analyst. You are given the complete conversation history between a user and an AI assistant.
The conversation includes approximately 15 questions, along with the user's responses and the AI's evaluations.

Your task is to carefully analyze the entire conversation and generate a **comprehensive feedback report** for the user.

The report must include the following sections:

1. **User Details**
   - Extract the user's name if mentioned in the conversation; if not, use "User".

2. **Conversation Summary**
   - Summarize the purpose and context of the entire conversation.
   - Highlight the main topics or goals discussed.

3. **Performance Overview**
   - Provide an overall assessment of how the user performed or interacted throughout the 15 questions.
   - Comment on clarity, accuracy, engagement, and understanding.

4. **Question-by-Question Feedback**
   For each of the 15 questions:
   - Question Number and Text  
   - User’s Response (summarized if lengthy)  
   - AI’s Evaluation (if given)  
   - Your synthesized feedback combining both the user’s response and AI’s evaluation.

5. **Strengths**
   - Highlight the user’s strongest areas or consistent positive traits.

6. **Areas for Improvement**
   - Point out patterns, mistakes, or topics the user can work on.

7. **Recommendations**
   - Give actionable, constructive suggestions to help the user improve in future conversations or assessments.

8. **Final Evaluation**
   - Provide an overall rating or qualitative summary of performance (e.g., "Excellent", "Good", "Needs Improvement"), and a short closing comment.

### Conversation History:
%s

Output: A well-formatted and human-readable report with the above sections.
`,
		conversationHistory)

	response := airesponse.CallOpenAIText(prompt)

	tempFile, err := os.CreateTemp("", "feedback_report_*.pdf")
	if err != nil {
		log.Printf("Error creating temp PDF file: %s", err)
		return "", err
	}
	defer tempFile.Close()

	pdfg := gofpdf.New("P", "mm", "A4", "")
	pdfg.AddPage()
	pdfg.SetFont("Arial", "", 12)
	pdfg.SetLeftMargin(15)
	pdfg.SetRightMargin(15)

	pdfg.MultiCell(0, 8, response, "", "L", false)

	err = pdfg.OutputFileAndClose(tempFile.Name())
	if err != nil {
		log.Printf("Error writing PDF to file: %s", err)
		return "", err
	}

	log.Printf("Successfully created feedback report at: %s", tempFile.Name())
	return tempFile.Name(), nil
}
