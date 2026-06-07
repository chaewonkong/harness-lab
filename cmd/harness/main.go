package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const OLLAMA_HOST = "OLLAMA_HOST"

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools"`
	Stream   bool      `json:"stream"`
	Think    bool      `json:"think"`
}

type Message struct {
	Role    string `json:"role"` // user
	Content string `json:"content"`
}
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}
type Function struct {
	Name       string     `json:"name"`
	Descripion string     `json:"description"`
	Parameters Parameters `json:"parameters"`
	Required   []string   `json:"required"`
}

type Parameters struct {
	Type       string              `json:"type"` // object
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func main() {
	host := os.Getenv(OLLAMA_HOST)
	if host == "" {
		log.Fatal("OLLAMA_HOST not found in environ")
	}

	apiURL := fmt.Sprintf("http://%s/api/chat", host)

	r := Request{
		Model:    "phi4-mini",
		Messages: []Message{{Role: "user", Content: "지금 영국 런던의 시간은 몇시 몇분이야?"}},
		Tools: []Tool{
			{Type: "function", Function: Function{
				Name:       "get_current_time",
				Descripion: "지정한 도시의 현재시각을 반환한다",
				Parameters: Parameters{
					Type: "object",
					Properties: map[string]Property{
						"city": {Type: "string", Description: "도시이름, 예) Seoul"},
					},
				},
				Required: []string{"city"},
			}},
		},
		Think: false,
	}

	buf, err := json.Marshal(r)
	if err != nil {
		slog.Error("failed to marshal data", "err", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		slog.Error("failed post data", "err", err)
		os.Exit(1)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed read response", "err", err)
		os.Exit(1)
	}

	fmt.Println(string(bytes))
}
