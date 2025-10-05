package handlers

import (
	airesponse "ai_interview/AiResponse"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type InterviewState struct {
	Greeting             string
	Questions            *airesponse.InterviewQuestions
	CurrentQuestionIndex int
	ConversationHistory  []string
}

type SessionWebSocketMap struct {
	SessionSocket *websocket.Conn
	State         *InterviewState
}

var (
	sessionMap  = make(map[string]*SessionWebSocketMap)
	ConvHistory = make(map[string][]string)
	mu          sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func AddSession(session string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	if s, ok := sessionMap[session]; ok {
		s.SessionSocket = conn
	} else {
		sessionMap[session] = &SessionWebSocketMap{
			SessionSocket: conn,
			State:         &InterviewState{},
		}
	}
}

func RemoveSession(session string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessionMap, session)
}

func GetConversationHistory(sessionID string) ([]string, bool) {
	mu.Lock()
	defer mu.Unlock()
	sessionState, ok := sessionMap[sessionID]
	if !ok {
		return nil, false
	}
	historyCopy := make([]string, len(sessionState.State.ConversationHistory))
	copy(historyCopy, sessionState.State.ConversationHistory)
	return historyCopy, true
}

func InitializeInterviewState(sessionID string, greeting string, questions *airesponse.InterviewQuestions) {
	mu.Lock()
	defer mu.Unlock()
	if s, ok := sessionMap[sessionID]; ok {
		s.State.Greeting = greeting
		s.State.Questions = questions
	} else {
		sessionMap[sessionID] = &SessionWebSocketMap{
			State: &InterviewState{
				Greeting:  greeting,
				Questions: questions,
			},
		}
	}
}

func WSHandler(c *gin.Context) {
	session := c.Query("session_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("failed to make Connection: %s", err)
		return
	}
	AddSession(session, conn)
	defer RemoveSession(session)

	mu.Lock()
	s, ok := sessionMap[session]
	if !ok {
		mu.Unlock()
		return
	}
	mu.Unlock()

	s.SessionSocket.WriteMessage(websocket.TextMessage, []byte(s.State.Greeting))
	s.State.ConversationHistory = append(s.State.ConversationHistory, s.State.Greeting)

	firstQuestion := s.State.Questions.Questions[0]
	s.SessionSocket.WriteMessage(websocket.TextMessage, []byte(firstQuestion))
	s.State.ConversationHistory = append(s.State.ConversationHistory, firstQuestion)
	s.State.CurrentQuestionIndex = 0

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message from Socket: %s", err)
			break
		}

		s.State.ConversationHistory = append(s.State.ConversationHistory, string(msg))

		evaluation := airesponse.EvaluateAnswer(s.State.ConversationHistory)
		s.State.ConversationHistory = append(s.State.ConversationHistory, evaluation)
		s.SessionSocket.WriteMessage(websocket.TextMessage, []byte(evaluation))

		s.State.CurrentQuestionIndex++
		if s.State.CurrentQuestionIndex < len(s.State.Questions.Questions) {
			nextQuestion := s.State.Questions.Questions[s.State.CurrentQuestionIndex]
			s.State.ConversationHistory = append(s.State.ConversationHistory, nextQuestion)
			s.SessionSocket.WriteMessage(websocket.TextMessage, []byte(nextQuestion))
		} else {
			s.SessionSocket.WriteMessage(websocket.TextMessage, []byte("Thank you for your time. The interview is now complete. Download Feedback Report"))
			conversationHistory, ok := GetConversationHistory(session)
			if !ok {
				log.Printf("Session ID or History is Empty.")
			}
			ConvHistory[session] = conversationHistory
			// _, err := downlaodreport.CreateFeedbackReportFile(s.State.ConversationHistory)
			// if err != nil {
			// 	log.Printf("Failed to make Report: %s", err)
			// }
			break
		}
	}
}
