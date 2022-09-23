package api

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	ServerName string `yaml:"ServerName"`
	Env        string `yaml:"Env"`
	Host       string `yaml:"Host"`
	Port       int    `yaml:"Port"`

	MongoConfig MongoDBConfiguration `yaml:"MongoConfig"`
}

type MongoDBConfiguration struct {
	DSN string `yaml:"DSN"`
	DB  string `yaml:"DB"`
}

func (c *Configuration) BindFile(filename string) error {
	contents, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(contents, c)
}
