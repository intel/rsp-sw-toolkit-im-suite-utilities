package helper

import "testing"

func TestIsEmpty(t *testing.T) {

	var target string

	if !IsEmpty(target) {
		t.Error("Initial string should be empty")
	}

	target = "Not Empty"
	if IsEmpty(target) {
		t.Error("Set string should not be empty")
	}

	target = ""
	if !IsEmpty(target) {
		t.Error("Cleared string should be empty")
	}
}
