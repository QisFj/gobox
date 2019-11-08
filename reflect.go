package gobox

import (
	"reflect"
)

func setField(v interface{}, fieldName string, fieldValue interface{}) (err error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		return errUnexpectedType(v)
	}
	field := vv.FieldByName(fieldName)
	if !field.IsValid() {
		return errFieldNotExist(fieldName)
	}
	if !field.CanSet() {
		return errUnaddressable
	}
	field.Set(reflect.ValueOf(fieldValue))
	return nil
}

func getField(v interface{}, fieldName string) (fieldValue interface{}, err error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		return nil, errUnexpectedType(v)
	}
	field := vv.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, errFieldNotExist(fieldName)
	}
	return field.Interface(), nil
}
