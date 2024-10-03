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
	if val.Kind() == reflect.Ptr && val.IsNil() {
		buffer.Write(w, "<nil>")
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		printer.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		buffer.Write(w, "[")
		for i := consts.ZeroValue; i < v.Len(); i++ {
			if i > consts.ZeroValue {
				buffer.Write(w, " ")
			}

			stringifyValue(w, v.Index(i))
		}
		buffer.Write(w, "]")
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			buffer.Write(w, v.Type().String())
		}

		// special handling of Timestamp values
		if v.Type() == timestampType {
			printer.Fprintf(w, "{%s}", v.Interface())
			return
		}
		buffer.Write(w, "{")
		var sep bool
		for i := consts.ZeroValue; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Map && fv.IsNil() {
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
	default:
		if v.CanInterface() {
			printer.Fprint(w, v.Interface())
		}
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
