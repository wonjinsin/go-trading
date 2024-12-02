package model

import "github.com/tmc/langchaingo/llms/openai"

// DB ...
type DB struct {
	OpenAI *openai.LLM
}

// WithOpenAI ...
func (db *DB) WithOpenAI() *openai.LLM {
	return db.OpenAI
}
