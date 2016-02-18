package cblib

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/bgentry/speakeasy"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
)

const (
	urlPrompt       = "Platform URL"
	systemKeyPrompt = "System Key"
	emailPrompt     = "Developer Email"
	passwordPrompt  = "Password: "
)

func init() {
	flag.StringVar(&URL, "platform-url", "", "Clearblade platform url for target system")
	flag.StringVar(&SystemKey, "system-key", "", "System key for target system")
	flag.StringVar(&Email, "email", "", "Developer email for login")
	flag.StringVar(&Password, "password", "", "Developer password")
}

func getOneItem(prompt string, isASecret bool) string {
	reader := bufio.NewReader(os.Stdin)
	if isASecret {
		pw, err := speakeasy.Ask("Developer password: ")
		fmt.Printf("\n")
		if err != nil {
			fmt.Printf("Error getting password: %s\n", err.Error())
			os.Exit(1)
		}
		thing := string(pw)
		return strings.TrimSpace(thing)
	}
	fmt.Printf("%s: ", prompt)
	thing, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading answer: %s\n", err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(thing)
}

func buildPrompt(basicPrompt, defaultValue string) string {
	if defaultValue == "" {
		return basicPrompt
	}
	return fmt.Sprintf("%s (%s)", basicPrompt, defaultValue)
}

func getAnswer(entered, defaultValue string) string {
	if entered != "" {
		return entered
	}
	return defaultValue
}

func fillInTheBlanks(defaults *DefaultInfo) {
	var defaultUrl, defaultEmail, defaultSys string
	if defaults != nil {
		defaultUrl, defaultEmail, defaultSys = defaults.url, defaults.email, defaults.systemKey
	}
	if URL == "" {
		URL = getAnswer(getOneItem(buildPrompt(urlPrompt, defaultUrl), false), defaultUrl)
		cb.CB_ADDR = URL
	}
	if SystemKey == "" {
		SystemKey = getAnswer(getOneItem(buildPrompt(systemKeyPrompt, defaultSys), false), defaultSys)
	}
	if Email == "" {
		Email = getAnswer(getOneItem(buildPrompt(emailPrompt, defaultEmail), false), defaultEmail)
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

func Authorize(defaults *DefaultInfo) (*cb.DevClient, error) {
	if MetaInfo != nil {
		DevToken = MetaInfo["token"].(string)
		Email = MetaInfo["developerEmail"].(string)
		URL = MetaInfo["platformURL"].(string)
		cb.CB_ADDR = URL
		fmt.Printf("Using ClearBlade platform at '%s'\n", cb.CB_ADDR)
		return &cb.DevClient{DevToken: DevToken, Email: Email}, nil
	}
	// No cb meta file -- get url, syskey, email passwd
	fillInTheBlanks(defaults)
	fmt.Printf("Using ClearBlade platform at '%s'\n", cb.CB_ADDR)
	cli := cb.NewDevClient(Email, Password)
	if err := cli.Authenticate(); err != nil {
		return nil, err
	}
	return cli, nil
}
