package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// PayLoad came from stdin,
// version is not needed.
type PayLoad struct {
	F []*file   `json:"files"`
	A *argument `json:"argument"`
	I string    `json:"stdin"`
	L string    `json:"language"`
}

type argument struct {
	Compile []string `json:"compile"`
	Run     []string `json:"run"`
}

type file struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Result of user's program
type Result struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitError string `json:"exiterror"`
}

func main() {

	var ar PayLoad
	var err error
	b, err := ioutil.ReadAll(os.Stdin)
	err = json.Unmarshal(b, &ar)
	if err != nil {
		exitF(err.Error())
	}

	ar.L = strings.ToLower(ar.L)
	if !ar.isSupport() {
		exitF("language %s is not supported.", ar.L)
	}

	// get result from user's program

	if ar.needCompile() {
		ar.compileAndRun()
	} else {
		ar.Run()
	}
}

func exitF(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

// send user's result to container's stdout
func returnStdOut(uOut, uErr, exError string) {
	result := &Result{
		Stdout:    uOut,
		Stderr:    uErr,
		ExitError: exError,
	}
	b, _ := json.Marshal(result)
	fmt.Print(string(b))
}

func writeFiles(files []*file) ([]string, error) {

	tmpPath, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(files), len(files))
	for i, file := range files {

		path, err := writeFile(tmpPath, file)
		if err != nil {
			return nil, err
		}

		paths[i] = path
	}
	return paths, nil
}

func writeFile(basePath string, file *file) (string, error) {

	absPath := filepath.Join(basePath, file.Name)

	err := os.MkdirAll(filepath.Dir(absPath), 0775)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(absPath, []byte(file.Content), 0664)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func errToStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
