package adk

import (
	"net/http"
	"strings"

	"github.com/run-bigpig/jcp/internal/adk/openai"
	"github.com/run-bigpig/jcp/internal/models"
)

func setOpenAIChatTokenLimit(body map[string]any, modelName string, mode models.OpenAITokenParamMode, limit int) {
	setOpenAIChatTokenLimitWithResolvedMode(body, openai.ResolveTokenParamMode(modelName, string(mode)), limit)
}

func setOpenAIChatTokenLimitWithResolvedMode(body map[string]any, resolvedMode string, limit int) {
	delete(body, "max_tokens")
	delete(body, "max_completion_tokens")

	switch resolvedMode {
	case openai.TokenParamModeMaxCompletionTokens:
		body["max_completion_tokens"] = limit
	default:
		body["max_tokens"] = limit
	}
}

func alternateOpenAIChatTokenParamMode(resolvedMode string) string {
	switch resolvedMode {
	case openai.TokenParamModeMaxCompletionTokens:
		return openai.TokenParamModeMaxTokens
	default:
		return openai.TokenParamModeMaxCompletionTokens
	}
}

func isOpenAIChatTokenParamRetryable(statusCode int, respBody []byte) bool {
	if statusCode < http.StatusBadRequest {
		return false
	}

	msg := strings.ToLower(string(respBody))
	if strings.Contains(msg, "please use maxcompletiontokens") ||
		strings.Contains(msg, "please use max_completion_tokens") {
		return true
	}

	tokenParamNames := []string{"max_tokens", "max_completion_tokens", "maxcompletiontokens"}
	tokenParamErrors := []string{"unsupported", "not supported", "unknown", "unrecognized", "invalid"}
	for _, tokenParamName := range tokenParamNames {
		if !strings.Contains(msg, tokenParamName) {
			continue
		}
		for _, tokenParamError := range tokenParamErrors {
			if strings.Contains(msg, tokenParamError) {
				return true
			}
		}
	}

	return false
}
