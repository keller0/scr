package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

var supportedLanguage = []string{
	"bash",
	"c",
	"cpp",
	"go",
	"haskell",
	"java",
	"perl",
	"php",
	"python",
	"scala",
	"ruby",
	"rust",
}

func goRun(workDir, stdin string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// args[0] is the program name
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func (pl *PayLoad) runCode() {

	switch {
	case pl.L == "c" || pl.L == "cpp" || pl.L == "rust":
		args := []string{}
		absFilePaths := pl.dealFiles()
		workDir := filepath.Dir(absFilePaths[0])

		switch pl.L {
		case "c":
			args = appendArgs(pl.A.Compile, "gcc")
		case "cpp":
			args = appendArgs(pl.A.Compile, "g++")
		case "rust":
			args = appendArgs(pl.A.Compile, "rustc")
		}

		binName := "a.out"
		args = append(args, []string{"-o", binName}...)
		args = append(args, absFilePaths...)
		// compile
		stdOut, stdErr, exitErr := goRun(workDir, "", args...)
		if exitErr != nil {
			returnStdOut(stdOut, stdErr, errToStr(exitErr))
			exitF("Compile Error")
		}

		// run
		binPath := filepath.Join(workDir, binName)
		args = appendArgs(pl.A.Run, binPath)

		stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
		returnStdOut(stdOut, stdErr, errToStr(exitErr))

	case pl.L == "java":
		runJava(pl)

	case pl.L == "scala":
		runScala(pl)

	default:

		var args = []string{}
		switch pl.L {
		case "haskell":
			args = appendArgs(pl.A.Compile, "runhaskell")
		case "go":
			// maybe there is a better way
			args = appendArgs(pl.A.Run, "run")
			args = append([]string{"go"}, args...)

		default:
			args = appendArgs(pl.A.Run, pl.L)
		}

		absFilePaths := pl.dealFiles()
		workDir := filepath.Dir(absFilePaths[0])
		args = append(args, absFilePaths...)

		stdOut, stdErr, exitErr := goRun(workDir, pl.I, args...)
		returnStdOut(stdOut, stdErr, errToStr(exitErr))
	}
}

func (pl *PayLoad) isSupport() bool {
	for _, l := range supportedLanguage {
		if pl.L == l {
			return true
		}
	}
	return false
}

func javaClassName(filename string) string {
	ext := filepath.Ext(filename)
	return filename[0 : len(filename)-len(ext)]
}

// deal files, save file and return those absolute paths
func (pl *PayLoad) dealFiles() []string {
	if len(pl.F) == 0 {
		exitF("No files given")
	}
	absFilePaths, err := writeFiles(pl.F)
	if err != nil {
		exitF("Write files failed")
	}
	return absFilePaths
}

func runJava(pl *PayLoad) {
	absFilePaths := pl.dealFiles()
	workDir := filepath.Dir(absFilePaths[0])

	args := appendArgs(pl.A.Compile, "javac")
	args = append(args, absFilePaths...)

	filename := filepath.Base(absFilePaths[0])

	// compile
	stdOut, stdErr, exitErr := goRun(workDir, "", args...)
	if exitErr != nil {
		returnStdOut(stdOut, stdErr, errToStr(exitErr))
		exitF("Compile Error")
	}

	args = appendArgs(pl.A.Run, "java")
	args = append(args, javaClassName(filename))

	stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
	returnStdOut(stdOut, stdErr, errToStr(exitErr))
}

func runScala(pl *PayLoad) {
	absFilePaths := pl.dealFiles()
	workDir := filepath.Dir(absFilePaths[0])

	args := appendArgs(pl.A.Compile, "scalac")
	args = append(args, absFilePaths...)

	filename := filepath.Base(absFilePaths[0])

	// compile
	stdOut, stdErr, exitErr := goRun(workDir, "", args...)
	if exitErr != nil {
		returnStdOut(stdOut, stdErr, errToStr(exitErr))
		exitF("Compile Error")
	}

	args = appendArgs(pl.A.Run, "scala")
	args = append(args, javaClassName(filename))

	stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
	returnStdOut(stdOut, stdErr, errToStr(exitErr))
}

// if arguments is empty use default argument, or the first argument
// is not right, put them behind default
func appendArgs(args []string, def string) []string {

	if argIsEmpty(args) {
		args = []string{def}
	} else if args[0] != def {
		args = append([]string{def}, args...)
	}
	return args
}

func argIsEmpty(args []string) bool {

	if len(args) == 0 {
		return true
	}

	var b = []string{}
	for _, a := range args {
		if a != "" {
			b = append(b, a)
		}
	}
	if len(b) == 0 {
		return true
	}
	return false
}
