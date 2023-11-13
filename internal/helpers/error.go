package helpers

import "errors"

var (
	ErrCategories   = errors.New("invalid category, please go to /api/v1/categories for list of all valid categories")
	ErrTransactions = errors.New("invalid transaction, please go to /api/v1/transactions for list of all valid transaction type")
	ErrInternal     = errors.New("it's our fault, not yours")
)

type ResponseError struct {
	Error error
	Code  int
}

func NewError(err error, code int) *ResponseError {
	return &ResponseError{err, code}
}
