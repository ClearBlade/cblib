package cblib

import (
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

type SubCommand struct {
	name            string
	usage           string
	needsAuth       bool
	mustBeInRepo    bool
	mustNotBeInRepo bool
	flags           flag.FlagSet
	run             func(cmd *SubCommand, client *cb.DevClient, args ...string) error
	example			string
}

// centralized state of the repo
type RepoStatus struct {
	startedInRepo 		bool
	hasNavigatedToRepo	bool
	isInRepo			bool
	foundSystemDotJSON	bool
	foundCBMeta			bool
	hasValidToken		bool
	isHealthy			bool
	areURLSConfigured	bool
}

var (
	subCommands = map[string]*SubCommand{}
)

func (c *SubCommand) Execute(args []string) error {
	var client *cb.DevClient
	var err error
	c.flags.Parse(args)

	repoStatus := determineRepoStatus()
	fmt.Println("Debug1")
	fmt.Printf("%+v",repoStatus)
	if ! repoStatus.areURLSConfigured {
		platformURL, messagingURL := FormatURLs(URL, MsgURL)

		// WHAT DOES THIS DO????
		cb.CB_ADDR = platformURL
		cb.CB_MSG_ADDR = messagingURL
	}
	if c.mustBeInRepo && ! repoStatus.isInRepo {
		// Nav to repo to execute

		// TODO Should we change dirs?
		// cdToRoot()
		// repoStatus.isInRepo = true
	} else if c.mustNotBeInRepo && repoStatus.startedInRepo {
		// throw error
	}
	// Should we even change directories?
	// Should we export into currDir or new folder?
	/*
	if err = GoToRepoRootDir(); err != nil {
		if  {
			return fmt.Errorf("You must be in an initialized repo to run the '%s' command\n", c.name)
		}
		if err.Error() != SpecialNoCBMetaError {
			return err
		}
	} else if c.mustNotBeInRepo {
		return fmt.Errorf("You cannot run the '%s' command in an existing ClearBlade repository", c.name)
	}
	*/
	if repoStatus.foundCBMeta {
		// Check validity of dev token on disk

		// TODO Handle error, which isnt thrown right now
		valid := ValidateDevToken()
		repoStatus.hasValidToken = valid
		fmt.Println("Debug2")
		fmt.Printf("Valid? %v",valid)
		fmt.Printf("%+v",repoStatus)
	}

	if(repoStatus.hasValidToken){
		IngestCBMeta()
		cb.CB_ADDR = URL
		cb.CB_MSG_ADDR = MsgURL
		fmt.Printf("Using ClearBlade platform at '%s'\n", cb.CB_ADDR)
		fmt.Printf("Using ClearBlade messaging at '%s'\n", cb.CB_MSG_ADDR)
		client, _ = makeClientFromMetaInfo(MetaInfo)
	} else {
		client, err = PromptForAuthorize(nil)
		if err != nil {
			return err
		}
	}
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

func PrintRootHelp(){
	var usage = `
The cb-cli (ClearBlade CLI) provides methods for interacting with ClearBlade platform

Usage: cb-cli <command>

Commands:

`
	for cmd, _ := range subCommands {
		usage += fmt.Sprintf("\t%v\n",cmd)
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

func determineRepoStatus() RepoStatus{
	isInRepo := IsInRepo()
	foundCBMeta := FoundCBMeta()
	return RepoStatus{
		startedInRepo: 			isInRepo,
		hasNavigatedToRepo: 	false,
		isInRepo: 				isInRepo,
		foundSystemDotJSON: 	FoundSystemDotJSON(),
		foundCBMeta: 			foundCBMeta,
		hasValidToken: 			false,
		isHealthy: 				false,
		areURLSConfigured: 		URLSAreConfigured(),
	}


}
