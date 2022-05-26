package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type malformedRequest struct {
	status int
	msg    string
}

func (m *malformedRequest) Error() string {
	return m.msg
}

func hasHeaderValue(r *http.Request, key string, value string) bool {
	if r.Header.Get(key) != "" {
		for _, v := range r.Header.Values(key) {
			if strings.Contains(v, value) {
				return true
			}
		}
	}

	return false
}

func DecodeRequestBody(r *http.Request, dst any) error {
	if !hasHeaderValue(r, "Content-Type", "application/json") {
		return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: "request Content-Type header is not application/json"}
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return &malformedRequest{status: http.StatusBadRequest, msg: fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &malformedRequest{status: http.StatusBadRequest, msg: "request body contains badly-formed JSON"}

		case errors.As(err, &unmarshalTypeError):
			return &malformedRequest{status: http.StatusBadRequest, msg: fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return &malformedRequest{status: http.StatusBadRequest, msg: fmt.Sprintf("request body contains unknown field %s", fieldName)}

		case errors.Is(err, io.EOF):
			return &malformedRequest{status: http.StatusBadRequest, msg: "request body must not be empty"}

		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return &malformedRequest{status: http.StatusBadRequest, msg: "request body must only contain a single JSON object"}
	}

	return nil
}
