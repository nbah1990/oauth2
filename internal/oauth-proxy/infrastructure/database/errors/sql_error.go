package errors

import (
	"fmt"
)

type SQLError string

const (
	NoRows SQLError = `sql: no rows in result set`
)

func NewNotFoundError(entity string) error {
	return fmt.Errorf(fmt.Sprintf(`%s not found`, entity))
}
