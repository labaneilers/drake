package docker

import (
	"os"
	"os/exec"
	"path"
	"syscall"
)

// Builds the build-time docker container containing all build dependencies
func createBuildImage(name string) error {
	buildCmd := exec.Command("docker", "build", ".", "-t", name, "--file", "./Dockerfile.build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Run()

	return nil
}

func RunCommandInBuildContainer(cwd string, buildImageDirectoryName string, buildImageName string, command string) {
	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	createBuildImage(path.Base(cwd))

	fullCommand := []string{
		"docker",
		"run",
		"-w",
		buildImageDirectoryName,
		"-v",
		cwd + ":" + buildImageDirectoryName,
		buildImageName,
		"/bin/sh",
		"-c",
		command}

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
