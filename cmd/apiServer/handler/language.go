package handler

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
			versions = append(versions, tVersion{v, "/api/v1/" + l + "/" + v})
		}
		all = append(all, tLanguage{l, &versions})
	}
	c.JSON(http.StatusOK, all)
}

// VersionsOfOne return all version of one language
func VersionsOfOne(c *gin.Context) {

	language := c.Params.ByName("language")
	if !LanIsSupported(language) {
		c.String(http.StatusNotFound, "%s is not supported", language)
	} else {
		var all []tLanguage
		var versions []tVersion
		for _, v := range VersionMap[language] {
			versions = append(versions, tVersion{v, "/api/v1/" + language + "/" + v})
		}
		all = append(all, tLanguage{language, &versions})
		c.JSON(http.StatusOK, all)
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
	"go": {"1.14"},
	"haskell": {
		"ghc-8.6",
	},
	"python": {
		"3.7",
		"2.7",
	},
	"php": {"7.4"},
	"java": {
		"14",
	},

	"perl":  {"5.28"},
	"perl6": {"latest"},
	"ruby":  {"2.8"},
	"rust":  {"latest"},
}

var imageMap = map[string]string{
	"bash-4.4": "gcc:8.3", // for now

	"c-gcc10": "gcc:10",

	"cpp-gcc10": "gcc:10",

	"php-7.2.5":  "php:7.2.5",
	"python-3.7": "python:3.7",
	"python-2.7": "python:2.7",

	"java-14": "openjdk:14",

	"go-1.14": "golang:1.14",

	"haskell-ghc-8.6": "haskell:8.6",
	"perl-5.28":       "perl:5.28",
	"perl6-latest":    "perl6",
	"ruby-2.8":        "ruby:2.8",
	"rust-latest":     "rust",
}

// V2Images return image name for one version of language
func V2Images(language, version string) string {

	return "yximages" + "/" + imageMap[language+"-"+version]

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
