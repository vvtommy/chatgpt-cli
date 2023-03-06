package openai

import "github.com/pkg/errors"

type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

var _roles = map[MessageRole]struct{}{
	RoleSystem:    {},
	RoleUser:      {},
	RoleAssistant: {},
}

type ModelType string

const (
	ModelChatGPT ModelType = "gpt-3.5-turbo"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

func (m *Message) Check() error {
	if m.Role == "" {
		return errors.New("message.role is required")
	}
	if _, exists := _roles[m.Role]; !exists {
		return errors.Errorf("invalid role: %s", m.Role)
	}
	if m.Content == "" {
		return errors.New("message.content is required")
	}

	return nil
}

type CreateChatCompletionOption struct {
	Model            ModelType          `json:"model"`
	Messages         []Message          `json:"messages"`
	Temperature      *float32           `json:"temperature,omitempty"`
	TopProbability   *float32           `json:"top_p,omitempty"`
	ChoicesNumber    *uint              `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	MaxTokens        *uint              `json:"max_tokens,omitempty"`
	PresencePenalty  *float32           `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float32           `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float32 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
}

const (
	TemperatureMin = 0
	TemperatureMax = 2
)

const (
	TopProbabilityMin = 0
	TopProbabilityMax = 1
)

const (
	PresencePenaltyMin = -2
	PresencePenaltyMax = 2
)
const (
	FrequencyPenaltyMin = -2
	FrequencyPenaltyMax = 2
)

func inRangeIfPresent[R ~int | ~float32](value *R, min, max R) bool {
	if value == nil {
		return true
	}
	if *value < min || *value > max {
		return false
	}
	return true
}

func (o CreateChatCompletionOption) prepared() (CreateChatCompletionOption, error) {
	if o.Model != ModelChatGPT {
		return CreateChatCompletionOption{}, errors.Errorf("unsupported model: %s", o.Model)
	}

	for i := range o.Messages {
		if err := o.Messages[i].Check(); err != nil {
			return CreateChatCompletionOption{}, nil
		}
	}

	if !inRangeIfPresent(o.Temperature, TemperatureMin, TemperatureMax) {
		return CreateChatCompletionOption{}, errors.Errorf("invalid temperature: %f", o.Temperature)
	}

	if !inRangeIfPresent(o.TopProbability, TopProbabilityMin, TopProbabilityMax) {
		return CreateChatCompletionOption{}, errors.Errorf("invalid top_p: %f", o.TopProbability)
	}

	if o.ChoicesNumber != nil && *o.ChoicesNumber < 1 {
		return CreateChatCompletionOption{}, errors.Errorf("at least one choice")
	}

	// Always false
	o.Stream = false

	if !inRangeIfPresent(o.PresencePenalty, PresencePenaltyMin, PresencePenaltyMax) {
		return CreateChatCompletionOption{}, errors.Errorf("invalid PresencePenalty: %f", o.PresencePenalty)
	}

	if !inRangeIfPresent(o.FrequencyPenalty, FrequencyPenaltyMin, FrequencyPenaltyMax) {
		return CreateChatCompletionOption{}, errors.Errorf("invalid FrequencyPenalty: %f", o.FrequencyPenalty)
	}

	return o, nil
}
