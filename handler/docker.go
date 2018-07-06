package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi-back/internal/docker"
)

type result struct {
	UserResult *uResult `json:"userResult"`
	TaskError  string   `json:"taskError"`
}

type uResult struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	Exiterror string `json:"exiterror"`
}

// RunCode depended on language type and version
func RunCode(c *gin.Context) {
	language := c.Params.ByName("language")
	version := c.Params.ByName("version")
	// check language an version
	if !LanIsSupported(language) {
		c.String(http.StatusBadRequest, language+" is not support")
		return
	}
	if version == "" {
		version = VersionMap[language][0]
	} else {
		if !LVIsSupported(language, version) {
			c.String(http.StatusBadRequest, language+" "+version+" is not support")
			return
		}
	}

	var ar docker.PayLoad
	if err := c.ShouldBindJSON(&ar); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// use docker to run ric
	res, err := workerRun(ar, strings.ToLower(language), version)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func workerRun(ar docker.PayLoad, language, version string) (*result, error) {

	var w docker.Worker

	// load info to worker
	err := w.LoadInfo(&ar, language, V2Images(language, version))
	if err != nil {
		return nil, err
	}

	userResult, taskError, err := w.Run()
	if err != nil {
		return nil, err
	}

	var res result
	e := json.Unmarshal([]byte(userResult), &res.UserResult)
	if e != nil {
		fmt.Println(e)
	}

	res.TaskError = taskError

	return &res, nil
}
