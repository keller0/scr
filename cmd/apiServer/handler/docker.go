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
	ExitError string `json:"exitError"`
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
	language := strings.ToLower(c.Params.ByName("language"))
	version := strings.ToLower(c.Params.ByName("version"))

	// check language and version
	if !LanIsSupported(language) {
		c.JSON(http.StatusBadRequest, LanguageNotSupported)
		return
	}
	if version == "" {
		version = runnerDefaultVersion(language)
	}

	if !LVIsSupported(language, version) {
		c.JSON(http.StatusBadRequest, LanguageNotSupported)
		return
	}

	img := V2Images(language, version)
	var pl PayLoad
	if err := c.ShouldBindJSON(&pl); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, PayloadNotValid)
		return
	}

	pl.L = strings.ToLower(language)
	log.Info("request language: ", pl.L, " version: ", img)
	// use docker to run ric
	userResult, ricResp, err := runJob(pl, img)
	if err != nil {
		// api error
		if err == docker.ErrWorkerTimeOut {
			c.JSON(http.StatusRequestTimeout, TimeOutErr)
		} else if err == docker.ErrTooMuchOutPut {
			c.JSON(http.StatusRequestTimeout, TooMuchOutPutErr)
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"code": RunCodeErr.Code,
					"msg":  err.Error(),
				})
		}
		return
	} else {
		if len(ricResp) > 0 {
			log.Error(ricResp)
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"code": RunCodeErr.Code,
					"msg":  ricResp,
				})
		} else {
			var res result
			res.TaskError = ricResp
			e := json.Unmarshal([]byte(userResult), &res.UserResult)
			if e != nil {
				// decode user result error
				log.Error(e)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": RunCodeErr.Code,
					"msg":  ricResp,
				})
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
