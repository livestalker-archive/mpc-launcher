package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("./fixtures/config_missing.yml")
	if err == nil {
		t.Error("Config file does not exists but it has not checked yet")
	}
}

func TestLoadConfig2(t *testing.T) {
	_, err := LoadConfig("./fixtures/config_not_yaml.yml")
	if err == nil {
		t.Error("Config file has not yaml format but it has not checked yet")
	}
}
