package helpers

import (
	"encoding/json"
	"net/http"
)

type Sender interface {
	Send(http.ResponseWriter) error
}

type Response[X any] struct {
	Error    string   `json:"error,omitempty"`
	Errors   []string `json:"errors,omitempty"`
	Data     X        `json:"data,omitempty"`
	Metadata any      `json:"metadata,omitempty"`
}

type wrapper[T any] struct {
	code        int
	contentType string
	response    Response[T]
}

func New[T any]() *wrapper[T] {
	return &wrapper[T]{
		code:        http.StatusOK,
		contentType: "application/json",
	}
}

func (wr *wrapper[T]) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", wr.contentType)
	w.WriteHeader(wr.code)
	if err := json.NewEncoder(w).Encode(wr.response); err != nil {
		return err
	}

	return nil
}

func (w *wrapper[T]) Code(code int) *wrapper[T] {
	w.code = code
	return w
}

func (w *wrapper[T]) Error(err error) Sender {
	w.response.Error = err.Error()
	return w
}

func (w *wrapper[T]) Errors(errs []string) Sender {
	w.response.Errors = errs
	return w
}
func (w *wrapper[T]) ContentType(contentType string) *wrapper[T] {
	w.contentType = contentType
	return w
}

func (w *wrapper[T]) Metadata(metadata any) *wrapper[T] {
	w.response.Metadata = metadata
	return w
}

func (w *wrapper[T]) Data(data T) Sender {
	w.response.Data = data
	return w
}
