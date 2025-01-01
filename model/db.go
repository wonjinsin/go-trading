package model

import (
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm"
)

// DB ...
type DB struct {
	MainDB *gorm.DB
	OpenAI *openai.LLM
}
