package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExecCommand runs a command and wires up stdout/stderr
func ExecCommand(path string, args ...string) {
	cmd := exec.Command(path, args...)

	fmt.Println(cmd.Path + " " + strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd.Path + " " + strings.Join(cmd.Args, " "))
}

//RunCommandInBuildContainer runs the specified command inside the build image container
func RunCommandInBuildContainer(cwd string, dockerImageDirectoryName string, dockerFile string, command string) {
	_, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	projectName := filepath.Base(cwd)
	dockerImageName := projectName + strings.Replace(strings.ToLower(dockerFile), "dockerfile.", "", -1) + "image"

	ExecCommand("docker", "build", ".", "-t", dockerImageName, "--file", "./"+dockerFile)

	ExecCommand("docker",
		"run",
		"-w",
		dockerImageDirectoryName,
		"-v",
		cwd+":"+dockerImageDirectoryName,
		dockerImageName,
		"/bin/sh",
		"-c",
		command)
}
