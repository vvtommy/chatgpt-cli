package openai

import (
	"github.com/imroc/req/v3"
)

const baseUrl = "https://api.openai.com/v1"

type Context struct {
	APIKey     string
	OrganizeID string
	client     *req.Client
}

func afterResponse(_ *req.Client, r *req.Response) error {
	if r.Err != nil {
		return r.Err
	}

	var err error
	if errorResponse, ok := r.ErrorResult().(*ErrorResponse); ok {
		switch r.StatusCode {
		case 400, 404, 405:
			err = InvalidRequestError{&errorResponse.ErrorItem}
		case 401:
			err = AuthenticationError{&errorResponse.ErrorItem}
		case 403:
			err = PermissionError{&errorResponse.ErrorItem}
		case 409:
			err = RetryError{&errorResponse.ErrorItem}
		case 429:
			err = RateLimitError{&errorResponse.ErrorItem}
		default:
			err = errorResponse.ErrorItem
		}
		return err
	}
	return r.Err
}
func NewOpenAIContext(APIKey string, organizeID string) *Context {
	client := req.NewClient().
		SetCommonBearerAuthToken(APIKey).
		SetCommonHeader("OpenAI-Organization", organizeID).
		SetBaseURL(baseUrl).
		DisableAutoDecode().
		SetCommonErrorResult(&ErrorResponse{}).
		OnAfterResponse(afterResponse)

	return &Context{
		APIKey:     APIKey,
		OrganizeID: organizeID,
		client:     client,
	}
}
