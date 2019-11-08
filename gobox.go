package gobox

import (
	"encoding/json"
	"reflect"
	"sync"
)

type Box interface {
	internal()
	Set(v interface{}) error                    // v must be {} or &{}
	Get(v interface{}, filters ...Filter) error // v must be &{} or &[]{}
	Count(v interface{}, filters ...Filter) (int, error)
	Update(v interface{}, filters ...Filter) error
	UpdateAttr(v interface{}, attr map[string]interface{}, filters ...Filter) error
}

func New() Box {
	return &box{
		storage: make(map[string][]interface{}),
	}
}

type box struct {
	rwLock  sync.RWMutex
	storage map[string][]interface{}
}

func (box box) internal() {}

func (box *box) Set(v interface{}) error {
	box.rwLock.Lock()
	defer box.rwLock.Unlock()
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	name := vv.Type().String()
	id := len(box.storage[name]) + 1 // id should start from 1
	err := setField(v, "ID", id)
	if err != nil {
		return err
	}
	box.storage[name] = append(box.storage[name], vv.Interface())
	return nil
}

func (box *box) filter(name string, filters ...Filter) ([]interface{}, error) {
	all := box.storage[name]
	chosenCount := len(all)
	unchosen := make([]bool, chosenCount)
	var chosen bool
	var err error
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		for i, v := range all {
			if unchosen[i] {
				continue
			}
			chosen, err = filter(v)
			if err != nil {
				return nil, errFilterError(err)
			}
			if !chosen {
				unchosen[i] = true
				chosenCount--
				if chosenCount == 0 {
					return nil, nil
				}
			}
		}
	}
	selected := make([]interface{}, 0, chosenCount)
	for i, unchosen := range unchosen {
		if unchosen {
			continue
		}
		selected = append(selected, all[i])
	}
	return selected, nil
}

func (box *box) Get(v interface{}, filters ...Filter) error {
	box.rwLock.RLock()
	defer box.rwLock.RUnlock()
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	var name string
	if vv.Kind() == reflect.Slice {
		name = vv.Type().Elem().String()
	} else {
		name = vv.Type().String()
	}
	chosen, err := box.filter(name, filters...)
	if err != nil {
		return err
	}
	if vv.Kind() == reflect.Slice {
		return deepCopy(chosen, v)
	}
	if len(chosen) == 0 {
		return errNotFound
	}
	return deepCopy(chosen[0], v)
}

func (box box) Count(v interface{}, filters ...Filter) (int, error) {
	box.rwLock.RLock()
	defer box.rwLock.RUnlock()
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	var name string
	if vv.Kind() == reflect.Slice {
		name = vv.Type().Elem().String()
	} else {
		name = vv.Type().String()
	}
	chosen, err := box.filter(name, filters...)
	if err != nil {
		return 0, err
	}
	return len(chosen), nil
}

func (box *box) Update(v interface{}, filters ...Filter) error {
	box.rwLock.Lock()
	defer box.rwLock.Unlock()
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	name := vv.Type().String()
	chosen, err := box.filter(name, filters...)
	if err != nil {
		return err
	}
	if len(chosen) == 0 {
		return errNotFound
	}
	if len(chosen) > 1 {
		return errUpdateMultiRecords
	}
	idI, err := getField(chosen[0], "ID")
	if err != nil {
		return err
	}
	id := idI.(int)
	err = setField(v, "ID", id)
	if err != nil {
		return err
	}
	box.storage[name][id-1] = vv.Interface()
	return nil
}

func (box *box) UpdateAttr(v interface{}, attr map[string]interface{}, filters ...Filter) error {
	box.rwLock.Lock()
	defer box.rwLock.Unlock()
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	name := vv.Type().String()
	chosen, err := box.filter(name, filters...)
	if err != nil {
		return err
	}
	if len(chosen) == 0 {
		return errNotFound
	}
	if len(chosen) > 1 {
		return errUpdateMultiRecords
	}
	idI, err := getField(chosen[0], "ID")
	if err != nil {
		return err
	}
	id := idI.(int)

	p := reflect.New(reflect.TypeOf(chosen[0]))
	p.Elem().Set(reflect.ValueOf(chosen[0]))
	for fieldName, fieldValue := range attr {
		field := p.Elem().FieldByName(fieldName)
		if !field.IsValid() {
			return errFieldNotExist(fieldName)
		}
		field.Set(reflect.ValueOf(fieldValue))
	}
	err = setField(p.Interface(), "ID", id)
	if err != nil {
		return err
	}
	box.storage[name][id-1] = p.Elem().Interface()
	return nil
}

func deepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
