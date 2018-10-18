package conf

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

// Store is conf backend storage
type Store interface {
	Name() string
	PreCheck() bool
	Load(config interface{}) error
}

type Conf struct {
	config interface{}
	err    error
}

func New(config interface{}) *Conf {
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		panic("Please use a pointer type as config reciver")
	}

	return &Conf{
		config: config,
	}
}

func (c *Conf) Load(stores ...Store) *Conf {
	if c.err != nil {
		return c
	}
	for i := range stores {
		if ok := stores[i].PreCheck(); !ok {
			continue
		}
		if err := stores[i].Load(c.config); err != nil {
			c.err = errors.New(stores[i].Name() + ": " + err.Error())
		}
		return c
	}
	c.err = os.ErrNotExist
	return c
}

func (c *Conf) Over(stores ...Store) *Conf {
	if c.err != nil {
		return c
	}
	for i := range stores {
		if ok := stores[i].PreCheck(); !ok {
			continue
		}
		clone := deepCopy(c.config)
		if err := stores[i].Load(clone); err != nil {
			c.err = errors.New(stores[i].Name() + ": " + err.Error())
		}

		j1, err := json.Marshal(c.config)
		if err != nil {
			c.err = err
		}
		j2, err := json.Marshal(clone)
		if err != nil {
			c.err = err
		}
		m := &map[string]interface{}{}
		if err := json.Unmarshal(j1, m); err != nil {
			c.err = err
		}
		if err := json.Unmarshal(j2, m); err != nil {
			c.err = err
		}
		j, err := json.Marshal(m)
		if err != nil {
			c.err = err
		}
		if err := json.Unmarshal(j, c.config); err != nil {
			c.err = err
		}
		return c
	}

	c.err = os.ErrNotExist
	return c
}

func (c *Conf) Err() error {
	return c.err
}
