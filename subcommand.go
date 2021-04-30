package cblib

import (
	"flag"
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type SubCommand struct {
	name      string
	usage     string
	needsAuth bool
	flags     flag.FlagSet
	run       func(cmd *SubCommand, client *cb.DevClient, args ...string) error
	example   string
}

var (
	subCommands = map[string]*SubCommand{}
)

func (c *SubCommand) Execute(args []string) error {
	var err error

	err = c.setup(args)
	if err != nil {
		return fmt.Errorf("Setup failed: %s", err)
	}

	err = c.beforeExecute(args)
	if err != nil {
		return fmt.Errorf("Before execute failed: %s", err)
	}

	err = c.execute(args)
	if err != nil {
		return err
	}

	err = c.afterExecute(args)
	if err != nil {
		return fmt.Errorf("After execute failed: %s", err)
	}

	return nil
}

func (c *SubCommand) setup(args []string) error {
	SetRootDir(".")
	return nil
}

func (c *SubCommand) beforeExecute(args []string) error {
	// TODO: read remotes
	// TODO: reconcile remotes from legacy if needed
	// TODO: activate current remote
	return nil
}

func (c *SubCommand) afterExecute(args []string) error {
	// TODO: persist remote changes
	return nil
}

func (c *SubCommand) execute(args []string) error {
	var client *cb.DevClient
	var err error
	c.flags.Parse(args)

	if URL != "" && MsgURL != "" {
		setupAddrs(URL, MsgURL)
	}

	// This is the most important part of initialization
	MetaInfo, _ = getCbMeta()

	if MetaInfo != nil {
		client, err = authorizeUsingGlobalMetaInfo()
		if err != nil {
			return fmt.Errorf("Authentication failed: %s", err)
		}

		if c.needsAuth {
			err := client.CheckAuth()
			if err != nil {
				return fmt.Errorf("Check authentication step failed. Please make sure that your token is valid. Hint: if you're inside an existing system you might want to run the 'target' command to re-authenticate. Error - %s", err.Error())
			}
		}
	} else if c.needsAuth {
		client, err = Authorize(nil)
		if err != nil {
			return err
		}
	}
	RootDirIsSet = false
	return c.run(c, client, c.flags.Args()...)

}

func PrintHelpFor(c *SubCommand, args ...string) {
	fmt.Printf("Usage: %s\n", c.usage)
	c.flags.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println(c.example)
}

func GetCommand(commandName string) (*SubCommand, error) {
	commandName = strings.ToLower(commandName)
	if theSubCommand, ok := subCommands[commandName]; ok {
		return theSubCommand, nil
	}
	return nil, fmt.Errorf("Subcommand %s not found", commandName)
}

func AddCommand(commandName string, stuff *SubCommand) {
	commandName = strings.ToLower(commandName)
	if _, alreadyExists := subCommands[commandName]; alreadyExists {
		fmt.Printf("Trying to add command %s but it already exists, ignoring...\n", commandName)
		return
	}
	subCommands[commandName] = stuff
}

func PrintRootHelp() {
	var usage = `
The cb-cli (ClearBlade CLI) provides methods for interacting with ClearBlade platform

Usage: cb-cli <command>

Commands:

`
	for cmd, _ := range subCommands {
		usage += fmt.Sprintf("\t%v\n", cmd)
	}
	usage += `

Examples:

	cb-cli init 							# inits your local workspace
	cb-cli export							# inits and exports your systems
	cb-cli push -service=Service1			# pushs an individual service up to Platform
	cb-cli pull -collection=Collection2	# pulls an individual collection to filesystem
	 `
	fmt.Println(usage)
}
