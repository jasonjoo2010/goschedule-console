package types

type JsonResponse map[string]interface{}

// NewNormalResponse creates normal response with data
func NewNormalResponse(src map[string]interface{}) JsonResponse {
	dst := make(JsonResponse, len(src))
	for k, v := range src {
		dst[k] = v
	}
	dst.Err(0, "")
	return dst
}

// NewErrorResponse creates an error response with given errno and msg
func NewErrorResponse(errno int, msg string) JsonResponse {
	dst := make(JsonResponse, 2)
	dst.Err(errno, msg)
	return dst
}

// NewEmptyResponse creates an empty response with errno=0
func NewEmptyResponse() JsonResponse {
	return NewErrorResponse(0, "")
}

func (resp JsonResponse) Err(errno int, msg string) {
	resp["errno"] = errno
	resp.ErrDesc(msg)
}

func (resp JsonResponse) ErrDesc(msg string) {
	resp["errdesc"] = msg
}
