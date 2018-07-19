// Module for working with preset configuration yaml file.
// For more details see conf/presets.yml.dist
package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// Slice of presets
type Presets []Preset

// Preset struct
type Preset struct {
	// Human readable name of preset
	Name  string   `yaml:"name"`
	// Options for light control
	Light Light `yaml:"light"`
	// Files for playing
	Files []PresetFile `yaml:"files"`
}

// Video file
type PresetFile struct {
	// Full path to video file
	Name string   `yaml:"name"`
	// Additional arguments for MPC player
	Args []string `yaml:"args,flow"`
}

// Light options
type Light struct{
	// Start light preset after N second
	Time int `yaml:"time"`
	// Light scene number
	Number int `yaml:"number"`
}

// Load presets file
func LoadPresets(filename string) (Presets, error) {
	var presets Presets
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &presets)
	if err != nil {
		return nil, err
	}
	return presets, nil
}

func (pf * PresetFile) GetFullArgs() []string {
	args := make([]string, len(pf.Args) + 1)
	args[0] = pf.Name
	copy(args[1:], pf.Args)
	return args
}