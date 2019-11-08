package gobox

type Filter func(v interface{}) (bool, error)

func FieldFilter(field string, filter Filter) Filter {
	if filter == nil {
		return nil
	}
	return func(v interface{}) (bool, error) {
		var err error
		v, err = getField(v, field)
		if err != nil {
			return false, err
		}
		return filter(v)
	}
}

func Eq(value interface{}) Filter {
	return func(v interface{}) (b bool, e error) {
		return v == value, nil
	}
}

func In(values ...interface{}) Filter {
	return func(v interface{}) (b bool, e error) {
		for _, value := range values {
			if v == value {
				return true, nil
			}
		}
		return false, nil
	}
}
