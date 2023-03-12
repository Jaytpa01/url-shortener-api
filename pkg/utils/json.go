package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jaytpa01/url-shortener-api/api"
	"github.com/go-chi/render"
)

const (
	MAX_REQUEST_BODY_SIZE = 1048576
)

// DecoceJSON tries to read a JSON request body and unmarshal it into dst.
// It expects the request to have the correct content type header,
// and checks for various errors and validates there are no invalid json fields in the request.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	if render.GetRequestContentType(r) != render.ContentTypeJSON {
		return api.NewUnsupportedMediaType(
			"incorrect-content-header",
			"Content-Type header is not application/json",
			api.WithAction(`Ensure the Content-Type header is "application/json".`),
		)
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_REQUEST_BODY_SIZE)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return api.NewBadRequest("json-syntax-error", msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return api.NewBadRequest("bad-json-eof-error", msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return api.NewBadRequest("invalid-field-value", msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return api.NewBadRequest("unknown-field", msg)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return api.NewBadRequest("no-empty-requests", msg)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return api.NewRequestPayloadTooLarge("request-payload-too-large", msg)

		default:
			return err
		}
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return api.NewBadRequest("too-many-json-objects", msg)
	}

	return nil
}
