package main

import (
	"testing"
	"fmt"
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
	if len(p) != 2 || len(p[0].Files) != 2{
		t.Error("Wrong fixture")
	}
}

func TestLoadPresets3(t *testing.T) {
	_, err := LoadPresets("./fixtures/presets_not_yaml.yml")
	if err == nil {
		t.Error("Preset file has not yaml format but it has not checked yet")
	}
	fmt.Println(err)
}
