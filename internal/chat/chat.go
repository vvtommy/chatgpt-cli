package chat

import (
	"chatgit-cli/openai"
	"fmt"
	"github.com/segmentio/ksuid"
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
