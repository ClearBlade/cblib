package cblib

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/maputil"
)

const (
	urlPrompt       = "Platform URL"
	msgurlPrompt    = "Messaging URL"
	systemKeyPrompt = "System Key"
	emailPrompt     = "Developer Email"
	passwordPrompt  = "Developer password: "
)

func initAuthFlags() {
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

// --------------------------------
// Prompt and fill functions
// --------------------------------
// TODO: Dangerous since they change global variables.

// PromptSet is used as a bitmask for configuring prompts.
type PromptSet uint8

const (
	// PromptSkipURL skips prompting the platform URL if used.
	PromptSkipURL PromptSet = 1 << iota
	// PromptSkipMsgURL skips prompting the messaging URL if used.
	PromptSkipMsgURL
	// PromptSkipEmail skips prompting the Email if used.
	PromptSkipEmail
	// PromptSkipSystemKey skips prompting the system key if used.
	PromptSkipSystemKey
	// PromptSkipPassword skips prompting the password if used.
	PromptSkipPassword
	// PromptAll prompts for all missing flags (a bit having a value of 1 means skip).
	PromptAll PromptSet = 0
)

// Has returns true if the given PromptSet has the desired flag.
func (p *PromptSet) Has(flag PromptSet) bool {
	return (*p)&flag != 0
}

func promptAndFillMissingURL(defaultURL string) bool {
	if URL == "" {
		URL = getAnswer(getOneItem(buildPrompt(urlPrompt, defaultURL), false), defaultURL)
		return true
	}
	return false
}

func promptAndFillMissingMsgURL(defaultMsgURL string) bool {
	if MsgURL == "" {
		MsgURL = getAnswer(getOneItem(buildPrompt(msgurlPrompt, defaultMsgURL), false), defaultMsgURL)
		return true
	}
	return false
}

func promptAndFillMissingURLAndMsgURL(defaultURL, defaultMsgURL string) (bool, bool) {
	promptedPlatformURL := promptAndFillMissingURL(defaultURL)
	if promptedPlatformURL {
		return true, promptAndFillMissingMsgURL(defaultMsgURL)
	}
	return false, false
}

func promptAndFillMissingEmail(defaultEmail string) bool {
	if Email == "" {
		Email = getAnswer(getOneItem(buildPrompt(emailPrompt, defaultEmail), false), defaultEmail)
		return true
	}
	return false
}

func promptAndFillMissingSystemKey(defaultSystemKey string) bool {
	if SystemKey == "" {
		SystemKey = getAnswer(getOneItem(buildPrompt(systemKeyPrompt, defaultSystemKey), false), defaultSystemKey)
		return true
	}
	return false
}

func promptAndFillMissingPassword() bool {
	if Password == "" {
		Password = getOneItem(passwordPrompt, true)
		return true
	}
	return false
}

func promptAndFillMissingAuth(defaults *DefaultInfo, promptSet PromptSet) {
	// var defaultURL, defaultMsgURL, defaultEmail, defaultSystemKey string
	var defaultURL, defaultEmail, defaultSystemKey string
	if defaults != nil {
		defaultURL = defaults.url
		// defaultMsgURL = defaults.msgUrl
		defaultEmail = defaults.email
		defaultSystemKey = defaults.systemKey
	}

	if !promptSet.Has(PromptSkipURL) {
		promptAndFillMissingURL(defaultURL)
	}

	// // TODO: messaging URL is optional since it can be derived from platform URL
	// // when not present
	// if !promptSet.Has(PromptSkipMsgURL) {
	// 	promptAndFillMissingMsgURL(defaultMsgURL)
	// }

	if !promptSet.Has(PromptSkipEmail) {
		promptAndFillMissingEmail(defaultEmail)
	}

	if !promptSet.Has(PromptSkipSystemKey) {
		promptAndFillMissingSystemKey(defaultSystemKey)
	}

	if !promptSet.Has(PromptSkipPassword) {
		promptAndFillMissingPassword()
	}
}

// --------------------------------
// Authorize
// --------------------------------
// Section focuses on authorizing and creating clearblade clients.

// authorizeUsingGlobalCLIFlags creates a new clearblade client by using the
// global flags passed to the CLI program.
func authorizeUsingGlobalCLIFlags() (*cb.DevClient, error) {
	return authorizeUsing(URL, MsgURL, Email, Password, DevToken)
}

// authorizeUsingGlobalMetaInfo creates a new clearblade client by using the
// the GLOBAL MetaInfo variable (cb meta must exists).
func authorizeUsingGlobalMetaInfo() (*cb.DevClient, error) {
	if MetaInfo == nil {
		return nil, fmt.Errorf("global meta info is nil")
	}
	return authorizeUsingMetaInfo(MetaInfo)
}

// authorizeUsingMetaInfo creates a new clearblade client by using the given
// meta info map.
func authorizeUsingMetaInfo(metaInfo map[string]interface{}) (*cb.DevClient, error) {

	if metaInfo == nil {
		return nil, fmt.Errorf("meta info should be non-nil")
	}

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

	return authorizeUsing(platformURL, messagingURL, email, "", token)
}

// authorizeUsing creates a new clearblade client using the given information.
// If a token is provided it takes precedence over the password.
func authorizeUsing(platformURL, messagingURL, email, password, token string) (*cb.DevClient, error) {

	platformURL, messagingURL, err := processURLs(platformURL, messagingURL)
	if err != nil {
		return nil, err
	}

	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	token = strings.TrimSpace(token)

	var cli *cb.DevClient

	if len(token) > 0 {
		cli = cb.NewDevClientWithTokenAndAddrs(platformURL, messagingURL, token, email)
		// TODO: commented out to preserve backward compatibility. Do we really need
		// to not check?
		// err = cli.CheckAuth()
		// if err != nil {
		// 	return nil, err
		// }

	} else if len(password) > 0 {
		cli = cb.NewDevClientWithAddrs(platformURL, messagingURL, email, password)
		err = verifyAuthentication(cli)
		if err != nil {
			return nil, err
		}

	} else {
		errmsg := fmt.Errorf("must provide either password or token")
		return nil, errmsg
	}

	return cli, nil
}

// verifyAuthentication verifies the given clearblade client, and prompts the
// user for a code if the given client requires two-factor auth.
func verifyAuthentication(cli *cb.DevClient) error {

	authResponse, err := cli.Authenticate()
	if err != nil {
		return err
	}

	info := authResponse.DevResponse

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
			return err
		}
	}

	return nil
}

// Authorize creates a new clearblade client by using the GLOBAL meta info, if
// it is not set, it will prompt the user for missing flags.
func Authorize(defaults *DefaultInfo) (*cb.DevClient, error) {

	var cli *cb.DevClient
	var err error

	if MetaInfo != nil {
		cli, err = authorizeUsingGlobalMetaInfo()

	} else {
		prompt := PromptAll
		if DevToken != "" {
			prompt |= PromptSkipEmail
			prompt |= PromptSkipPassword
		}
		promptAndFillMissingAuth(defaults, prompt)
		cli, err = authorizeUsingGlobalCLIFlags()
	}

	if err != nil {
		return nil, fmt.Errorf("Authorize failed: %s", err)
	}

	fmt.Printf("Using ClearBlade platform at '%s'\n", cli.HttpAddr)
	fmt.Printf("Using ClearBlade messaging at '%s'\n", cli.MqttAddr)

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
			// TODO: settings platform and messaging URL(s) using client rather
			// than globals from the clearblade Go SDK.
			// "platform_url":    cb.CB_ADDR,
			// "messaging_url":   cb.CB_MSG_ADDR,
			"platform_url":    client.HttpAddr,
			"messaging_url":   client.MqttAddr,
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
