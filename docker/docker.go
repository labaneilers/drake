package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

// Builds the build-time docker container containing all build dependencies
func createBuildImage(buildFile string, name string) error {
	buildCmd := exec.Command("docker", "build", ".", "-t", name, "--file", "./"+buildFile)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Run()

	return nil
}

//RunCommandInBuildContainer runs the specified command inside the build image container
func RunCommandInBuildContainer(cwd string, dockerImageDirectoryName string, dockerFile string, command string) {
	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	projectName := path.Base(cwd)
	dockerImageName := projectName + strings.Replace(strings.ToLower(dockerFile), "dockerfile.", "", -1) + "image"

	createBuildImage(dockerFile, dockerImageName)

	fullCommand := []string{
		"docker",
		"run",
		"-w",
		dockerImageDirectoryName,
		"-v",
		cwd + ":" + dockerImageDirectoryName,
		dockerImageName,
		"/bin/sh",
		"-c",
		command}

	fmt.Println(strings.Join(fullCommand, " "))

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
