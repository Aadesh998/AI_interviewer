# AI Interview

This project is an AI-powered interview platform that allows candidates to practice their interviewing skills. The platform uses OpenAI's GPT-4o to generate interview questions, evaluate candidate answers, and provide feedback.

## Features

- **Resume-based greeting**: The platform greets the candidate based on their resume.
- **AI-generated questions**: The platform generates a list of 15 interview questions for a given domain.
- **Real-time feedback**: The platform provides real-time feedback on the candidate's answers.
- **Audio transcription**: The platform can transcribe audio files using the Deepgram API.

## Getting Started

To run the project, you will need to have Go and an OpenAI API key.

1. Clone the repository:

```
git clone https://github.com/aadesh-d-kumar/ai_interview.git
```

2. Install the dependencies:

```
go mod tidy
```

3. Set your OpenAI API key as an environment variable:

```
export OPENAI_API_KEY=<your-api-key>
```

4. Run the project:

```
go run main.go
```

The application will be available at `http://localhost:8000`.

## Project Structure

```
.
├── AiResponse
│   └── Aichat.go
├── handlers
│   ├── audiohandler.go
│   ├── create_session.go
│   └── Wshandler.go
├── models
│   └── deepgram.go
├── utils
│   ├── cleanOutput.go
│   ├── getextractedtext.go
│   └── getUUID.go
├── dashboard.html
├── go.mod
├── go.sum
├── index.html
└── main.go
```

- **AiResponse**: This directory contains the code for interacting with the OpenAI API.
- **handlers**: This directory contains the HTTP handlers for the API endpoints.
- **models**: This directory contains the data structures for the application.
- **utils**: This directory contains utility functions.
- **dashboard.html**: This is the main dashboard for the application.
- **index.html**: This is the main entry point for the application.
- **main.go**: This is the main entry point for the application.

## API Endpoints

- `POST /create-session`: Creates a new interview session.
- `GET /ws`: Handles the WebSocket connection for the interview.
- `POST /audio`: Transcribes an audio file.

## Frontend

The frontend is built with HTML, CSS, and JavaScript. It allows the user to upload their resume, start a new interview, and answer questions.

## Backend

The backend is built with Go and the Gin framework. It handles the API requests, interacts with the OpenAI API, and manages the interview state.

## Dependencies

- [Gin](https://github.com/gin-gonic/gin)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [OpenAI Go](https://github.com/openai/openai-go)
- [pdf](https://github.com/ledongthuc/pdf)
- [uuid](https://github.com/google/uuid)
