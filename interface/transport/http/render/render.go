package render

import (
	"encoding/json"
	"net/http"
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
func JSON(w http.ResponseWriter, s int, v interface{}) error {
	w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
