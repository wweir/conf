package json

import (
	"encoding/json"
	"errors"
	"os"
)

type Store struct {
	file string
}

func New(file string) *Store {
	return &Store{file: file}
}

func (s *Store) Name() string {
	return "json"
}

func (s *Store) PreCheck() bool {
	if s == nil {
		return false
	}
	_, err := os.Stat(s.file)
	return err == nil
}

func (s *Store) Load(config interface{}) error {
	if s == nil {
		return errors.New("json store should not be nil")
	}

	f, err := os.Open(s.file)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(config)
}
