package apperrorcodeimpl

import (
	"github.com/tbe-team/raybot/internal/services/apperrorcode"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/xerror"
)

func init() {
	register(command.ErrCommandNotFound)
	register(command.ErrNoNextExecutableCommand)
	register(command.ErrNoCommandBeingProcessed)
	register(command.ErrCommandInProcessingCanNotBeDeleted)
}

var errorCodes = []apperrorcode.ErrorCode{}

func register(err xerror.XError) {
	errorCodes = append(errorCodes, apperrorcode.ErrorCode{
		Code:    err.MsgID(),
		Message: err.Msg(),
	})
}

func GetAll() []apperrorcode.ErrorCode {
	return errorCodes
}
