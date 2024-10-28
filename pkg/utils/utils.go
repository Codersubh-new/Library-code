package utils

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

func ParseBody(r *http.Request, x interface{}) error {

	if err := r.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	err := decoder.Decode(x, r.PostForm)
	if err != nil {
		return fmt.Errorf("error decoding form data: %w", err)
	}
	return nil
}
