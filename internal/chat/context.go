package chat

import (
	"chatgit-cli/openai"
)

type HistoryItem struct {
	Option      openai.CreateChatCompletionOption `json:"option"`
	ResponseRaw string                            `json:"responseRaw"`
	Response    openai.ChatCompletionResponse     `json:"response"`
	Failed      bool                              `json:"succeed"`
}

type History []HistoryItem

type Context struct {
	APIKey    string
	SessionID string
	OpenAI    *openai.Context
	History   History
}

func (c *Context) Input(text string) (string, error) {
	return c.Chat(text)
}

var _defaultCompletionOption = openai.CreateChatCompletionOption{
	Model: openai.ModelChatGPT,
}

func newDefaultCompletionOption() openai.CreateChatCompletionOption {
	return _defaultCompletionOption
}

func (c *Context) Chat(text string) (string, error) {
	response, err := c.ChatRequest(text)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

func (c *Context) ChatRequest(text string) (openai.ChatCompletionResponse, error) {
	option := newDefaultCompletionOption()
	option.Messages = []openai.Message{{Role: "user", Content: text}}
	response, err := c.OpenAI.CreateChatCompletion(option)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	return response, nil
}
