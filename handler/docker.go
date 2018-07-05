package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

// Arun - a run request from api client
type Arun struct {
	Files    []*oneFile `json:"files"`
	Argument *argument  `json:"argument"`
	Stdin    string     `json:"stdin"`
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

	var ar Arun
	if err := c.ShouldBindJSON(&ar); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// use docker to run ric
	res, err := ar.workerRun(strings.ToLower(language), version)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func (ar *Arun) workerRun(language, version string) (*result, error) {

	var w Worker
	var err error
	// load info to worker
	err = w.loadInfo(ar, language, V2Images(language, version))
	if err != nil {
		return nil, err
	}

	containerJSON, err := w.createContainer()
	defer func() {
		err = w.cli.ContainerRemove(w.ctx, w.tmpID, types.ContainerRemoveOptions{})
		fmt.Println("Container", w.tmpID, "removed")
		if err != nil {
			fmt.Println("failed to remove container ", w.tmpID)
		}
	}()

	if err != nil {
		return nil, err
	}
	w.containerID = containerJSON.ID
	err = w.attachContainer()
	if err != nil && w.ricErr.Len() == 0 {
		return nil, err
	}

	userResult := w.ricOut.String()
	taskError := w.ricErr.String()

	var res result

	e := json.Unmarshal([]byte(userResult), &res.UserResult)
	if e != nil {
		fmt.Println(e)
	}

	res.TaskError = taskError

	return &res, nil
}
