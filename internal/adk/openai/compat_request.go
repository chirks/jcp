package openai

import (
	"errors"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"
)

func buildCompatRetryRequest(req goopenai.ChatCompletionRequest, err error) (goopenai.ChatCompletionRequest, bool) {
	if !isCompatRetryableError(err) {
		return req, false
	}

	changed := false
	if req.MaxTokens > 0 && req.MaxCompletionTokens == 0 {
		req.MaxCompletionTokens = req.MaxTokens
		req.MaxTokens = 0
		changed = true
	}
	if req.Temperature != 0 && req.Temperature != 1 {
		req.Temperature = 1
		changed = true
	}
	if req.TopP != 0 && req.TopP != 1 {
		req.TopP = 1
		changed = true
	}
	if req.N != 0 && req.N != 1 {
		req.N = 1
		changed = true
	}
	if req.PresencePenalty != 0 {
		req.PresencePenalty = 0
		changed = true
	}
	if req.FrequencyPenalty != 0 {
		req.FrequencyPenalty = 0
		changed = true
	}

	return req, changed
}

func isCompatRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, goopenai.ErrReasoningModelMaxTokensDeprecated) ||
		errors.Is(err, goopenai.ErrReasoningModelLimitationsOther) ||
		errors.Is(err, goopenai.ErrO1MaxTokensDeprecated) ||
		errors.Is(err, goopenai.ErrO1BetaLimitationsOther) {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "please use maxcompletiontokens") ||
		strings.Contains(msg, "please use max_completion_tokens") ||
		strings.Contains(msg, "temperature, top_p and n are fixed at 1")
}
