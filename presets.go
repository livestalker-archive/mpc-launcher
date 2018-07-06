package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type Presets struct {
	Presets []Preset `yaml: presets`
}

type Preset struct {
	Name  string   `yaml: "name"`
	Files []string `yaml: "files"`
}

func LoadPresets(filename string) *Presets {
	var presets Presets
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &presets)
	return &presets
}

func main() {
	LoadPresets("./conf/presets.yml")
	fmt.Println("hello")
}
