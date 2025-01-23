package repository

import (
	"context"
	"encoding/json"
	"magmar/model"
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
func (o *openAPIQaRepository) Ask(ctx context.Context, prompt string) (decision *model.Decision, err error) {
	zlog.With(ctx).Infow(util.LogRepo, "prompt", prompt)
	resp, err := llms.GenerateFromSinglePrompt(ctx, o.conn, prompt, llms.CallOption(func(opts *llms.CallOptions) {
		opts.Model = "gpt-4o"
		opts.JSONMode = true
	}))
	if err != nil {
		zlog.With(ctx).Errorw("Generate from single prompt failed", "err", err)
		return nil, err
	}

	if err := json.Unmarshal([]byte(resp), &decision); err != nil {
		zlog.With(ctx).Errorw("Unmarshal failed")
		return nil, err
	}
	return decision, nil
}
