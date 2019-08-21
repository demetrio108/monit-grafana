package config

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type MonitInstance struct {
	Interval int    `yaml:"interval"`
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
}

type Config struct {
	Instances []MonitInstance
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %v", err)
	}
	defer f.Close()

	config_decoder := yaml.NewDecoder(f)
	cfg := Config{}
	err = config_decoder.Decode(&cfg.Instances)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config: %v", err)
	}

	config_yml, err := yaml.Marshal(cfg.Instances)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling config: %v", err)
	}

	log.Printf("Loaded config: \n%s", string(config_yml))
	return &cfg, nil
}
