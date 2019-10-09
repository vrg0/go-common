package args

import (
	"testing"
)

func TestGetDefault(t *testing.T) {
	defaultStr := "default"
	if GetDefault("none", defaultStr) != defaultStr {
		t.Error("get default none")
	}
	if GetDefault("Shell", defaultStr) == defaultStr {
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