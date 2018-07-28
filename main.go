package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Version number (input at build time)
var Version string

// Structure representing the CLI arguments taken by this program
type cliArgs struct {
	Version     bool `short:"v" long:"version" description:"Get version"`
	Verbose     bool `long:"verbose" description:"Show verbose debug information"`
	Help        bool `short:"h" long:"help" description:"Shows help"`
	New         bool `short:"n" long:"new" description:"Creates a template"`
	Interactive bool `short:"i" long:"interactive" description:"Opens an interactive shell to the docker container for the command, but doesn't execute it"`
	Args        struct {
		BuildCommand string   `description:"The alias for a build task to run in the docker build container"`
		Rest         []string `description:"Additional arguments"`
	} `positional-args:"yes" `
}

// Parses CLI arguments to a struct
func parseArgs(args []string) cliArgs {
	var opts cliArgs
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)
	_, err := parser.ParseArgs(args[1:])
	if err != nil {
		parser.WriteManPage(os.Stderr)
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if opts.Help || opts.Args.BuildCommand == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	return opts
}

// Given the config and CLI arguments, constructs the command to run inside the Docker build container
func getTaskCommand(config DrkConfig, opts *cliArgs) DrkConfigBuildCommand {
	commandData := config.GetBuildCommand(opts.Args.BuildCommand)
	commandData.Command = commandData.Command + " " + strings.Join(opts.Args.Rest, " ")
	return commandData
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	opts := parseArgs(os.Args)

	if opts.New {
		WriteConfig(cwd)
		os.Exit(0)
	}

	config := GetConfig(cwd)

	taskCommand := getTaskCommand(config, &opts)

	fmt.Println("drk: Running command: " + taskCommand.Command)

	if taskCommand.NoDocker {
		splitCommand := strings.Split(taskCommand.Command, " ")
		ExecCommand(splitCommand[0], splitCommand[1:]...)
	} else {
		RunCommandInBuildContainer(
			cwd,
			taskCommand.DockerImageDir,
			taskCommand.DockerFile,
			taskCommand.Command,
			opts.Interactive,
			taskCommand.Env)
	}
}
