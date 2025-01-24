package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
)

const (
	ErrTypeInvalidAuth      pyrin.ErrorType = "INVALID_AUTH"
	ErrTypeApiTokenNotFound pyrin.ErrorType = "API_TOKEN_NOT_FOUND"

	ErrTypeProjectNotFound pyrin.ErrorType = "PROJECT_NOT_FOUND"
	ErrTypeBoardNotFound   pyrin.ErrorType = "BOARD_NOT_FOUND"
	ErrTypeTaskNotFound    pyrin.ErrorType = "TASK_NOT_FOUND"

	ErrTypeUserAlreadyExists pyrin.ErrorType = "USER_ALREADY_EXISTS"
)

func InvalidAuth(message string) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidAuth,
		Message: "Invalid auth: " + message,
	}
}

func ProjectNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeProjectNotFound,
		Message: "Project not found",
	}
}

func BoardNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeBoardNotFound,
		Message: "Board not found",
	}
}

func TaskNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeTaskNotFound,
		Message: "Task not found",
	}
}

func ApiTokenNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeApiTokenNotFound,
		Message: "Api Token not found",
	}
}

func UserAlreadyExists() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeUserAlreadyExists,
		Message: "User already exists",
	}
}
