package main

import (
	"fmt"
	"os"
	"strings"

	"drk/config"
	"drk/docker"

	"github.com/jessevdk/go-flags"
)

// Structure representing the CLI arguments taken by this program
type cliArgs struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Help    bool `short:"h" long:"help" description:"Shows help"`
	New     bool `short:"n" long:"new" description:"Creates a template"`
	Args    struct {
		BuildCommand string   `description:"The alias for a build task to run in the docker build container"`
		Rest         []string `description:"Additional arguments" `
	} `positional-args:"yes" `
}

// Parses CLI arguments to a struct
func parseArgs(args []string) cliArgs {
	var opts cliArgs
	parser := flags.NewParser(&opts, flags.IgnoreUnknown|flags.PrintErrors|flags.HelpFlag)
	_, err := parser.ParseArgs(args[1:])
	if err != nil {
		parser.WriteManPage(os.Stderr)
		os.Exit(1)
	}

	if opts.Help {
		// BUG: Colors are not being written to console properly
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	return opts
}

// Given the config and CLI arguments, constructs the command to run inside the Docker build container
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

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	opts := parseArgs(os.Args)

	if opts.New {
		config.WriteConfig(cwd)
		os.Exit(0)
	}

	config := config.GetConfig(cwd)

	taskCommand := getTaskCommand(config, &opts)

	fmt.Println("drk: Running command: " + taskCommand)

	docker.RunCommandInBuildContainer(cwd, config.BuildImageDirectoryName, config.BuildImageName, taskCommand)
}
