package utils

import "reflect"

func Entries(input any) map[string]any {
	entries := make(map[string]any)
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := range v.NumField() {
		field := t.Field(i)
		value := v.Field(i)
		entries[field.Name] = value.Interface()
	}

	return entries
}

func ExtractNonNilFieldsByTag(st any, tag string) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(st)
	typ := reflect.TypeOf(st)
	for i := range typ.NumField() {
		field := typ.Field(i)
		dbTag := field.Tag.Get(tag)
		fieldValue := val.Field(i)
		if !fieldValue.IsNil() {
			result[dbTag] = fieldValue.Elem().Interface()
		}
	}
	return result
}
