package bqrfgemini

type promptRequest struct {
	PromptInput  string `json:"prompt_input"`
	Model        string `json:"model"`
	PromptOutput string `json:"prompt_output"`
}
