package config

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/ghodss/yaml"
)

// DrkConfig is a struct representing a drk.yaml config file
type DrkConfig struct {
	BuildFile               string
	BuildImageName          string
	BuildImageDirectoryName string
	BuildCommand            map[string]string
}

// Gets a default config
func defaultConfig(projectName string) DrkConfig {
	return DrkConfig{"Dockerfile.build", strings.ToLower(projectName + "build"), "/code", map[string]string{
		"default": "npm run",
	}}
}

// GetConfig gets a config based on the drk.yaml for the current directory
func GetConfig(cwd string) DrkConfig {
	projectName := path.Base(cwd)
	configData, err := ioutil.ReadFile("drk.yaml")
	if err != nil {
		return defaultConfig(projectName)
	}

	config := defaultConfig(projectName)
	parseErr := yaml.Unmarshal(configData, &config)
	if parseErr != nil {
		panic(err)
	}

	return config
}
