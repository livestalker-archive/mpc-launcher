package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strconv"
)

type Config struct {
	MpcPath   string   `yaml:"mpc_path"`
	Args      []string `yaml:"args,flow"`
	MonCount  int      `yaml:"mon_count"`
	StartPort int      `yaml:"start_port"`
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

func (config *Config) GetNArgsFilename(filename string, monNumber int) []string{
	args := make([]string, 0)
	args = append(args, filename)
	realPort := config.StartPort + monNumber
	for ix, el := range config.Args {
		args[ix] = el
	}
	// add port arg
	args = append(args, "/monitor")
	args = append(args, strconv.Itoa(monNumber))
	// add monitor arg
	args = append(args, "/webport")
	args = append(args, strconv.Itoa(realPort))
	return args
}

func (config *Config) GetNArgs(monNumber int) []string{
	args := make([]string, 0)
	realPort := config.StartPort + monNumber
	for _, el := range config.Args {
		args = append(args, el)
	}
	// add port arg
	args = append(args, "/monitor")
	args = append(args, strconv.Itoa(monNumber))
	// add monitor arg
	args = append(args, "/webport")
	args = append(args, strconv.Itoa(realPort))
	return args
}
