package cblib

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/maputil"
)

const (
	urlPrompt       = "Platform URL"
	msgurlPrompt    = "Messaging URL"
	systemKeyPrompt = "System Key"
	emailPrompt     = "Developer Email"
	passwordPrompt  = "Developer password: "
)

func init() {
	flag.StringVar(&URL, "platform-url", "", "Clearblade platform url for target system")
	flag.StringVar(&MsgURL, "messaging-url", "", "Clearblade messaging url for target system")
	flag.StringVar(&SystemKey, "system-key", "", "System key for target system")
	flag.StringVar(&Email, "email", "", "Developer email for login")
	flag.StringVar(&Password, "password", "", "Developer password")
}

func getOneItem(prompt string, isASecret bool) string {
	reader := bufio.NewReader(os.Stdin)
	if isASecret {
		pw, err := speakeasy.Ask(prompt)
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

// fillInTheBlanks will prompt the user for GLOBALS that are not
// provided via flags.
func fillInTheBlanks(defaults *DefaultInfo) {
	var defaultURL, defaultMsgURL, defaultEmail, defaultSys string

	if defaults != nil {
		defaultURL = defaults.url
		defaultMsgURL = defaults.msgUrl
		defaultEmail = defaults.email
		defaultSys = defaults.systemKey
	}

	if URL == "" {
		URL = getAnswer(getOneItem(buildPrompt(urlPrompt, defaultURL), false), defaultURL)
	}

	if MsgURL == "" {
		MsgURL = getAnswer(getOneItem(buildPrompt(msgurlPrompt, defaultMsgURL), false), defaultMsgURL)
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

	setupAddrs(URL, MsgURL)
}

// newClientFromMetaInfo creates a new clearblade client from the given meta
// info. The meta info should contain the following fields:
// - "token"
// - "developerEmail" or "developer_email"
// - "platformURL" or "platform_url"
// - "messagingURL" or "platform_url"
func newClientFromMetaInfo(metaInfo map[string]interface{}) (*cb.DevClient, error) {

	// Mix and match old schema vs new schema
	// LookupString defaults to empty string when it doesn't find any of the
	// given keys.
	token, tokenOk := maputil.LookupString(metaInfo, "token")
	email, emailOk := maputil.LookupString(metaInfo, "developerEmail", "developer_email")
	platformURL, platformURLOk := maputil.LookupString(metaInfo, "platformURL", "platform_url")
	messagingURL, _ := maputil.LookupString(metaInfo, "messagingURL", "messaging_url")

	if !tokenOk {
		return nil, fmt.Errorf("missing token from meta info")
	}

	if !emailOk {
		return nil, fmt.Errorf("missing email from meta info")
	}

	if !platformURLOk {
		return nil, fmt.Errorf("missing platform url from meta info")
	}

	// WARNING: changes globals in clearblade SDK
	setupAddrs(platformURL, messagingURL)

	return cb.NewDevClientWithToken(token, email), nil
}

// newClientFromGlobalMetaInfo is similar to newClientFromMetaInfo but uses
// the GLOBAL MetaInfo instead of having to pass your own meta info. Use with
// caution.
func newClientFromGlobalMetaInfo() (*cb.DevClient, error) {
	return newClientFromMetaInfo(MetaInfo)
}

// Authorize creates a new clearblade client by using the GLOBAL meta info, if
// it is not set, it will prompt the user for missing fields.
func Authorize(defaults *DefaultInfo) (*cb.DevClient, error) {

	if MetaInfo != nil {
		return newClientFromGlobalMetaInfo()
	}

	// No cb meta file -- get url, syskey, email passwd
	fillInTheBlanks(defaults)

	fmt.Printf("Using ClearBlade platform at '%s'\n", cb.CB_ADDR)
	fmt.Printf("Using ClearBlade messaging at '%s'\n", cb.CB_MSG_ADDR)

	cli := cb.NewDevClient(Email, Password)
	authResp, err := cli.Authenticate()
	if err != nil {
		fmt.Printf("Authenticate failed: %s\n", err)
		return nil, err
	}
	info := authResp.DevResponse
	if info.IsTwoFactor {
		prompt := getPromptBasedOnTwoFactorMethod(info.TwoFactorMethod)
		code := getAnswer(getOneItem(buildPrompt(prompt, ""), false), "")
		err := cli.VerifyAuthentication(cb.VerifyAuthenticationParams{
			Code:            code,
			TwoFactorMethod: info.TwoFactorMethod,
			OtpID:           info.OtpID,
			OtpIssued:       info.OtpIssued,
		})
		if err != nil {
			fmt.Printf("Authentication verification failed: %s\n", err.Error())
			return nil, err
		}
	}
	return cli, nil
}

func getPromptBasedOnTwoFactorMethod(method string) string {
	switch method {
	case "email":
		return "Please enter the code sent to your email inbox"
	case "sms":
		return "Please enter the code sent to your device"
	case "email_sms":
		return "Please enter the code sent to your email inbox and device"
	}
	return "Please enter the code"
}

func checkIfTokenHasExpired(client *cb.DevClient, systemKey string) (*cb.DevClient, error) {
	err := client.CheckAuth()
	if err != nil {
		fmt.Printf("Token has probably expired. Please enter details for authentication again...\n")
		MetaInfo = nil
		client, _ = Authorize(nil)
		metaStuff := map[string]interface{}{
			"platform_url":    cb.CB_ADDR,
			"messaging_url":   cb.CB_MSG_ADDR,
			"developer_email": Email,
			"token":           client.DevToken,
		}
		if err = storeCBMeta(metaStuff); err != nil {
			return nil, err
		}
		return client, nil
	}
	return client, nil
}
