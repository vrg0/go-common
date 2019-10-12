package conf

import "testing"

func TestGet(t *testing.T) {
	v, ok := Get("blacklist.url", "kkk")
	if !ok {
		t.Error("not found conf")
	} else {
		t.Log(v)
	}
}

func TestGet2(t *testing.T) {
	v, ok := Get("blacklist.url", "none")
	if ok {
		t.Error(v)
	} else {
		t.Log(v)
	}
}

func TestGetNamespace(t *testing.T) {
	v := GetNamespace("blacklist.url")
	t.Log(v)
}

func TestGetOrDefault(t *testing.T) {
	v := GetOrDefault("blacklist.url", "none", "123")
	if v != "123" {
		t.Log(v)
	}

	v2 := GetOrDefault("blacklist.url", "kkk", "345")
	if v2 == "345" {
		t.Log(v2)
	}
}

func TestRefreshKvMap(t *testing.T) {
	RefreshKvMap(map[string]string{"123":"aaa"})
	v := GetOrDefault("blacklist.url", "kv", "xxx")
	if v != "/aaa" {
		t.Error(v)
	} else {
		t.Log(v)
	}
}