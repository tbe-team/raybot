package executor

import "fmt"

//nolint:revive
type ExecutorError struct {
	parent error
	msg    string
}

func NewExecutorError(parent error, msg string) *ExecutorError {
	return &ExecutorError{
		parent: parent,
		msg:    msg,
	}
}

func (e ExecutorError) Error() string {
	if e.parent != nil {
		return fmt.Sprintf("Msg=%s, Parent=(%v)", e.msg, e.parent)
	}
	return fmt.Sprintf("Msg=%s", e.msg)
}
