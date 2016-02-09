package cblib

import (
	"bufio"
	"code.google.com/p/gopass"
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
)

const (
	urlPrompt       = "Platform URL: "
	systemKeyPrompt = "System Key: "
	emailPrompt     = "Developer Email: "
	passwordPrompt  = "Password: "
)

func init() {
	flag.StringVar(&URL, "platform-url", "", "Clearblade platform url for target system")
	flag.StringVar(&SystemKey, "system-key", "", "System key for target system")
	flag.StringVar(&Email, "email", "", "Developer email for login")
	flag.StringVar(&Password, "password", "", "Developer password")
}

func getOneItem(prompt string, isASecret bool) string {
	if isASecret {
		thing, err := gopass.GetPass(prompt)
		if err != nil {
			fmt.Printf("Error getting password: %s\n", err.Error())
			os.Exit(1)
		}
		return thing
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s", prompt)
	thing, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading answer: %s\n", err.Error())
		os.Exit(1)
	}
	return strings.Trim(thing, "\n")
}

func fillInTheBlanks() {
	if URL == "" {
		URL = getOneItem(urlPrompt, false)
		cb.CB_ADDR = URL
	}
	if SystemKey == "" {
		SystemKey = getOneItem(systemKeyPrompt, false)
	}
	if Email == "" {
		Email = getOneItem(emailPrompt, false)
	}
	if Password == "" {
		Password = getOneItem(passwordPrompt, true)
	}
}

func GoToRepoRootDir() error {
	var err error
	whereIReallyAm, _ := os.Getwd()
	for {
		dirname, dirErr := os.Getwd()
		if dirErr != nil {
			return dirErr
		}
		if dirname == "/" {
			os.Chdir(whereIReallyAm) //  go back in case this err is ignored
			return fmt.Errorf(SpecialNoCBMetaError)
		}
		if MetaInfo, err = getDict(".cbmeta"); err != nil {
			if err = os.Chdir(".."); err != nil {
				return fmt.Errorf("Error changing directory: %s", err.Error())
			}
		} else {
			return nil
		}
	}
}

func Authorize() (*cb.DevClient, error) {
	/*
		err := GoToRepoRootDirAndGetMeta()
		if err == nil {
			// Auth using the .cbmeta file
			return &cb.DevClient{DevToken: DevToken}
		}

		if err.Error() != SpecialNoCBMetaError {
			fmt.Printf("Error trying to setup authorization: %s\n", err.Error())
			os.Exit(1)
		}
	*/
	if MetaInfo != nil {
		DevToken = MetaInfo["token"].(string)
		Email = MetaInfo["developerEmail"].(string)
		URL = MetaInfo["platformURL"].(string)
		cb.CB_ADDR = URL
		fmt.Printf("Using ClearBlade platform at '%s'\n", cb.CB_ADDR)
		return &cb.DevClient{DevToken: DevToken, Email: Email}, nil
	}
	// No cb meta file -- get url, syskey, email passwd
	fillInTheBlanks()
	fmt.Printf("Using ClearBlade platform at '%s'\n", cb.CB_ADDR)
	cli := cb.NewDevClient(Email, Password)
	if err := cli.Authenticate(); err != nil {
		return nil, err
	}
	return cli, nil
}
