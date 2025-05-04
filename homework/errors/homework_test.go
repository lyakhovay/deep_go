package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}
	strError := strings.Builder{}
	strError.WriteString(strconv.Itoa(len(e.errors)))
	strError.WriteString(" errors occured:\n")
	for _, err := range e.errors {
		strError.WriteString("\t* " + err.Error())
	}
	strError.WriteString("\n")
	return strError.String()
}

func Append(err error, errs ...error) *MultiError {
	merr, ok := err.(*MultiError)
	if !ok {
		merr = &MultiError{errs}
		return merr
	}
	merr.errors = append(merr.errors, errs...)
	return merr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
