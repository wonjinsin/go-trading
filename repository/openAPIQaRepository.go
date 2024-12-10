package repository

import (
	"context"
	"fmt"
	"magmar/util"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type openAPIQaRepository struct {
	conn *openai.LLM
}

// NewOpenAPIQaRepository ...
func NewOpenAPIQaRepository(llm *openai.LLM) QaRepository {
	return &openAPIQaRepository{conn: llm}
}

// Ask ...
func (o *openAPIQaRepository) Ask(ctx context.Context) (err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	prompt := "Hello"
	completion, err := llms.GenerateFromSinglePrompt(ctx, o.conn, prompt)
	if err != nil {
		zlog.With(ctx).Errorw("Generate from single prompt failed")
		return err
	}
	fmt.Println(completion)
	return nil
}
