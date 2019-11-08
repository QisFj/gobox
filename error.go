package gobox

import (
	"errors"
	"fmt"
)

var (
	errNotFound      = errors.New("record not found")
	errFieldNotExist = func(field string) error {
		return fmt.Errorf("field `%s` not exist", field)
	}
	errUnaddressable      = errors.New("using unaddressable value")
	errUpdateMultiRecords = errors.New("update more than one record")
	errFilterError        = func(err error) error {
		return fmt.Errorf("filter error: %s", err)
	}
	errUnexpectedType = func(v interface{}) error {
		return fmt.Errorf("unexpected type: %T", v)
	}
)
