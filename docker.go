package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// ExecCommand runs a command and wires up stdout/stderr
func ExecCommand(path string, args ...string) {
	cmd := exec.Command(path, args...)

	color.Cyan(strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

//RunCommandInBuildContainer runs the specified command inside the build image container
func RunCommandInBuildContainer(cwd string, taskCommand DrkConfigBuildCommand, interactive bool) {
	_, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	projectName := filepath.Base(cwd)
	dockerImageName := projectName + strings.Replace(strings.ToLower(taskCommand.DockerFile), "dockerfile.", "", -1) + "image"

	ExecCommand("docker", "build", ".", "-t", dockerImageName, "--file", "./"+taskCommand.DockerFile)

	args := []string{
		"run",
		"-w",
		taskCommand.DockerImageDir,
		"-v",
		cwd + ":" + taskCommand.DockerImageDir}

	for _, portMapping := range taskCommand.Ports {
		args = append(args, "-p")
		args = append(args, portMapping)
	}

	for _, envVar := range taskCommand.Env {
		args = append(args, "-e")
		args = append(args, envVar)
	}

	if interactive {
		args = append(args, "-it")
	}

	args = append(args, dockerImageName)
	args = append(args, "/bin/sh")

	if !interactive {
		args = append(args, "-c")
		args = append(args, taskCommand.Command)
	}

	ExecCommand("docker", args...)
}
