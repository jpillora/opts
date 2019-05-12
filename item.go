package opts

import (
	"encoding"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"time"

	"github.com/posener/complete"
)

//item group represents a single "Options" block
//in the help text ouput
type itemGroup struct {
	name  string
	flags []*item
}

const defaultGroup = ""

//item is the structure representing a
//an opt item. it also implements flag.Value
//generically using reflect.
type item struct {
	val       reflect.Value
	mode      string
	name      string
	shortName string
	envName   string
	useEnv    bool
	help      string
	defstr    string
	slice     bool
	min       int //valid if slice
	noarg     bool
	predictor complete.Predictor
	set       bool
}

func newItem(val reflect.Value) (*item, error) {
	if !val.IsValid() {
		return nil, fmt.Errorf("invalid value")
	}
	i := &item{}
	supported := false
	//take interface value, and attempt to
	//make it (or pointer to it) a flag.Value
	v := val.Interface()
	pv := interface{}(nil)
	if val.CanAddr() {
		pv = val.Addr().Interface()
	}
	//convert other types into a flag value:
	if t, ok := v.(texter); ok {
		v = &textValue{t}
	} else if t, ok := pv.(texter); ok {
		v = &textValue{t}
	} else if t, ok := v.(binaryer); ok {
		v = &binaryValue{t}
	} else if t, ok := pv.(binaryer); ok {
		v = &binaryValue{t}
	} else if d, ok := pv.(*time.Duration); ok {
		v = newDurationValue(d)
	} else if fv, ok := pv.(flag.Value); ok {
		v = fv
	}
	//implements flag value?
	if fv, ok := v.(flag.Value); ok {
		supported = true
		//NOTE: replacing val removes our ability to set
		//the value, resolved by flag.Value handling all Set calls.
		val = reflect.ValueOf(fv)
	}
	//implements predictor?
	if p, ok := v.(complete.Predictor); ok {
		i.predictor = p
	}
	//val must be concrete at this point
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	//lock in val
	i.val = val
	i.slice = val.Kind() == reflect.Slice
	//prevent defaults on slices (should vals be appended? should it be reset? how to display defaults?)
	if i.slice && val.Len() > 0 {
		return nil, fmt.Errorf("slices cannot have default values")
	}
	//type checks
	t := i.elemType()
	if t.Kind() == reflect.Ptr {
		return nil, fmt.Errorf("slice elem (%s) cannot be a pointer", t.Kind())
	} else if i.slice && t.Kind() == reflect.Bool {
		return nil, fmt.Errorf("slice of bools not supported")
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String, reflect.Bool:
		supported = true
	}
	//use the inner bool flag, if defined, otherwise if bool
	if bf, ok := v.(interface{ IsBoolFlag() bool }); ok {
		i.noarg = bf.IsBoolFlag()
	} else if t.Kind() == reflect.Bool {
		i.noarg = true
	}
	if !supported {
		return nil, fmt.Errorf("field type not supported: %s", t.Kind())
	}
	return i, nil
}

func (i *item) elemType() reflect.Type {
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
	//can only set singles once
	if i.set && !i.slice {
		return errors.New("already set")
	}
	//set has two modes, slice and inplace.
	// when slice, create a new zero value, scan into it, append to slice
	// when inplace, take pointer, scan into it
	var elem reflect.Value
	if i.slice {
		elem = reflect.New(i.elemType()) //ptr
	} else if i.val.CanAddr() {
		elem = i.val.Addr() //ptr
	} else {
		elem = i.val //possibly a flag.Value
	}
	v := elem.Interface()
	//convert string into value
	if fv, ok := v.(flag.Value); ok {
		//addr implements set
		if err := fv.Set(s); err != nil {
			return err
		}
	} else if elem.Kind() == reflect.Ptr {
		//set addr with scan
		n, err := fmt.Sscanf(s, "%v", v)
		if err != nil {
			return err
		} else if n == 0 {
			return errors.New("could not be parsed")
		}
	} else {
		return errors.New("could not be set")
	}
	//slice? append!
	if i.slice {
		//no pointer elems
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		//append!
		i.val.Set(reflect.Append(i.val, elem))
	}
	//mark item as set!
	i.set = true
	//done
	return nil
}

//IsBoolFlag implements the hidden interface
//documented here https://golang.org/pkg/flag/#Value
func (i *item) IsBoolFlag() bool {
	return i.noarg
}

//noopValue defines a flag value which does nothing
var noopValue = noopValueType(0)

type noopValueType int

func (noopValueType) String() string {
	return ""
}

func (noopValueType) Set(s string) error {
	return nil
}

//textValue wraps [un]marshaller into a flag value
type textValue struct {
	texter
}

type texter interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

func (t textValue) String() string {
	b, err := t.MarshalText()
	if err == nil {
		return string(b)
	}
	return ""
}

func (t textValue) Set(s string) error {
	return t.UnmarshalText([]byte(s))
}

//binaryValue wraps [un]marshaller into a flag value
type binaryValue struct {
	binaryer
}

type binaryer interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func (t binaryValue) String() string {
	b, err := t.MarshalBinary()
	if err == nil {
		return string(b)
	}
	return ""
}

func (t binaryValue) Set(s string) error {
	return t.UnmarshalBinary([]byte(s))
}

//borrowed from the stdlib :)
type durationValue time.Duration

func newDurationValue(p *time.Duration) *durationValue {
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = durationValue(v)
	return nil
}

func (d *durationValue) Get() interface{} {
	return time.Duration(*d)
}

func (d *durationValue) String() string {
	return (*time.Duration)(d).String()
}
