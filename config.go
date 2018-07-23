package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Config stores configuration of server
type Config struct {
	ListenAddr    string
	ReadTimeout   int
	WriteTimeout  int
	ClientTimeout int
	BaseProxyURL  string
	Paths         []Path
}

// Path stores details of an individual path
type Path struct {
	Path    string
	Handler string
}

var config Config

func loadConfig() {
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &config)
}

func getPathFromConfig(path string) Path {
	for i := range config.Paths {
		if strings.HasPrefix(path, config.Paths[i].Path) {
			// Found!
			return config.Paths[i]
		}
	}
	ret := Path{}
	ret.Handler = "proxy"
	return ret
}
