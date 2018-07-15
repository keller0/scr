package internal

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

	if GetEnv(testName, "good") != "good" {
		t.Error("GetEnv fall back error")
	}

	err = os.Setenv(testName, "good")
	defer os.Unsetenv(testName)
	if err != nil {
		t.Error(err)
	}
	if GetEnv(testName, "foo") != "good" {
		t.Error("GetEnv error")
	}

}
