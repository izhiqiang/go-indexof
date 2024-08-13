package config

import (
	"encoding/json"
	"os"
)

var (
	Global *Config
)

type Config struct {
	Port    int     `json:"port"`
	Debug   bool    `json:"debug"`
	IndexOf IndexOf `json:"index_of"`
}

type IndexOf struct {
	Name string `json:"name"`
	Root string `json:"root"`
}

func LoadConfig(file string) error {
	var (
		err error
	)
	fileByte, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileByte, &Global); err != nil {
		return err
	}
	return nil
}
