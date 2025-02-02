package gemini

import "time"

type messageResponse struct {
	CreatedAt    time.Time `json:"created_at"`
	Done         bool      `json:"done"`
	DoneReason   string    `json:"done_reason"`
	EvalCount    int       `json:"eval_count"`
	EvalDuration int64     `json:"eval_duration"`
	LoadDuration int       `json:"load_duration"`
	Message      struct {
		Content string `json:"content"`
		Role    string `json:"role"`
	} `json:"message"`
	Model              string `json:"model"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int    `json:"prompt_eval_duration"`
	TotalDuration      int64  `json:"total_duration"`
}
