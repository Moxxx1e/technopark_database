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
	CodeForumAlreadyExist: {
		Code:         CodeForumAlreadyExist,
		HTTPCode:     http.StatusConflict,
		DebugMessage: "slug already exist in database",
		UserMessage:  "Forum with this slug already exists",
	},
	CodeForumDoesNotExist: {
		Code:         CodeForumDoesNotExist,
		HTTPCode:     http.StatusNotFound,
		DebugMessage: "forum with this slug doesn't exist in database",
		UserMessage:  "Can't find forum with this slug",
	},
	CodeCantDeleteDatabase: {
		Code:         CodeCantDeleteDatabase,
		HTTPCode:     http.StatusInternalServerError,
		DebugMessage: "fail in truncate tables",
		UserMessage:  "Can't delete user data",
	},
	CodeThreadDoesNotExist: {
		Code:         CodeThreadDoesNotExist,
		HTTPCode:     http.StatusNotFound,
		DebugMessage: "fail to select from thread",
		UserMessage:  "Can't find thread with this slug/id",
	},
	CodePostDoesNotExist: {
		Code:         CodePostDoesNotExist,
		HTTPCode:     http.StatusNotFound,
		DebugMessage: "fail to select post by id",
		UserMessage:  "Can't find post with this id",
	},
	CodeVoteDoesNotExist: {
		Code:         CodeVoteDoesNotExist,
		HTTPCode:     http.StatusNotFound,
		DebugMessage: "fail to select vote",
		UserMessage:  "Can't find existed vote",
	},
	CodeVoteAlreadyExist: {
		Code:         CodeVoteAlreadyExist,
		HTTPCode:     http.StatusConflict,
		DebugMessage: "this vote already exists",
		UserMessage:  "Vote already exists",
	},
	CodeParentPostDoesNotExistInThread: {
		Code:         CodeParentPostDoesNotExistInThread,
		HTTPCode:     http.StatusConflict,
		DebugMessage: "Can't find parent post in thread",
		UserMessage:  "Parent post doesn't exist in thread",
	},
	CodeThreadAlreadyExist: {
		Code:         CodeThreadAlreadyExist,
		HTTPCode:     http.StatusConflict,
		DebugMessage: "thread with this slug already exist",
		UserMessage:  "thread already exist",
	},
}
