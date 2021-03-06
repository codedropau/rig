package stream

import (
	"encoding/json"
	"fmt"
	"io"
)

// PushStream returned by the Docker daemon.
type PushStream struct {
	Error string `json:"error"`
}

// Push output which was provided by the Docker daemon.
func Push(w io.Writer, stream io.ReadCloser) error {
	d := json.NewDecoder(stream)

	var s *PushStream

	for {
		if err := d.Decode(&s); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		// @todo, Stream push output.

		if s.Error != "" {
			return fmt.Errorf("push failed: %s", s.Error)
		}
	}

	return nil
}
