package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type BaseName struct {
	List []int `json:"list"`
}

type Foo struct {
	Names []string `json:"names"`
	Values []int `json:"values"`
	BaseNames []BaseName `json:"base_names"`
	Mapping map[int]BaseName `json:"mapping"`
}


func _initNullSlices(v reflect.Value) {
	fmt.Println(v.Kind())

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i += 1 {
			_initNullSlices(v.Field(i))
		}
	case reflect.Ptr:
		if v.IsValid() {
			_initNullSlices(v.Elem())
		}
	case reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		} else {
			for i := 0; i < v.Len(); i++ {
				_initNullSlices(v.Index(i))
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			_initNullSlices(v.MapIndex(key))
		}
	default:
		panic("unsupported")
	}
}

func nullToTmptySlices(obj interface{}) {
	_initNullSlices(reflect.ValueOf(obj))
}

func main() {
	x := Foo {
		BaseNames: []BaseName{{}},
	}

	isAnythingDeepNull(&x)

	repr, err := json.Marshal(&x)
	if err != nil {
		os.Exit(1)
	}

	fmt.Print("Marshalled: ", string(repr))
}