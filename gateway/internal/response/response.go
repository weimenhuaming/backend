package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const ()

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

func ErrorUnauthorized(msg string) error {
	return &CodeError{Code: http.StatusUnauthorized, Msg: msg}
}

func ErrorForbidden(mag string) error {
	return &CodeError{Code: http.StatusForbidden, Msg: mag}
}

func ErrorNotFound(mag string) error {
	return &CodeError{Code: http.StatusNotFound, Msg: mag}
}

func ErrorInternalServer(mag string) error {
	return &CodeError{Code: http.StatusInternalServerError, Msg: mag}
}

func ErrorBadRequest(mag string) error {
	return &CodeError{Code: http.StatusBadRequest, Msg: mag}
}

func ErrorConflict(msg string) error {
	return &CodeError{Code: http.StatusConflict, Msg: msg}
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
		body.Code = http.StatusOK
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
