package repository

import (
	"context"
	"magmar/config"
	"magmar/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramMsgRepository struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegramMsgRepository ...
func NewTelegramMsgRepository(magmar *config.ViperConfig, bot *tgbotapi.BotAPI) MsgRepository {
	return &telegramMsgRepository{
		bot:    bot,
		chatID: magmar.GetInt64(util.TelegramChatID),
	}
}

// SendMessage ...
func (t *telegramMsgRepository) SendMessage(ctx context.Context, message string) error {
	zlog.With(ctx).Infow(util.LogRepo, "message", message)
	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.bot.Send(msg)
	if err != nil {
		zlog.With(ctx).Errorw("Send message failed", "err", err)
		return err
	}
	return nil
}
