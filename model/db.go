package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm"
)

// DB ...
type DB struct {
	MainDB      *gorm.DB
	OpenAI      *openai.LLM
	TelegramBot *tgbotapi.BotAPI
}
