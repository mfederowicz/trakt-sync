// Package str used for structs
package str

import (
	"reflect"
)

// IDs represents JSON ids object with ids of object from other services
type IDs struct {
	Trakt  *int64  `json:"trakt,omitempty"`
	Slug   *string `json:"slug,omitempty"`
	Imdb   *string `json:"imdb,omitempty"`
	Tmdb   *int    `json:"tmdb,omitempty"`
	Tvdb   *int    `json:"tvdb,omitempty"`
	Tvrage *string `json:"tvrage,omitempty"`
}

// HaveID checks if id for key exists in object
func (i *IDs) HaveID(key string) bool {
	v := reflect.ValueOf(i)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("Input is not a struct or a pointer to a struct")
	}

	field, found := v.Type().FieldByName(key)
	if !found {
		return false // Field not found
	}

	fieldValue := v.FieldByName(key)

	const (
		EmptyStringLen = 0
	)
	// Check if the field is set and not nil or an empty string
	switch field.Type.Kind() {
	case reflect.String:
		return len(fieldValue.String()) != EmptyStringLen
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan:
		return !fieldValue.IsNil()
	default:
		zeroValue := reflect.Zero(field.Type)
		return !reflect.DeepEqual(fieldValue.Interface(), zeroValue.Interface())
	}
}
