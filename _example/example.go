package main

import (
	"log"

	"github.com/wweir/conf"
	"github.com/wweir/conf/json"
	"github.com/wweir/conf/toml"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	type Config struct {
		Prefix string
		Users  []struct {
			ID       string
			Password string
		}
	}

	cfg := new(Config)
	if err := conf.New(cfg).
		Load(toml.New("fake.toml"), toml.New("example.toml")).
		Over(json.New("example.json")).
		Err(); err != nil {
		log.Fatalln(err)
	}

	log.Println(*cfg)
}
