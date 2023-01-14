package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Applications struct {
	Apps map[string]App `yaml:"apps"`
}
type App struct {
	Shh string `yaml:"ssh"`
	Bdd string `yaml:"bdd",omitempty`
	Web string `yaml:"web", omitempty`
}

func NewConfig(path string) (conf *Applications, err error) {
	conf = &Applications{}

	yamlReader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)

	}
	defer yamlReader.Close()

	decoder := yaml.NewDecoder(yamlReader)

	if err = decoder.Decode(conf); err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	return conf, nil
}
