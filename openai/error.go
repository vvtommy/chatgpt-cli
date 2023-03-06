package openai

import "fmt"

// APIError error field in OpenAIError
type APIError struct {
	Message   string `json:"message"`
	ErrorType string `json:"type"`
	Code      string `json:"code"`
}

func (e APIError) Error() string {
	return e.Message
}

// ErrorResponse error response from OpenAI
type ErrorResponse struct {
	ErrorItem APIError `json:"error"`
}

type AuthenticationError struct {
	*APIError
}

type InvalidRequestError struct {
	*APIError
}

type PermissionError struct {
	*APIError
}

type RateLimitError struct {
	*APIError
}

type RetryError struct {
	*APIError
}

type SyntaxError struct {
	Raw []byte
}

func (s SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s", s.Raw)
}
