package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keller0/scr/cmd/apiServer/handler"
	"github.com/stretchr/testify/assert"
)

func TestRunCpp(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST("/:language", handler.RunCode)
	w := httptest.NewRecorder()

	buf := bytes.NewBufferString(cppHelloWorld)
	req, _ := http.NewRequest("POST", "/cpp", buf)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"userResult":{"stdout":"Hello, World! cpp","stderr":"","exiterror":""},"taskError":""}`, w.Body.String())
}

func TestRunGo(t *testing.T) {
	router := gin.New()
	router.POST("/:language", handler.RunCode)

	w := httptest.NewRecorder()

	buf := bytes.NewBufferString(goHelloWorld)
	req, _ := http.NewRequest("POST", "/go", buf)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"userResult":{"stdout":"Hello, World! go\n","stderr":"","exiterror":""},"taskError":""}`, w.Body.String())
}

const cppHelloWorld = `
{
	"files":[
		{
			"content":"#include <iostream>\n\nint main()\n{\n    std::cout << \"Hello, World! cpp\";\n}",
			"name":"main.cpp"
		}
	],
	"stdin":"",
	"argument":{
		"compile":[],
		"run":[]
	}
}
`
const goHelloWorld = `
{
	"files" : [
		{
		"name":"hi.go",
		"content":"package main\n\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"Hello, World! go\")\n}"
		}
	],
	"language" : "go",
	"argument" : {
		"compile":[],
		"run" :[]
	}
 }
`
