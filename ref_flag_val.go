package opts

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
)

var flagValueType = reflect.TypeOf((*flag.Value)(nil)).Elem()

func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func newReflectFlagVal(v reflect.Value) (flag.Value, error) {
	if !v.CanAddr() {
		return nil, fmt.Errorf("[rfv] cannot address value")
	} else if !v.IsValid() {
		return nil, fmt.Errorf("[rfv] invalid value")
	}
	t := v.Type()
	etype := t
	slice := t.Kind() == reflect.Slice
	if slice {
		etype = t.Elem()
		if etype.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("slice elem (%s) cannot be a pointer", etype.Kind())
		}
	}
	supported := false
	switch etype.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String, reflect.Bool:
		supported = true
	}
	//TODO
	//Text Unmarsal
	//time.Duration
	//time.Time (hardcoded RFC3339Nano)
	//inner flag.Value
	//opts.File (string, with auto os.Stat, possibly add file predict -> how filter?)
	if !supported {
		return nil, fmt.Errorf("field type not supported: %s", etype.Kind())
	}
	return &reflectFlagVal{
		slice: slice,
		etype: etype,
		v:     v,
	}, nil
}

type reflectFlagVal struct {
	slice bool
	etype reflect.Type
	v     reflect.Value
}

func (r *reflectFlagVal) String() string {
	if !r.v.IsValid() {
		return ""
	}
	return fmt.Sprintf("%v", r.v.Interface())
}

func (r *reflectFlagVal) Set(s string) error {
	//set has two modes, slice and inplace.
	// when slice, create a new zero value, scan into it, append to slice
	// when inplace, take pointer, scan into it
	var ptr reflect.Value
	if r.slice {
		ptr = reflect.New(r.etype)
	} else {
		ptr = r.v.Addr()
	}
	addr := ptr.Interface()
	//scan into this address
	n, err := fmt.Sscanf(s, "%v", addr)
	if err != nil {
		return err
	} else if n == 0 {
		return errors.New("could not be parsed")
	}
	//slice? append!
	if r.slice {
		v := ptr.Elem()
		r.v.Set(reflect.Append(r.v, v))
	}
	//done
	return nil
}
