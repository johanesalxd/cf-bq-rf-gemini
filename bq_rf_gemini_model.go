package bqrfgemini

import "encoding/json"

type promptRequest struct {
	PromptInput  string          `json:"prompt_input"`
	Model        string          `json:"model"`
	PromptOutput json.RawMessage `json:"prompt_output"`
}
