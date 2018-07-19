package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	MpcPath   string   `yaml:"mpc_path"`
	Args      []string `yaml:"args,flow"`
	MonCount  int      `yaml:"mon_count"`
	StartPort int      `yaml:"start_port"`
	WebUIPort int      `yaml:"webui_port"`
}

func LoadConfig(filename string) *Config {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &config)
	return &config
}
