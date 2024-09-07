package bqrfgemini

import "encoding/json"

type PromptRequest struct {
	PromptInput  string          `json:"promptInput"`
	Model        string          `json:"model"`
	ModelConfig  ModelConfig     `json:"modelConfig"`
	PromptOutput json.RawMessage `json:"promptOutput"`
}

type ModelConfig struct {
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens int32   `json:"maxOutputTokens"`
	TopP            float32 `json:"topP"`
	TopK            int32   `json:"topK"`
}

func newPromptRequest() PromptRequest {
	return PromptRequest{
		PromptInput:  "",
		Model:        "",
		ModelConfig:  newModelConfig(),
		PromptOutput: json.RawMessage(""),
	}
}

func newModelConfig() ModelConfig {
	return ModelConfig{
		Temperature:     1,
		MaxOutputTokens: 1000,
		TopP:            0.95,
		TopK:            1,
	}
}
