package args

import (
	"testing"
)

func TestGetOrDefault(t *testing.T) {
	defaultStr := "default"
	if GetOrDefault("none", defaultStr) != defaultStr {
		t.Error("get default none")
	}
	if GetOrDefault("Shell", defaultStr) == defaultStr {
		t.Error("get default shell")
	}
}

func TestGet(t *testing.T) {
	if _, ok := Get("none"); ok {
		t.Error("get default none")
	}
	if _, ok := Get("shell"); !ok {
		t.Error("get default shell")
	}
}