package config

import (
	"fmt"
	"io/ioutil"
	"os"
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
	imageName := ""
	if projectName != "" {
		imageName = strings.ToLower(projectName + "build")
	}
	return DrkConfig{"Dockerfile.build", imageName, "/code", map[string]string{
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

// WriteConfig writes a default config file template
func WriteConfig(cwd string) {
	configPath := path.Join(cwd, "drk.yaml")

	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("drk.yaml file already exists. Aborting.")
		os.Exit(1)
	}

	projectName := path.Base(cwd)
	data, _ := yaml.Marshal(defaultConfig(projectName))
	ioutil.WriteFile(configPath, data, os.ModeAppend)
	fmt.Println("Created drk.yaml with default values.")
}
