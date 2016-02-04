package cblib

import (
	"bufio"
	"code.google.com/p/gopass"
	"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	//"io/ioutil"
	"os"
	"os/user"
	"strings"
)

var (
	CommandLineEmail     bool
	DevToken             string
	AuthInfoFile         string
	MetaInfoFile         string
	SpecialNoCBMetaError = "No cbmeta file"
)

func init() {
	flag.StringVar(&AuthInfoFile, "authinfo", homedir()+"/.cbauth", "File in which you wish to store auth info")
	flag.StringVar(&MetaInfoFile, "metainfo", "./.cbmeta", "File in which you wish to store meta?! info")
}

func homedir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func Auth_prompt() (*cb.DevClient, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your email: ")
	email, _ := reader.ReadString('\n')
	email = strings.Trim(email, "\n")
	return AuthPromptPass(email)
}

func AuthPromptPass(email string) (*cb.DevClient, error) {
	pass, pass_err := gopass.GetPass(fmt.Sprintf("Enter password for '%s': ", email))
	if pass_err != nil {
		return nil, pass_err
	}
	return AuthUserPass(email, pass)
}

func AuthUserPass(email, password string) (*cb.DevClient, error) {
	cli := cb.NewDevClient(email, password)
	if err := cli.Authenticate(); err != nil {
		return nil, err
	}
	return cli, save_auth_info(AuthInfoFile, cli.DevToken)
}

func Auth(devToken string) (*cb.DevClient, error) {
	var cli *cb.DevClient
	if devToken != "" {
		cli = &cb.DevClient{
			DevToken: devToken,
		}
		return cli, nil
	}
	if _, err := os.Stat(AuthInfoFile); os.IsNotExist(err) {
		return Auth_prompt()
		/*
			email, pass, prompt_err := Auth_prompt()
			if prompt_err != nil {
				return nil, prompt_err
			}
			cli = cb.NewDevClient(email, pass)

			if err := cli.Authenticate(); err != nil {
				return nil, err
			} else {
				return cli, save_auth_info(AuthInfoFile, cli.DevToken)
			}
		*/
	} else {
		token, err := Load_auth_info(AuthInfoFile)
		if err != nil {
			return nil, err
		}
		cli = &cb.DevClient{
			DevToken: token,
		}
		//fmt.Println("Using developer token from " + homedir() + "/.cbauth")
		return cli, nil
	}
}

func save_auth_info(filename, token string) error {
	DevToken = token
	return nil
	//return ioutil.WriteFile(filename, []byte(token), 0600)
}

func Load_auth_info(filename string) (string, error) {
	return DevToken, nil
	/*
		if data, err := ioutil.ReadFile(filename); err != nil {
			return "", err
		} else {
			return string(data), nil
		}
	*/
}

func GoToRepoRootDir() error {
	var err error
	whereIReallyAm, _ := os.Getwd()
	MetaInfo = map[string]interface{}{}
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
