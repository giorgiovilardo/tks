package internal

import (
	"embed"
	"log"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

//go:embed config.toml
var configFS embed.FS

func LoadConf() Config {
	k := koanf.New(".")

	configBytes, err := configFS.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("Error reading embedded config: %v", err)
	}

	if err := k.Load(rawbytes.Provider(configBytes), toml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return config
}
