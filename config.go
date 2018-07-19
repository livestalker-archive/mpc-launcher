package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// Application config
type Config struct {
	// Full path to MPC binary
	MpcPath   string   `yaml:"mpc_path"`
	// General cmd arguments
	Args      []string `yaml:"args,flow"`
	// Monitor count
	MonCount  int      `yaml:"mon_count"`
	// Start port. First MPC instance will be listening on StartPort+1 ...
	StartPort int      `yaml:"start_port"`
	// Port for web interface
	WebUIPort int      `yaml:"webui_port"`
}

// Load application config
func LoadConfig(filename string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
