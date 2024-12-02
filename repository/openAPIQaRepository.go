package repository

import (
	"context"

	"github.com/tmc/langchaingo/llms/openai"
)

// OpenAPIQaRepository ...
type openAPIQaRepository struct {
	conn *openai.LLM
}

// NewOpenAPIQaRepository ...
func NewOpenAPIQaRepository(llm *openai.LLM) QaRepository {
	return &openAPIQaRepository{conn: llm}
}

// Ask ...
func (o *openAPIQaRepository) Ask(ctx context.Context) (err error) {
	return nil
}
