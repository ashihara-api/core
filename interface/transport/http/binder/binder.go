package binder

import (
	"encoding/json"
	"io"

	"github.com/ashihara-api/core/domain/errors"
)

func FromJSON[T any](body io.ReadCloser, v T) (err error) {
	if err := json.NewDecoder(body).Decode(v); err != nil {
		return errors.NewCause(err, errors.CaseBadRequest)
	}
	return nil
}
