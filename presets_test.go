package main

import (
	"testing"
	"strings"
)

func TestLoadPresets(t *testing.T) {
	_, err := LoadPresets("./fixtures/presets_missing.yml")
	if err == nil {
		t.Error("Preset file does not exists but it has not checked yet")
	}
}

func TestLoadPresets2(t *testing.T) {
	p, err := LoadPresets("./fixtures/presets.yml")
	if err != nil {
		t.Error(err)
	}
	if len(p) != 2 || len(p[0].Files) != 2 {
		t.Error("Wrong fixture")
	}
}

func TestLoadPresets3(t *testing.T) {
	_, err := LoadPresets("./fixtures/presets_not_yaml.yml")
	if err == nil {
		t.Error("Preset file has not yaml format but it has not checked yet")
	}
}

func TestPresetFile_GetFullArgs(t *testing.T) {
	pf := PresetFile{
		Name: "test preset",
		Args: []string{"arg1", "arg2", "arg3"},
	}
	args1 := "test preset arg1 arg2 arg3"
	args2 := pf.GetFullArgs()
	if args1 != strings.Join(args2, " ") {
		t.Error("Wrong full arguments slice")
	}
}
