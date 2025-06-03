package common

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

func DecodeJSONFromBytes(content []byte, v interface{}) error {
	r := bytes.NewReader(content)
	return errors.WithStack(json.NewDecoder(r).Decode(v))
}
