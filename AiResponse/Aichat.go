package airesponse

import (
	"ai_interview/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var client = openai.NewClient(
	option.WithAPIKey(""), // -> PUT YOUR OPENAI API KEY
)

func ResumeGreeter(resumeText string) string {
	prompt := fmt.Sprintf(`
		You are a professional senior-level job interviewer.
		You have access to the candidate's resume below.
		Greet the candidate in a friendly and professional way based on their background.
		Keep the greeting concise — no more than 4 sentences.

		Candidate Resume:
		%s
	`, resumeText)

	return CallOpenAIText(prompt)
}

type InterviewQuestions struct {
	Questions []string `json:"questions"`
}

func GenerateInterviewQuestions(domain string) (*InterviewQuestions, error) {

	prompt := fmt.Sprintf(`
		You are an experienced interviewer specializing in the %s domain.

		Generate JSON output **only**, with one arrays of questions.
		And array should contain exactly 15 short questions with level of difficulty is low, mid, high.
		At the end of the Question not to add difficulty of question.
		The JSON must have this exact structure:

		{
			"questions": ["q1", "q2", "q3", "q4", "q5", "q6", "q7", "q8", "q9", "q10", "q11", "q12", "q13", "q14", "q15"]
		}
	`, domain)

	response := CallOpenAIText(prompt)
	log.Printf("Response: %s", response)
	response = utils.CleanResp(response)

	var questions InterviewQuestions
	err := json.Unmarshal([]byte(response), &questions)
	if err != nil {
		log.Printf("JSON unmarshal error: %v\nResponse:\n%s", err, response)
		return nil, err
	}

	return &questions, nil
}

func EvaluateAnswer(conversationHistory []string) string {
	history := strings.Join(conversationHistory, "\n")

	prompt := fmt.Sprintf(`
		You are a professional senior-level job interviewer.
		Below is the conversation history of an interview. The last message is the candidate's answer to your question.
		Please evaluate the candidate's answer in a friendly and professional way and say "Let's move to the next question".
		Keep the evaluation concise — no more than 2 sentences.
		Do not include any preamble or explanation.

		Conversation History:
		%s
		`, history)

	return CallOpenAIText(prompt)
}

func CallOpenAIText(prompt string) string {
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
		Temperature: openai.Float(0.9),
	})
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}
