package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keller0/scr/internal/docker"
	log "github.com/sirupsen/logrus"
)

type result struct {
	UserResult *uResult `json:"userResult"`
	TaskError  string   `json:"taskError"`
}

type uResult struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitError string `json:"exiterror"`
}

// PayLoad as stdin pass to ric container's stdin
type PayLoad struct {
	F []*oneFile `json:"files"`
	A *argument  `json:"argument"`
	I string     `json:"stdin"`
	L string     `json:"language"`
}

type argument struct {
	Compile []string `json:"compile"`
	Run     []string `json:"run"`
}

// file type
type oneFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
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

	var pl PayLoad
	if err := c.ShouldBindJSON(&pl); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Payload not valid"]})
		return
	}

	img := V2Images(strings.ToLower(language), version)
	// use docker to run ric
	res, err := runJob(pl, img)
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

func runJob(pl PayLoad, img string) (*result, error) {

	var job docker.Job

	bs, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	job.Payload = bytes.NewBuffer(bs)
	job.Image = img

	userResult, taskError, err := job.Do()
	if err != nil {
		return nil, err
	}

	var res result
	e := json.Unmarshal([]byte(userResult), &res.UserResult)
	if e != nil {
		log.Info(e)
		return nil, err
	}

	res.TaskError = taskError

	return &res, nil
}
