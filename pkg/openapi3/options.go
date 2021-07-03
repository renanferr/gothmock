package openapi3

type HandlerOptions struct {
	ResponseStatus  int
	ResponseContent string
}

func NewHandlerOptions(status int, content string) *HandlerOptions {
	return &HandlerOptions{
		ResponseStatus:  status,
		ResponseContent: content,
	}
}
