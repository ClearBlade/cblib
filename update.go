package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	updateCommand := &SubCommand{
		name:         "update",
		usage:        "update a specified resource from the remote system",
		needsAuth:    true,
		mustBeInRepo: true,
		run:          doUpdate,
	}
	updateCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to update")
	updateCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to update")
	updateCommand.flags.StringVar(&CollectionName, "collection", "", "Unique id of collection to update")
	updateCommand.flags.StringVar(&CollectionId, "collectionID", "", "Unique id of collection to update")
	updateCommand.flags.StringVar(&User, "user", "", "Unique id of user to update")
	updateCommand.flags.StringVar(&UserId, "userID", "", "Unique id of user to update")
	updateCommand.flags.StringVar(&RoleName, "role", "", "Name of role to update")
	updateCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to update")
	updateCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to update")
	AddCommand("update", updateCommand)
}

func checkUpdateArgsAndFlags(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("There are no arguments to the update command, only command line options\n")
	}
	return nil
}

func doUpdate(cmd *SubCommand, cli *cb.DevClient, args ...string) error {
	if err := checkUpdateArgsAndFlags(args); err != nil {
		return err
	}
	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}
	setRootDir(".")

	didSomething := false

	if ServiceName != "" {
		didSomething = true
		if err := pushOneService(systemInfo, cli); err != nil {
			return err
		}
	}

	if LibraryName != "" {
		didSomething = true
		if err := pushOneLibrary(systemInfo, cli); err != nil {
			return err
		}
	}

	if CollectionName != "" {
		didSomething = true
		if err := pushOneCollection(systemInfo, cli); err != nil {
			return err
		}
	}

	if CollectionId != "" {
		didSomething = true
		if err := pushOneCollectionById(systemInfo, cli); err != nil {
			return err
		}
	}

	if User != "" {
		didSomething = true
		if err := pushOneUser(systemInfo, cli); err != nil {
			return err
		}
	}

	if UserId != "" {
		didSomething = true
		if err := pushOneUserById(systemInfo, cli); err != nil {
			return err
		}
	}

	if RoleName != "" {
		didSomething = true
		if err := pushOneRole(systemInfo, cli); err != nil {
			return err
		}
	}

	if TriggerName != "" {
		didSomething = true
		if err := pushOneTrigger(systemInfo, cli); err != nil {
			return err
		}
	}

	if TimerName != "" {
		didSomething = true
		if err := pushOneTimer(systemInfo, cli); err != nil {
			return err
		}
	}

	if !didSomething {
		fmt.Printf("Nothing to update -- you must specify something to update (ie, -service=<svc_name>)\n")
	}

	return nil
}
