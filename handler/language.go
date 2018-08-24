package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type tLanguage struct {
	Language string      `json:"name"`
	Versions *[]tVersion `json:"versions"`
}

type tVersion struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// AllVersion return all supported languages and their versions
func AllVersion(c *gin.Context) {
	var all []tLanguage
	for l, vs := range VersionMap {
		var versions []tVersion
		for _, v := range vs {
			versions = append(versions, tVersion{v, "/v1/" + l + "/" + v})
		}
		all = append(all, tLanguage{l, &versions})
	}
	c.JSON(http.StatusOK, all)
}

// VersionsOfOne return all version of one language
func VersionsOfOne(c *gin.Context) {

	language := c.Params.ByName("language")
	if !LanIsSupported(language) {
		c.String(http.StatusNotFound, "%s is not supported yet...", language)
	} else {
		var all []tLanguage
		var versions []tVersion
		for _, v := range VersionMap[language] {
			versions = append(versions, tVersion{v, "/v1/run/" + language + "/" + v})
		}
		all = append(all, tLanguage{language, &versions})
		c.JSON(http.StatusOK, all)
	}
}

// VersionMap stands for all languages and versions
var VersionMap = map[string][]string{
	"c": {
		"gcc8.1",
		"gcc7.3",
	},
	"cpp": {
		"g++8.1",
		"g++7.3",
	},
	"python":  {
		"2.7",
		"3.5",
	},
	"bash":    {"4.4"},
	"php":     {"7.2.5"},
	"java":    {"openjdk-8"},
	"go":      {"1.8", "1.10"},
	"haskell": {"ghc-8"},
	"perl":    {"5.28"},
	"ruby":    {"2.5"},
	"rust":    {"1.27"},
}

var imageMap = map[string]string{
	"bash-4.4":       "gcc:8.1", // for now
	"c-gcc8.1":       "gcc:8.1",
	"c-gcc7.3":       "gcc:7.3",
	"cpp-g++8.1":     "gcc:8.1",
	"cpp-g++7.3":     "gcc:7.3",
	"php-7.2.5":      "php:7.2.5",
	"python-3.5":     "python:3.5",
	"python-2.7":     "python:2.7-slim",
	"java-openjdk-8": "java:8",
	"go-1.8":         "golang:1.8",
	"go-1.10":        "golang:1.10",
	"haskell-ghc-8":  "haskell:8",
	"perl-5.28":      "perl:5.28",
	"ruby-2.5":       "ruby:2.5",
	"rust-1.27":      "rust:1.27",
}

// V2Images return image name for one version of language
func V2Images(language, version string) string {

	return "keller0" + "/" + imageMap[language+"-"+version]

}

// LanIsSupported check if the language is supported
func LanIsSupported(language string) bool {
	_, got := VersionMap[language]
	return got
}

// LVIsSupported check if the version of a language is supported
func LVIsSupported(lan, version string) bool {
	if !LanIsSupported(lan) {
		return false
	}
	vs := VersionMap[lan]
	for _, v := range vs {
		if v == version {
			return true
		}
	}
	return false
}
