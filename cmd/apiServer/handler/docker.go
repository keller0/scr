package handler

import (
	"bytes"
	"encoding/json"
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
	img := V2Images(strings.ToLower(language), version)
	var pl PayLoad
	if err := c.ShouldBindJSON(&pl); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"errNumber": responseErr["Payload not valid"]})
		return
	}

	pl.L = strings.ToLower(language)
	log.Info("request language: ", pl.L, " version: ", img)
	// use docker to run ric
	userResult, ricResp, err := runJob(pl, img)
	if err != nil {
		// api error
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
		return
	} else {
		if len(ricResp) > 0 {
			log.Error(ricResp)
			c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["Run code error"], "msg": ricResp})
		} else {
			var res result
			res.TaskError = ricResp
			e := json.Unmarshal([]byte(userResult), &res.UserResult)
			if e != nil {
				// decode user result error
				log.Error(e)
				c.JSON(http.StatusInternalServerError, gin.H{"errNumber": responseErr["Run code error"], "msg": e.Error()})
			} else {
				c.JSON(http.StatusOK, res)
			}

		}
		return
	}
}

func runJob(pl PayLoad, img string) (string, string, error) {

	var job docker.Job

	bs, err := json.Marshal(pl)
	if err != nil {
		return "", "", err
	}

	job.Payload = bytes.NewBuffer(bs)
	job.Image = img

	return job.Do()
}
