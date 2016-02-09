package cblib

import (
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

type SubCommand struct {
	name         string
	usage        string
	needsAuth    bool
	mustBeInRepo bool
	flags        flag.FlagSet
	run          func(cmd *SubCommand, client *cb.DevClient, args ...string) error
}

var (
	subCommands = map[string]*SubCommand{}
)

func (c *SubCommand) Execute( /*client *cb.DevClient,*/ args []string) error {
	var client *cb.DevClient
	var err error
	c.flags.Usage = func() {
		helpFunc(c, c.name)
	}
	c.flags.Parse(args)
	if URL != "" {
		cb.CB_ADDR = URL
	}
	if err = GoToRepoRootDir(); err != nil {
		if c.mustBeInRepo {
			return fmt.Errorf("You must be in an initialized repo to run the '%s' command\n", c.name)
		}
		if err.Error() != SpecialNoCBMetaError {
			return err
		}
	}
	if c.needsAuth {
		client, err = Authorize()
		if err != nil {
			return err
		}
	}

	return c.run(c, client, c.flags.Args()...)
}

func helpFunc(c *SubCommand, args ...string) {
	fmt.Printf("Usage: %s\n", c.usage)
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
