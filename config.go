package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// DrkConfig is a struct representing a drk.yaml config file
type DrkConfig struct {
	DockerFile     string                            `yaml:"dockerFile,omitempty"`
	DockerImageDir string                            `yaml:"dockerDir,omitempty"`
	Commands       map[string]*DrkConfigBuildCommand `yaml:"commands"`
}

// DrkConfigBuildCommand is a struct representing a build command
type DrkConfigBuildCommand struct {
	Command        string   `yaml:"command,omitempty"`
	DockerFile     string   `yaml:"dockerFile,omitempty"`
	DockerImageDir string   `yaml:"dockerDir,omitempty"`
	NoDocker       bool     `yaml:"noDocker,omitempty"`
	Env            []string `yaml:"env,omitempty"`
}

// Gets a default config
func defaultConfig() DrkConfig {
	defaultCmd := defaultConfigBuildCommand("npm start")

	config := DrkConfig{
		"Dockerfile.build",
		"/code",
		map[string]*DrkConfigBuildCommand{
			"default": &defaultCmd,
		},
	}

	return config
}

func defaultConfigBuildCommand(command string) DrkConfigBuildCommand {
	return DrkConfigBuildCommand{command, "", "", false, []string{}}
}

// GetBuildCommand creates a build command coalesced from the specified command keyword specified by the user and the configuration file
func (config *DrkConfig) GetBuildCommand(command string) DrkConfigBuildCommand {
	commandData := *config.Commands["default"]
	if len(command) > 0 {
		if val, ok := config.Commands[command]; ok {
			commandData = *val
		} else {
			commandData.Command = commandData.Command + " " + command
		}
	}

	return commandData
}

// GetConfig gets a config based on the drk.yaml for the current directory
func GetConfig(cwd string) DrkConfig {
	configData, err := ioutil.ReadFile("drk.yaml")
	if err != nil {
		return defaultConfig()
	}

	config := defaultConfig()
	parseErr := yaml.Unmarshal(configData, &config)
	if parseErr != nil {
		panic(err)
	}

	for _, commandData := range config.Commands {
		// Coalesce each build command with defaults
		if commandData.DockerFile == "" {
			commandData.DockerFile = config.DockerFile
		}
		if commandData.DockerImageDir == "" {
			commandData.DockerImageDir = config.DockerImageDir
		}
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

	config := defaultConfig()

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(configPath, data, os.ModePerm)
	fmt.Println("Created drk.yaml with default values.")
}
