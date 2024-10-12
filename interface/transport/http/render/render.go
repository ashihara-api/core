package render

import (
	"encoding/json"
	"net/http"

	"github.com/ashihara-api/core/domain/errors"
)

const (
	// MIME types

	// MIMEApplicationJSON ...
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONCharsetUTF8 ...
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8

	charsetUTF8 = "charset=UTF-8"
)

// JSON ...
func JSON[T any](w http.ResponseWriter, s int, v T) error {
	w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}

// ErrorJSON ...
func ErrorJSON(w http.ResponseWriter, err error) {
	var ec *errors.Cause
	if !errors.As(err, &ec) {
		errors.As(errors.NewCause(err, errors.CaseBackendError), &ec)
	}
	JSON(w, ec.Code(), ec) // nolint: errcheck
}
