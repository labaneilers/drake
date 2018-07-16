package main

import (
	"fmt"
	"os"
	"strings"

	"drk/config"
	"drk/docker"

	"github.com/jessevdk/go-flags"
)

type cliArgs struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Args    struct {
		BuildCommand string
		Rest         []string
	} `positional-args:"yes" required:"no"`
}

func getTaskCommand(config config.DrkConfig, opts *cliArgs) string {
	command := config.BuildCommand["default"]
	if len(opts.Args.BuildCommand) > 0 {
		if val, ok := config.BuildCommand[opts.Args.BuildCommand]; ok {
			command = val
		} else {
			command = command + " " + opts.Args.BuildCommand
		}
	}
	return command + " " + strings.Join(opts.Args.Rest, " ")
}

func parseArgs(args []string) cliArgs {
	var opts cliArgs
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)
	_, argsErr := parser.ParseArgs(args[1:])
	if argsErr != nil {
		panic(argsErr)
	}

	return opts
}

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	opts := parseArgs(os.Args)

	config := config.GetConfig(cwd)

	taskCommand := getTaskCommand(config, &opts)

	fmt.Println("drk: Running command: " + taskCommand)

	docker.RunCommandInBuildContainer(cwd, config.BuildImageDirectoryName, config.BuildImageName, taskCommand)
}
