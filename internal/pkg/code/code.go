package code

import (
	"net/http"

	"github.com/novalagung/gubrak"
	"github.com/rose839/IAM/pkg/errors"
)

// ErrCode implements `github.com/rose839/IAM/pkg/errors`.Coder interface.
type ErrCode struct {
	// C refers to the code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

var _ errors.Coder = &ErrCode{}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return coder.Ref
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.C == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

func register(code int, httpStatus int, message string, ref string) {
	if found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus); !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  ref,
	}

	errors.MustRegister(coder)
}
