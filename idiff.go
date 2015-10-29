package idiff

import (
	"fmt"
	"reflect"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FormatTest formats a diff so it's test friendly.
// You can create your own function like this one to format
// how you want to.
func FormatTest(d *DiffResult) string {
	if reflect.TypeOf(d.A) != reflect.TypeOf(d.B) {
		return fmt.Sprintf("types differ. got: %T, expected: %T\n", d.B, d.A)
	}

	var str []string
	for _, added := range d.Added {
		str = append(str, fmt.Sprintf("%T%s: added: %#v\n", d.A, added.Path, added.Value))
	}
	for _, removed := range d.Removed {
		str = append(str, fmt.Sprintf("%T%s: removed: %#v\n", d.A, removed.Path, removed.Value))
	}
	for _, modified := range d.Modified {
		if reflect.TypeOf(modified.A).Kind() == reflect.Func && reflect.ValueOf(modified.A).IsNil() {
			str = append(str, fmt.Sprintf("%T%s: got: not nil, expected: nil\n", d.A, modified.Path))
		} else if reflect.TypeOf(modified.B).Kind() == reflect.Func && reflect.ValueOf(modified.B).IsNil() {
			str = append(str, fmt.Sprintf("%T%s: got: nil, expected: not nil\n", d.A, modified.Path))
		} else {
			str = append(str, fmt.Sprintf("%T%s: got: %#v, expected: %#v\n", d.A, modified.Path, modified.B, modified.A))
		}
	}

	if len(str) > 0 {
		result := strings.Join(str, "")
		return result[:len(result)-1]
	}
	return ""
}

// Diff does a difference of the two passed interfaces.
func Diff(a, b interface{}) (*DiffResult, bool) {
	dr := newDiffResult(a, b)
	equal := dr.diff(reflect.ValueOf(a), reflect.ValueOf(b), "")
	return dr, equal
}

// Added represents a part of the diff where a thing was added.
type Added struct {
	Path  string
	Value interface{}
	A     interface{}
	B     interface{}
}

// Removed represents a part of the diff where a thing was removed.
type Removed struct {
	Path  string
	Value interface{}
	A     interface{}
	B     interface{}
}

// Modified represents a part of the diff where a thing was modified.
type Modified struct {
	Path string
	A    interface{}
	B    interface{}
}

// DiffResult represents a diff of values of two
// interface{}s
type DiffResult struct {
	Added    []Added
	Removed  []Removed
	Modified []Modified

	A interface{}
	B interface{}
}

func newDiffResult(a, b interface{}) *DiffResult {
	return &DiffResult{A: a, B: b}
}

func (dr *DiffResult) diff(v1, v2 reflect.Value, path string) bool {
	if !v1.IsValid() && !v2.IsValid() {
		return true
	}

	if !v1.IsValid() || !v2.IsValid() {
		dr.Modified = append(dr.Modified, Modified{
			Path: path,
			A:    v1.Interface(),
			B:    v2.Interface(),
		})
		return false
	}

	if v1.Type() != v2.Type() {
		dr.Modified = append(dr.Modified, Modified{
			Path: path,
			A:    v1.Interface(),
			B:    v2.Interface(),
		})
		return false
	}

	kind := v1.Type().Kind()
	equal := true
	switch kind {
	// case reflect.Func:
	// 	if !v1.IsNil() || !v2.IsNil() {
	// 		dr.Modified = append(dr.Modified, Modified{
	// 			Path:   path,
	// 			A: v1.Interface(),
	// 			B: v2.Interface(),
	// 		})
	// 		equal = false
	// 	}
	case reflect.Array, reflect.Slice:
		v1Len := v1.Len()
		v2Len := v2.Len()
		for i := 0; i < min(v1Len, v2Len); i++ {
			local := path + fmt.Sprintf("[%d]", i)
			if eq := dr.diff(v1.Index(i), v2.Index(i), local); !eq {
				equal = false
			}
		}
		if v1Len > v2Len {
			for i := v2Len; i < v1Len; i++ {
				local := path + fmt.Sprintf("[%d]", i)
				dr.Removed = append(dr.Removed, Removed{
					Path:  local,
					Value: v1.Index(i).Interface(),
					A:     v1.Interface(),
					B:     v2.Interface(),
				})
				equal = false
			}
		} else if v1Len < v2Len {
			for i := v1Len; i < v2Len; i++ {
				local := path + fmt.Sprintf("[%d]", i)
				dr.Added = append(dr.Added, Added{
					Path:  local,
					Value: v2.Index(i).Interface(),
					A:     v1.Interface(),
					B:     v2.Interface(),
				})
				equal = false
			}
		}
	case reflect.Ptr:
		v1 = v1.Elem()
		v2 = v2.Elem()
		if !v1.IsValid() && !v2.IsValid() {
			equal = true
		} else if !v1.IsValid() || !v2.IsValid() {
			dr.Modified = append(dr.Modified, Modified{
				Path: path,
				A:    v1.Interface(),
				B:    v2.Interface(),
			})
			equal = false
		} else {
			equal = dr.diff(v1, v2, path)
		}
	case reflect.Struct:
		typ := v1.Type()
		for i := 0; i < typ.NumField(); i++ {
			index := []int{i}
			field := typ.FieldByIndex(index)
			local := path + fmt.Sprintf(".%s", field.Name)
			if eq := dr.diff(v1.FieldByIndex(index), v2.FieldByIndex(index), local); !eq {
				equal = false
			}
		}
	case reflect.Map:
		for _, key := range v1.MapKeys() {
			av := v1.MapIndex(key)
			bv := v2.MapIndex(key)
			local := path + fmt.Sprintf("[%#v]", key.Interface())
			if !bv.IsValid() {
				dr.Removed = append(dr.Removed, Removed{
					Path:  local,
					Value: av.Interface(),
					A:     v1.Interface(),
					B:     v2.Interface(),
				})
				equal = false
			} else if eq := dr.diff(av, bv, local); !eq {
				equal = false
			}
		}
		for _, key := range v2.MapKeys() {
			aI := v1.MapIndex(key)
			if !aI.IsValid() {
				bI := v2.MapIndex(key)
				local := path + fmt.Sprintf("[%#v]", key.Interface())
				dr.Added = append(dr.Added, Added{
					Path:  local,
					Value: bI.Interface(),
					A:     v1.Interface(),
					B:     v2.Interface(),
				})
				equal = false
			}
		}
	default:
		if reflect.DeepEqual(v1.Interface(), v2.Interface()) {
			equal = true
		} else {
			dr.Modified = append(dr.Modified, Modified{
				Path: path,
				A:    v1.Interface(),
				B:    v2.Interface(),
			})
			equal = false
		}
	}
	return equal
}
