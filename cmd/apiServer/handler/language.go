package handler

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

type tVersion struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

type runner struct {
	Language string     `json:"name"`
	Versions []tVersion `json:"versions"`
}

var expectRunners []runner
var availableRunners []runner

func init() {

	for name, vs := range VersionMap {

		var versions []tVersion
		for _, v := range vs {
			versions = append(versions, tVersion{v, lv2Url(name, v)})
		}
		expectRunners = append(expectRunners, runner{name, versions})
	}
	// check available
	allLocalImages := getAllDockerImages()
	for _, r := range expectRunners {
		tmpRunner := &runner{Language: r.Language}

		for _, v := range r.Versions {
			tmpImg := V2Images(r.Language, v.Version)
			if containsString(allLocalImages, tmpImg) {
				tmpRunner.Versions = append(tmpRunner.Versions, v)
			}
		}
		if len(tmpRunner.Versions) > 0 {
			availableRunners = append(availableRunners, *tmpRunner)
		}
	}
}

// AllRunners return all supported languages and their versions
func AllRunners(c *gin.Context) {

	c.JSON(http.StatusOK, availableRunners)
}

// VersionsOfOne return all version of one language
func VersionsOfOne(c *gin.Context) {

	language := c.Params.ByName("language")
	if !LanIsSupported(language) {
		c.String(http.StatusNotFound, "%s is not supported", language)
	} else {

		var versions []tVersion
		for _, r := range availableRunners {
			if r.Language == language {
				versions = r.Versions
			}
		}

		c.JSON(http.StatusOK, versions)
	}
}

// VersionMap stands for all languages and versions
var VersionMap = map[string][]string{
	"bash": {"4.4"},
	"c": {
		"gcc10",
	},
	"cpp": {
		"gcc10",
	},
	"go": {"1.18"},
	"haskell": {
		"ghc-8.6",
	},
	"python": {
		"3.7",
		"2.7",
	},
	"php": {
		"7.4",
	},
	"java": {
		"14",
	},

	"perl":  {"5.28"},
	"perl6": {"latest"},
	"ruby":  {"2.7"},
	"rust":  {"latest"},
}

var imageMap = map[string]string{
	"bash-4.4": "gcc:10", // for now

	"c-gcc10": "gcc:10",

	"cpp-gcc10": "gcc:10",

	"php-7.4":    "php:7.4",
	"python-3.7": "python:3.7",
	"python-2.7": "python:2.7",

	"java-14": "openjdk:14",

	"go-1.18": "golang:1.18",

	"haskell-ghc-8.6": "haskell:8.6",
	"perl-5.28":       "perl:5.28",
	"perl6-latest":    "perl6",
	"ruby-2.7":        "ruby:2.7",
	"rust-latest":     "rust:latest",
}

// V2Images return image name for one version of language
func V2Images(language, version string) string {

	return "yximages" + "/" + imageMap[language+"-"+version]

}

func lv2Url(language, version string) string {
	return "/v1/" + language + "/" + version
}

// LanIsSupported check if the language is supported
func LanIsSupported(language string) bool {
	var supported bool
	for _, r := range availableRunners {
		if r.Language == language {
			supported = true
		}
	}
	return supported
}

// LVIsSupported check if the version of a language is supported
func LVIsSupported(lan, version string) bool {
	if !LanIsSupported(lan) {
		return false
	}
	var versions []tVersion
	for _, r := range availableRunners {
		if r.Language == lan {
			versions = r.Versions
			break
		}
	}

	for _, v := range versions {
		if v.Version == version {
			return true
		}
	}
	return false
}

func runnerDefaultVersion(language string) string {
	for _, r := range availableRunners {
		if r.Language == language {
			return r.Versions[0].Version
		}
	}
	return ""
}

func getAllDockerImages() []string {
	var imgS []string

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, i := range images {
		imgS = append(imgS, i.RepoTags...)
	}
	return imgS
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
