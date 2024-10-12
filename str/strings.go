// Package str used for structs
package str

import (
	"bytes"
	"reflect"

	"github.com/mfederowicz/trakt-sync/buffer"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/wissance/stringFormatter"
)

var timestampType = reflect.TypeOf(Timestamp{})

// Stringify attempts to create a reasonable string representation of types in
// the Trakt library. It does things like resolve pointers to their values
// and omits struct fields with nil values.
// inspired by the go-github library.
func Stringify(message any) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

// stringifyValue was heavily inspired by the goprotobuf library.
func stringifyValue(w *bytes.Buffer, val reflect.Value) {
	// Handle nil pointers
	if isNilPointer(val) {
		buffer.Write(w, "<nil>")
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		printer.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		stringifyvalueSlice(w, v)
		return
	case reflect.Struct:
		stringifyStructValue(w, v)
	default:
		if v.CanInterface() {
			printer.Fprint(w, v.Interface())
		}
	}
}

func stringifyStructValue(w *bytes.Buffer, v reflect.Value) {
	if v.Type().Name() != "" {
		buffer.Write(w, v.Type().String())
	}

	// special handling of Timestamp values
	if v.Type() == timestampType {
		printer.Fprintf(w, "{%s}", v.Interface())
		return
	}
	stringifyValueStruct(w, v)

}

// Helper function to check if a value is a nil pointer
func isNilPointer(val reflect.Value) bool {
	return val.Kind() == reflect.Ptr && val.IsNil()
}

func stringifyvalueSlice(w *bytes.Buffer, v reflect.Value) {
	buffer.Write(w, "[")
	for i := consts.ZeroValue; i < v.Len(); i++ {
		if i > consts.ZeroValue {
			buffer.Write(w, " ")
		}

		stringifyValue(w, v.Index(i))
	}
	buffer.Write(w, "]")
}

func stringifyValueStruct(w *bytes.Buffer, v reflect.Value) {
	buffer.Write(w, "{")
	var sep bool
	for i := consts.ZeroValue; i < v.NumField(); i++ {
		fv := v.Field(i)
		// Skip nil fields
		if isNilValue(fv) {
			continue
		}

		if sep {
			buffer.Write(w, ", ")
		} else {
			sep = true
		}
		buffer.Write(w, v.Type().Field(i).Name)
		buffer.Write(w, ":")
		stringifyValue(w, fv)
	}
	buffer.Write(w, "}")
}

func isNilValue(fv reflect.Value) bool {
	switch fv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		return fv.IsNil()
	default:
		return false
	}
}

// ContainString check if string exists in slice
func ContainString(key string, s []string) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}

	return false
}

// ContainInt check if int exists in slice
func ContainInt(key int, s []int) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}

	return false
}

// Formatc helper function for FormatComplex in stringFormatter
func Formatc(pattern string, data map[string]any) string {
	return stringFormatter.FormatComplex(pattern, data)
}

// Format helper function for Fomat in stringFormatter
func Format(pattern string, args ...any) string {
	return stringFormatter.Format(pattern, args...)
}
