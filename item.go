package opts

import (
	"errors"
	"fmt"
	"reflect"
)

//item is the structure representing a
//an opt item. it also implements flag.Value
//generically using reflect.
type item struct {
	val       reflect.Value
	name      string
	shortName string
	envName   string
	useEnv    bool
	typeName  string
	help      string
	defstr    string
	slice     bool
	min       int //valid if slice
}

func newItem(val reflect.Value) (*item, error) {
	if !val.CanAddr() {
		return nil, fmt.Errorf("[rfv] cannot address value")
	} else if !val.IsValid() {
		return nil, fmt.Errorf("[rfv] invalid value")
	}
	i := &item{
		val:   val,
		slice: val.Kind() == reflect.Slice,
	}
	t := i.ElemType()
	if t.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("slice elem (%s) cannot be a pointer", t.Kind())
	}
	supported := false
	switch t.Kind() {
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
	// var flagValueType = reflect.TypeOf((*flag.Value)(nil)).Elem()
	//opts.File (string, with auto os.Stat, possibly add file predict -> how filter?)
	if !supported {
		return nil, fmt.Errorf("field type not supported: %s", t.Kind())
	}
	return i, nil
}

func (i *item) ElemType() reflect.Type {
	t := i.val.Type()
	if i.slice {
		t = t.Elem()
	}
	return t
}

func (i *item) String() string {
	if !i.val.IsValid() {
		return ""
	}
	return fmt.Sprintf("%v", i.val.Interface())
}

func (i *item) Set(s string) error {
	//set has two modes, slice and inplace.
	// when slice, create a new zero value, scan into it, append to slice
	// when inplace, take pointer, scan into it
	var ptr reflect.Value
	if i.slice {
		ptr = reflect.New(i.ElemType())
	} else {
		ptr = i.val.Addr()
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
	if i.slice {
		i.val.Set(reflect.Append(i.val, ptr.Elem()))
	}
	//done
	return nil
}
