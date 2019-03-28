package main

import (
	"encoding/json"
	"testing"
)

func TestRunC(t *testing.T) {

	j := `
{
	"files" : [
	   {
		  "content" : "#include <stdio.h>\n\nint main(void) {\n    printf(\"C: Hello, World!\\n\");\n    return 0;\n}",
		  "name" : "main.c"
	   }
	],
	"language" : "c",
	"argument" : {
		"compile":[
			"gcc"
		],
		"run" :[]
	}
 }
`
	var pl PayLoad
	err := json.Unmarshal([]byte(j), &pl)
	if err != nil {
		t.Error(err)
	}
	pl.runCode()

}

func TestRunPHP(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"main.php",
		"content":"<?php\n    echo \"PHP: Hello, World!\"; exit(10);"
		}
	],
	"language" : "php",
	"argument" : {
		"compile":[],
		"run" :[]
	}
 }
`
	var pl PayLoad
	err := json.Unmarshal([]byte(j), &pl)
	if err != nil {
		t.Error(err)
	}
	pl.runCode()
}

func TestRunJava(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"Hi.java",
		"content":"public class Hi {\n \tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Java: Hello, World!\");\n\t}\n}"
		}
	],
	"language" : "java",
	"argument" : {
		"compile":[],
		"run" :[]
	}
 }
`
	var pl PayLoad
	err := json.Unmarshal([]byte(j), &pl)
	if err != nil {
		t.Error(err)
	}
	pl.runCode()

}

func TestRunGo(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"hi.go",
		"content":"package main\n\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"Go: hello, world\")\n}"
		}
	],
	"language" : "go",
	"argument" : {
		"compile":[],
		"run" :[]
	}
 }
`
	var pl PayLoad
	err := json.Unmarshal([]byte(j), &pl)
	if err != nil {
		t.Error(err)
	}
	pl.runCode()
}

//func TestRunScala(t *testing.T) {
//
//	j := `
//{
//	"files" : [
//		{
//		"name":"Hi.scala",
//		"content":"object Hi {\n    def main(args: Array[String]) {\n        println(\"Hello, world!\")\n    }\n}"
//		}
//	],
//	"language" : "scala",
//	"argument" : {
//		"compile":[],
//		"run" :[]
//	}
//}
//`
//	var ar PayLoad
//	err := json.Unmarshal([]byte(j), &ar)
//	if err != nil {
//		t.Error(err)
//	}
//	ar.compileAndRun()
//}

func TestRunPerl(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"Hi.pl",
		"content":"#!/usr/bin/perl\n\nuse strict;\nuse warnings;\n\nprint \"Perl: Hello, World!\\n\";"
		}
	],
	"language" : "perl",
	"argument" : {
		"compile":[],
		"run" :[]
	}
}
`
	var pl PayLoad
	err := json.Unmarshal([]byte(j), &pl)
	if err != nil {
		t.Error(err)
	}
	pl.runCode()
}

func TestRunRuby(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"Hi.rb",
		"content":"#!/usr/bin/env ruby\n\nputs 'Ruby: Hello world'"
		}
	],
	"language" : "ruby",
	"argument" : {
		"compile":[],
		"run" :[]
	}
}
`
	var pl PayLoad
	err := json.Unmarshal([]byte(j), &pl)
	if err != nil {
		t.Error(err)
	}
	pl.runCode()
}
