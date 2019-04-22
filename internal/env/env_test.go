package env

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	testName := "testGetEnv"
	err := os.Unsetenv(testName)
	if err != nil {
		t.Error(err)
	}

	if Get(testName, "good") != "good" {
		t.Error("env.Get fall back error")
	}

	err = os.Setenv(testName, "good")
	defer os.Unsetenv(testName)
	if err != nil {
		t.Error(err)
	}
	if Get(testName, "foo") != "good" {
		t.Error("env.Get error")
	}

}
