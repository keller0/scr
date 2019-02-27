package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keller0/yxi.io/internal/docker"
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
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Language not support"]})
		return
	}
	if version == "" {
		version = VersionMap[language][0]
	} else {
		if !LVIsSupported(language, version) {
			c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Language not support"]})
			return
		}
	}

	var ar docker.PayLoad
	if err := c.ShouldBindJSON(&ar); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Payload not valid"]})
		return
	}
	// use docker to run ric
	res, err := workerRun(ar, strings.ToLower(language), version)
	if err != nil {
		if err == docker.ErrWorkerTimeOut {
			c.JSON(http.StatusRequestTimeout, gin.H{"errNumber": responseErr["Time out"]})
		} else if err == docker.ErrTooMuchOutPut {
			c.JSON(http.StatusRequestTimeout, gin.H{"errNumber": responseErr["Too much output"]})
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"errNumber": responseErr["Run code error"],
					"msg":       err.Error(),
				})
		}
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
