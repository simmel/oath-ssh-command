package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	config := find_config()
	if config != "lol" {
		t.Error("Expected lol, got", config)
	}
}
