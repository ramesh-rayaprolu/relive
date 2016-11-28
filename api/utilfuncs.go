package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// writeResponse utility for writing to ResponseWriter
func writeResponse(data interface{}, w http.ResponseWriter) error {
	var (
		enc []byte
		err error
	)
	enc, err = json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Failure to marshal, err = %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	n, err := w.Write(enc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Failure to write, err = %s", err)
	}
	if n != len(enc) {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("Short write sent = %d, wrote = %d", len(enc), n)
	}
	return nil
}
