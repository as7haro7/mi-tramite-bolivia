package handlers

import (
	"strconv"
)

// parseError representa un error de parseo simple
type parseError struct{}

func (e *parseError) Error() string { return "parse error" }

// parseInt64 convierte un string a int64, devolviendo error si falla.
func parseInt64(s string) (int64, error) {
	if s == "" {
		return 0, &parseError{}
	}
	return strconv.ParseInt(s, 10, 64)
}
