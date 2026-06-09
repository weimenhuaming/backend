package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	CodeOK = 200
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// CodeError represents a business error returned in the response body.
type CodeError struct {
	Code int
	Msg  string
}

func (e *CodeError) Error() string {
	return e.Msg
}

// NewError 处理错误请求
func NewError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// Response 统一处理响应
func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		if codeErr, ok := err.(*CodeError); ok {
			body.Code = codeErr.Code
			body.Msg = codeErr.Msg
		} else {
			body.Code = -1
			body.Msg = err.Error()
		}
	} else {
		body.Code = CodeOK
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
