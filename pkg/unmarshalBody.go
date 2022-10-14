package pkg

import (
	"encoding/json"
	"io"
)

func Unmarshal(model interface{}, body io.Reader) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &model)
}

