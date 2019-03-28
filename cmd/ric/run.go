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
		absFilePaths := pl.dealFiles()
		workDir := filepath.Dir(absFilePaths[0])

		if len(pl.A.Compile) == 0 {
			switch pl.L {
			case "c":
				pl.A.Compile = []string{"gcc"}
			case "cpp":
				pl.A.Compile = []string{"g++"}
			case "rust":
				pl.A.Compile = []string{"rustc"}
			}
		}
		binName := "a.out"

		args := append(pl.A.Compile, []string{"-o", binName}...)
		args = append(args, absFilePaths...)
		// compile
		stdOut, stdErr, exitErr := goRun(workDir, "", args...)
		if exitErr != nil {
			if _, ok := exitErr.(*exec.ExitError); ok {
				returnStdOut(stdOut, stdErr, errToStr(exitErr))
				exitF("Compile Error")
			}
			exitF("Ric goRun Failed")
		}

		// run
		binPath := filepath.Join(workDir, binName)
		args = append(pl.A.Run, binPath)

		stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
		returnStdOut(stdOut, stdErr, errToStr(exitErr))

	case pl.L == "java":
		runJava(pl)

	case pl.L == "scala":
		runScala(pl)

	case pl.L == "go":
		runGo(pl)

	default:

		if len(pl.A.Run) == 0 {
			if pl.L == "haskell" {
				pl.A.Run = []string{"runhaskell"}
			} else {
				pl.A.Run = []string{pl.L}
			}
		}
		args := pl.A.Run[0:]
		var workDir string
		if len(pl.F) != 0 {
			absFilePaths, err := writeFiles(pl.F)
			if err != nil {
				exitF("Write files failed")
			}
			workDir = filepath.Dir(absFilePaths[0])
			args = append(pl.A.Run[0:], absFilePaths...)
		}

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

func runGo(pl *PayLoad) {
	absFilePaths := pl.dealFiles()
	workDir := filepath.Dir(absFilePaths[0])
	if len(pl.A.Compile) == 0 {
		pl.A.Compile = []string{"go", "build"}
	}
	binName := "main"

	args := append(pl.A.Compile, []string{"-o", binName}...)
	args = append(args, absFilePaths...)
	// compile
	stdOut, stdErr, exitErr := goRun(workDir, "", args...)
	if exitErr != nil {
		if _, ok := exitErr.(*exec.ExitError); ok {
			returnStdOut(stdOut, stdErr, errToStr(exitErr))
			exitF("Compile Error")
		}
		exitF("Ric goRun Failed")
	}

	// run
	binPath := filepath.Join(workDir, binName)
	args = append(pl.A.Run, binPath)

	stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
	returnStdOut(stdOut, stdErr, errToStr(exitErr))
}

func runJava(pl *PayLoad) {
	absFilePaths := pl.dealFiles()
	workDir := filepath.Dir(absFilePaths[0])

	if len(pl.A.Compile) == 0 {
		pl.A.Compile = []string{"javac"}
	}

	args := append(pl.A.Compile, absFilePaths...)

	filename := filepath.Base(absFilePaths[0])

	// compile
	stdOut, stdErr, exitErr := goRun(workDir, "", args...)
	if exitErr != nil {
		returnStdOut(stdOut, stdErr, errToStr(exitErr))
		exitF("Compile Error")
	}

	if len(pl.A.Run) == 0 {
		pl.A.Run = []string{"java"}
	}
	args = append(pl.A.Run, javaClassName(filename))
	stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
	returnStdOut(stdOut, stdErr, errToStr(exitErr))
}

func runScala(pl *PayLoad) {
	absFilePaths := pl.dealFiles()
	workDir := filepath.Dir(absFilePaths[0])

	if len(pl.A.Compile) == 0 {
		pl.A.Compile = []string{"scalac"}
	}

	args := append(pl.A.Compile, absFilePaths...)

	filename := filepath.Base(absFilePaths[0])

	// compile
	stdOut, stdErr, exitErr := goRun(workDir, "", args...)
	if exitErr != nil {
		returnStdOut(stdOut, stdErr, errToStr(exitErr))
		exitF("Compile Error")
	}

	if len(pl.A.Run) == 0 {
		pl.A.Run = []string{"scala"}
	}
	args = append(pl.A.Run, javaClassName(filename))
	stdOut, stdErr, exitErr = goRun(workDir, pl.I, args...)
	returnStdOut(stdOut, stdErr, errToStr(exitErr))
}
