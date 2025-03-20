package xerror

import (
	"errors"
	"fmt"
)

// XError represents the base error structure.
// Use [NewXError] to create a new instance of XError.
type XError struct {
	parent error // Now use for logging purpose only
	status Status
	msgID  string
	msg    string
}

// NewXError initializes a Base instance.
//
// msgID should be structured as [component].[errorBrief], e.g.
//
//	login.failedAuthentication
//	login.invalidCredentials
func NewXError(parent error, status Status, msgID, msg string) XError {
	return XError{
		parent: parent,
		status: status,
		msgID:  msgID,
		msg:    msg,
	}
}

func (e XError) Error() string {
	if e.parent != nil {
		return fmt.Sprintf("ID=%s, Msg=%s, Parent=(%v)", e.msgID, e.msg, e.parent)
	}
	return fmt.Sprintf("ID=%s, Msg=%s", e.msgID, e.msg)
}

func (e *XError) WithParent(parent error) {
	e.parent = parent
}

func (e *XError) Unwrap() error {
	return e.parent
}

func (e XError) Status() Status {
	return e.status
}

func (e XError) MsgID() string {
	return e.msgID
}

func (e XError) Msg() string {
	return e.msg
}

func (e XError) Parent() error {
	return e.parent
}

// IsStatus checks if the error is an XError with the given status.
func IsStatus(err error, status Status) bool {
	var xErr XError
	if errors.As(err, &xErr) {
		return xErr.Status() == status
	}
	return false
}

func Unauthorized(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusUnauthorized, msgID, msg)
}

func Forbidden(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusForbidden, msgID, msg)
}

func NotFound(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusNotFound, msgID, msg)
}

func UnprocessableEntity(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusUnprocessableEntity, msgID, msg)
}

func Conflict(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusConflict, msgID, msg)
}

func TooManyRequests(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusTooManyRequests, msgID, msg)
}

func BadRequest(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusBadRequest, msgID, msg)
}

func ValidationFailed(parent error, msg string) XError {
	return NewXError(parent, StatusValidationFailed, "validationFailed", msg)
}

func InternalServerError(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusInternalServerError, msgID, msg)
}

func Timeout(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusTimeout, msgID, msg)
}

func NotImplemented(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusNotImplemented, msgID, msg)
}

func BadGateway(parent error, msgID, msg string) XError {
	return NewXError(parent, StatusBadGateway, msgID, msg)
}
