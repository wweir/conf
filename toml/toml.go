package toml

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type Store struct {
	file string
}

func New(file string) *Store {
	return &Store{file: file}
}

func (s *Store) Name() string {
	return "toml"
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
		return errors.New("toml store should not be nil")
	}

	_, err := toml.DecodeFile(s.file, config)
	return err
}
