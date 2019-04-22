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
		"gcc8.3",
		"gcc7.4",
		"gcc6.5",
		"gcc5.5",
	},
	"cpp": {
		"gcc8.3",
		"gcc7.4",
		"gcc6.5",
		"gcc5.5",
	},
	"go": {"1.12", "1.11"},
	"haskell": {
		"ghc-8.6",
	},
	"python": {
		"3.7",
		"2.7",
	},
	"php": {"7.2.5"},
	"java": {
		"13",
		"12",
		"11",
		"8",
	},

	"perl":  {"5.28"},
	"perl6": {"latest"},
	"ruby":  {"2.6"},
	"rust":  {"latest"},
}

var imageMap = map[string]string{
	"bash-4.4": "gcc:8.3", // for now

	"c-gcc8.3": "gcc:8.3",
	"c-gcc7.4": "gcc:7.4",
	"c-gcc6.5": "gcc:6.5",
	"c-gcc5.5": "gcc:5.5",

	"cpp-gcc8.3": "gcc:8.3",
	"cpp-gcc7.4": "gcc:7.4",
	"cpp-gcc6.5": "gcc:6.5",
	"cpp-gcc5.5": "gcc:5.5",

	"php-7.2.5":  "php:7.2.5",
	"python-3.7": "python:3.7",
	"python-2.7": "python:2.7",

	"java-13": "openjdk:13",
	"java-12": "openjdk:12",
	"java-11": "openjdk:11",
	"java-8":  "openjdk:8",

	"go-1.11": "golang:1.11",
	"go-1.12": "golang:1.12",

	"haskell-ghc-8.6": "haskell:8.6",
	"perl-5.28":       "perl:5.28",
	"perl6-latest":    "perl6",
	"ruby-2.6":        "ruby:2.6",
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
