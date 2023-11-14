package helpers

import (
	"errors"
	"strings"
)

var (
	ErrBusinessNotFound = errors.New("business with specified ID/Alias can't be found")
	ErrCategories       = errors.New("invalid category, please go to /api/v1/categories for list of all valid categories")
	ErrTransactions     = errors.New("invalid transaction, please go to /api/v1/transactions for list of all valid transaction type")
	ErrInternal         = errors.New("it's our fault, not yours")
)

type ResponseError struct {
	Err  error
	code int
}

func NewError(err error, code int) ResponseError {
	return ResponseError{err, code}
}

func (e ResponseError) Error() string {
	return e.Err.Error()
}

func (e ResponseError) Code() int {
	return e.code
}

type ValidationError struct {
	msg error
}

func NewValdiationError(err error) ValidationError {
	return ValidationError{
		msg: err,
	}
}

func (ve ValidationError) Error() string {
	return ve.msg.Error()
}

func (ve ValidationError) ErrSlice() []string {
	var errSlice []string
	for _, err := range strings.Split(ve.msg.Error(), ";") {
		errSlice = append(errSlice, strings.TrimSpace(err))
	}

	return errSlice
}
