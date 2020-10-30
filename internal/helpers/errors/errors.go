package errors

import (
	. "github.com/technopark_database/internal/consts"
	"net/http"
)

type Error struct {
	Code         ErrorCode `json:"code"`
	HTTPCode     int       `json:"-"`
	DebugMessage string    `json:"debug_message"`
	UserMessage  string    `json:"message"`
}

var WrongErrorCode = &Error{
	HTTPCode:     http.StatusTeapot,
	DebugMessage: "wrong error code",
	UserMessage:  "Что-то пошло не так",
}

func New(code ErrorCode, err error) *Error {
	customErr, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	customErr.DebugMessage = err.Error()
	return customErr
}

func Get(code ErrorCode) *Error {
	err, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	return err
}

var Errors = map[ErrorCode]*Error{
	CodeBadRequest: {
		Code:         CodeBadRequest,
		HTTPCode:     http.StatusBadRequest,
		DebugMessage: "wrong request data",
		UserMessage:  "Incorrect format of request",
	},
	CodeInternalServerError: {
		Code:         CodeInternalServerError,
		HTTPCode:     http.StatusInternalServerError,
		DebugMessage: "something went wrong",
		UserMessage:  "Error on server",
	},
	CodeUserDoesNotExist: {
		Code:         CodeUserDoesNotExist,
		HTTPCode:     http.StatusNotFound,
		DebugMessage: "user with this nickname doesn't exist",
		UserMessage:  "Can't find user with this nickname",
	},
	CodeUserNicknameConflicts: {
		Code:         CodeUserEmailConflicts,
		HTTPCode:     http.StatusConflict,
		DebugMessage: "user with this nickname already exists",
		UserMessage:  "Input nickname already exists",
	},
	CodeUserEmailConflicts: {
		Code:         CodeUserEmailConflicts,
		HTTPCode:     http.StatusConflict,
		DebugMessage: "user with this email already exists",
		UserMessage:  "Input email already exists",
	},
}
