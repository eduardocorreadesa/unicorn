package configuration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Port                string
	Host                string
	Prefix              string
	PathPetNames        string
	PathPetAdj          string
	PathPetCapabilities string
	TimeProcess         int
}

type Config struct {
	Server Server `toml:"server"`
}

func NewEnvConfig(ctx context.Context, env string) Config {

	if len(strings.TrimSpace(env)) == 0 {
		env = "local"
	}

	file := fmt.Sprintf("./config/config-%s.toml", env)

	return NewEnvConfigFile(ctx, file)
}

func NewEnvConfigFile(ctx context.Context, file string) Config {

	var config Config
	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatal("Fail to decode toml config", err)
		panic(err)
	}

	return config
}
