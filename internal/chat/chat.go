package chat

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/vvtommy/chatgpt-cli/openai"
)

func newSessionID() string {
	return fmt.Sprintf("chatgpt-cli-%s", ksuid.New().String())
}

type NewChatContextOption struct {
	OpenAI *openai.Context
}

func NewChatContext(option NewChatContextOption) *Context {
	return &Context{
		SessionID: newSessionID(),
		OpenAI:    option.OpenAI,
		History:   History{},
	}
}
