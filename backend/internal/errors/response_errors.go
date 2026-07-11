package errors

import (
	"fmt"
	"net/http"
)

type ResponseError interface {
	Status() int
}

type BadRequestError struct {
	err error
}

func NewBadRequestError(format string, a ...any) *BadRequestError {
	return &BadRequestError{
		err: fmt.Errorf(format, a...),
	}
}

func (e *BadRequestError) Error() string {
	return e.err.Error()
}

func (e *BadRequestError) Status() int {
	return http.StatusBadRequest
}

type BadGatewayError struct {
	err error
}

func NewBadGatewayError(format string, a ...any) *BadGatewayError {
	return &BadGatewayError{
		err: fmt.Errorf(format, a...),
	}
}

func (e *BadGatewayError) Error() string {
	return e.err.Error()
}

func (e *BadGatewayError) Status() int {
	return http.StatusBadGateway
}

type InternalServerError struct {
	err error
}

func NewInternalServerError(format string, a ...any) *InternalServerError {
	return &InternalServerError{
		err: fmt.Errorf(format, a...),
	}
}

func (e *InternalServerError) Error() string {
	return e.err.Error()
}

func (e *InternalServerError) Status() int {
	return http.StatusInternalServerError
}
