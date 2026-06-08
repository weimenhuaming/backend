package syserrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParamsInvalidResponse(c *gin.Context) {
	c.JSON(http.StatusOK, NewErrWithCode(ParamsInvalid))
}
func SystemErrorResponse(c *gin.Context) {
	c.JSON(http.StatusOK, NewErrWithCode(SystemError))

}
func UnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusOK, NewErrWithCode(Unauthorized))
}
func CustomSystemErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, CustomSystemErrorError(msg))
}

func Response(c *gin.Context, err *Err) {
	c.JSON(http.StatusOK, err)
}

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func New() *Err {
	return &Err{}
}

func NewErr(code int, msg string) *Err {
	return &Err{
		Code: code,
		Msg:  msg,
	}
}

func NewErrWithCode(code int) *Err {
	return &Err{
		Code: code,
		Msg:  MsgMap[code],
	}
}

func ParamsInvalidError() *Err {
	return NewErrWithCode(ParamsInvalid)
}
func SystemErrorError() *Err {
	return NewErrWithCode(SystemError)
}
func UnauthorizedError() *Err {
	return NewErrWithCode(Unauthorized)
}
func CustomSystemErrorError(msg string) *Err {
	return NewErr(SystemError, msg)
}

func (e *Err) String() string {
	return e.Msg
}
