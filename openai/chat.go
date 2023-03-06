package openai

import (
	"encoding/json"
	"github.com/imroc/req/v3"
)

const _chatPath = "/chat/completions"

func (c *Context) CreateChatCompletion(option CreateChatCompletionOption) (ChatCompletionResponse, error) {
	var err error
	var response *req.Response
	request := c.client.R()

	option, err = option.prepared()
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	response, err = request.SetBody(option).Post(_chatPath)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	parsedResponse := ChatCompletionResponse{}
	if err := json.Unmarshal(response.Bytes(), &parsedResponse); err != nil {
		return ChatCompletionResponse{}, &SyntaxError{Raw: response.Bytes()}
	}

	return parsedResponse, nil
}

type ChatCompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Choices []CompletionChoice `json:"choices"`
	Usage   CompletionUsage    `json:"usage"`
}

func (c *ChatCompletionResponse) ToJSON() (string, error) {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type CompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
