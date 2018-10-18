package conf

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func deepCopy(iface interface{}) interface{} {
	buf := &bytes.Buffer{}
	typ := reflect.TypeOf(iface)
	newIface := reflect.New(typ).Interface()
	if err := gob.NewEncoder(buf).Encode(iface); err != nil {
		panic(err)
	}
	if err := gob.NewDecoder(buf).Decode(newIface); err != nil {
		panic(err)
	}

	return newIface
}
