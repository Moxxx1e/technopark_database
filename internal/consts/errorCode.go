package consts

type ErrorCode uint16

const (
	CodeBadRequest = iota + 101
	CodeInternalServerError
	CodeUserDoesNotExist
	CodeUserNicknameConflicts
	CodeUserEmailConflicts
)
