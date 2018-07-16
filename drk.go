package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Args    struct {
		BuildCommand string
		Rest         []string
	} `positional-args:"yes" required:"no"`
}

type drkConfig struct {
	BuildFile               string
	BuildImageName          string
	BuildImageDirectoryName string
	BuildCommand            map[string]string
}

func defaultConfig(projectName string) drkConfig {
	return drkConfig{"Dockerfile.build", strings.ToLower(projectName + "build"), "/code", map[string]string{
		"default": "npm run",
	}}
}

func getConfig(projectName string) drkConfig {
	configData, err := ioutil.ReadFile("drk.json")
	if err != nil {
		return defaultConfig(projectName)
	}

	config := defaultConfig(projectName)
	parseErr := json.Unmarshal(configData, &config)
	if parseErr != nil {
		panic(err)
	}

	return config
}

// Builds the build-time docker container containing all build dependencies
func createBuildImage(name string) error {
	buildCmd := exec.Command("docker", "build", ".", "-t", name, "--file", "./Dockerfile.build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Run()

	return nil
}

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	parser := flags.NewParser(&opts, flags.IgnoreUnknown)
	_, argsErr := parser.ParseArgs(os.Args[1:])
	if argsErr != nil {
		panic(argsErr)
	}

	config := getConfig(path.Base(cwd))

	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	buildErr := createBuildImage(config.BuildImageName)
	if buildErr != nil {
		panic(buildErr)
	}

	buildCommand := config.BuildCommand["default"]
	if len(opts.Args.BuildCommand) > 0 {
		if val, ok := config.BuildCommand[opts.Args.BuildCommand]; ok {
			buildCommand = val
		} else {
			buildCommand = buildCommand + " " + opts.Args.BuildCommand
		}
	}

	command := buildCommand + " " + strings.Join(opts.Args.Rest, " ")

	fullCommand := []string{
		"docker",
		"run",
		"-w",
		config.BuildImageDirectoryName,
		"-v",
		cwd + ":" + config.BuildImageDirectoryName,
		config.BuildImageName,
		"/bin/sh",
		"-c",
		command}

	fmt.Println("drk: Running command: " + command)

	if opts.Verbose {
		fmt.Println("drk (verbose): " + strings.Join(fullCommand, " "))
	}

	// Here's the actual `syscall.Exec` call. If this call is
	// successful, the execution of our process will end
	// here and be replaced by the `/bin/ls -a -l -h`
	// process. If there is an error we'll get a return
	// value.
	execErr := syscall.Exec(
		binary,
		fullCommand,
		os.Environ())

	if execErr != nil {
		panic(execErr)
	}
}
