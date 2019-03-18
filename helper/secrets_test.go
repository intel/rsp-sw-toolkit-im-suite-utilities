package helper

import (
	"strings"
	"testing"
)

func TestGetSecret_ShouldReturnErrWhenSecretNotFound(t *testing.T) {
	_, err := GetSecret("MySecret")
	if err != nil {
		if !strings.ContainsAny(err.Error(), "The system cannot find the path specified.") {
			t.Error("Wrong error message")
		}
	}
}
func TestGetSecret_ShouldReturnSecretWhenSecretFound(t *testing.T) {
	result, err := GetSecret("./secret.txt")
	if err != nil {
		t.Error("Error trying to read secret")
	}
	if result != "hello" {
		t.Error("Incorrect secret value")
	}
}
