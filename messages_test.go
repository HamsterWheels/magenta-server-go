package main

import "testing"

func TestTrimMessage(t *testing.T) {
	const s = "Test Sentence\r\n"
	got := trimMessage(s)
	if got == s {
		t.Errorf("%v still contains returns and newlines", s)
	}
}
