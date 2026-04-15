package core

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type b64string struct {
	Bytes []byte
}

func (b64 *b64string) UnmarshalText(text []byte) error {
	bytes := make([]byte, len(text))
	_, err := base64.StdEncoding.Decode(bytes, text)
	if err != nil {
		return fmt.Errorf("failed to decode base64 value: %s", err)
	}

	b64.Bytes = bytes
	return nil
}

type AppConfig struct {
	BaseUrl    string    `toml:"base_url"`
	SessionKey b64string `toml:"session_key"`
}

type AuthConfig struct {
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

type DatabaseConfig struct {
	Url string
}

type ServerConfig struct {
	Bind string
	Port uint16
}

type Config struct {
	App      AppConfig
	Auth     AuthConfig
	Database DatabaseConfig
	Server   ServerConfig
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file %s: %s", path, err.Error())
	}
	defer f.Close()

	var config Config
	decoder := toml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file %s: %s", path, err.Error())
	}

	return &config, nil
}
