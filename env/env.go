package env

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var errType = errors.New("Expect a pointer to a struct")

func New(prefix string) *Store {
	return &Store{
		prefix: prefix,
	}
}

type Store struct {
	prefix string
}

func (s *Store) Name() string {
	return "env"
}
func (s *Store) PreCheck() bool {
	if s == nil {
		return false
	}
	return true
}

func (s *Store) Load(config interface{}) error {
	if s == nil {
		prefix := filepath.Base(os.Args[0])
		prefix = strings.ToUpper(prefix)
		s = &Store{
			prefix: prefix + "_",
		}
	}

	ptrVal := reflect.ValueOf(config)
	if ptrVal.Kind() != reflect.Ptr {
		return errType
	}
	val := ptrVal.Elem()
	if val.Kind() != reflect.Struct {
		return errType
	}

	return parseStruct(s.prefix, &val)
}

func parseStruct(prefix string, val *reflect.Value) error {
	typ := val.Type()
	fieldCount := val.Type().NumField()
	for i := 0; i < fieldCount; i++ {
		name := parseName(prefix, typ.Field(i))
		env := os.Getenv(name)
		env = strings.TrimSpace(env)

		switch val.Field(i).Kind() {
		case reflect.String:
			val.SetString(env)

		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			x, err := strconv.ParseInt(env, 10, 0)
			if err != nil {
				return errors.New(name + " " + err.Error())
			}
			val.SetInt(x)

		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			x, err := strconv.ParseUint(env, 10, 0)
			if err != nil {
				return errors.New(name + " " + err.Error())
			}
			val.SetUint(x)

		case reflect.Bool:
			switch env {
			case "Y", "y", "YES", "Yes", "yes":
				val.SetBool(true)
			case "N", "n", "NO", "No", "no":
				val.SetBool(false)
			default:
				x, err := strconv.ParseBool(env)
				if err != nil {
					return errors.New(name + " " + err.Error())
				}
				val.SetBool(x)
			}

		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(env, 0)
			if err != nil {
				return errors.New(name + " " + err.Error())
			}
			val.SetFloat(x)

		case reflect.Struct:
			v := val.Field(i)
			parseStruct(name, &v)

		case reflect.Ptr:
			v := val.Field(i).Elem()
			parseStruct(name, &v)
		}
	}
	return nil
}

func parseName(prefix string, field reflect.StructField) string {
	name := field.Tag.Get("env")
	if name != "" {
		return name
	}

	name = field.Name
	if name == "" {
		return prefix
	}

	if strings.Title(name) == name {
		name = strings.ToUpper(name)
	}

	name = strings.TrimLeft(name, "_")
	name = strings.TrimRight(name, "_")

	if prefix == "" {
		return name
	}

	return prefix + "_" + name
}
