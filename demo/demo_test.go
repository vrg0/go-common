package demo

import "testing"

func TestDemoF(t *testing.T) {
	t.Log(F("abc"))
	t.Log(F("aaaaac"))
}

func TestDemoF1(t *testing.T) {
	t.Log(F1("abc"))
	t.Log(F1("aba"))
	t.Log(F1("ccaba"))
}

func TestDemoF2(t *testing.T) {
	t.Log(F2("aaa"))
	t.Log(F2("acc"))
	t.Log(F2("aaaa"))
}
