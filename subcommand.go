package cblib

import (
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

type SubCommand struct {
	name  string
	usage string
	flags flag.FlagSet
	run   func(cmd *SubCommand, client *cb.DevClient, args ...string) error
}

var (
	subCommands = map[string]*SubCommand{}
)

func (c *SubCommand) Execute(client *cb.DevClient, args []string) error {
	c.flags.Usage = func() {
		helpFunc(c, c.name)
	}
	c.flags.Parse(args)

	return c.run(c, client, c.flags.Args()...)
}

func helpFunc(c *SubCommand, args ...string) {
	fmt.Printf("Usage: %s\n", c.usage)
}

func GetCommand(commandName string) (*SubCommand, error) {
	if theSubCommand, ok := subCommands[commandName]; ok {
		return theSubCommand, nil
	}
	return nil, fmt.Errorf("Subcommand %s not found", commandName)
}

func AddCommand(commandName string, stuff *SubCommand) {
	if _, alreadyExists := subCommands[commandName]; alreadyExists {
		fmt.Printf("Trying to add command %s but it already exists, ignoring...\n", commandName)
		return
	}
	subCommands[commandName] = stuff
}
