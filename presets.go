package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Presets []Preset

type Preset struct {
	Name  string   `yaml:"name"`
	Files []PresetFile `yaml:"files"`
}

type PresetFile struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args,flow"`
}

func LoadPresets(filename string) Presets {
	var presets Presets
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &presets)
	return presets
}

func (preset * PresetFile) GetFullArgs() []string {
	args := make([]string, 0)
	args = append([]string{preset.Name}, preset.Args...)
	return args
}