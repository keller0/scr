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
		  "content" : "#include <stdio.h>\n\nint main(void) {\n    printf(\"Hello, World!\\n\");\n    return 0;\n}",
		  "name" : "main.c"
	   }
	],
	"language" : "c",
	"argument" : {
		"compile":[
			"gcc"
		],
		"run" :[
		
		]
	}
 }
`
	var ar PayLoad
	err := json.Unmarshal([]byte(j), &ar)
	if err != nil {
		t.Error(err)
	}
	ar.compileAndRun()

}

func TestRunPHP(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"main.php",
		"content":"<?php\n    echo \"Hello, World!\"; exit(10);"
		}
	],
	"language" : "php",
	"argument" : {
		"compile":[],
		"run" :[]
	}
 }
`
	var ar PayLoad
	err := json.Unmarshal([]byte(j), &ar)
	if err != nil {
		t.Error(err)
	}
	ar.Run()
}

func TestRunJava(t *testing.T) {

	j := `
{
	"files" : [
		{
		"name":"Hi.java",
		"content":"public class Hi {\n \tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello, World!\");\n\t}\n}"
		}
	],
	"language" : "java",
	"argument" : {
		"compile":[],
		"run" :[]
	}
 }
`
	var ar PayLoad
	err := json.Unmarshal([]byte(j), &ar)
	if err != nil {
		t.Error(err)
	}
	ar.compileAndRun()

}
