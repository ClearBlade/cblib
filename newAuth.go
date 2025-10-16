package cblib

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/chromedp/chromedp"
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/maputil"
)

const (
	urlPrompt          = "Platform URL"
	msgurlPrompt       = "Messaging URL"
	systemKeyPrompt    = "System Key"
	browserLoginPrompt = "Login using Browser? (n|Y - Only Google Chrome supported.)"
	emailPrompt        = "Developer Email"
	passwordPrompt     = "Developer Password (will be hidden): "
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

func isBlankOrNull(param string) bool {
	return param == "" || param == "null"
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

func promptIfSkipBrowserLogin() bool {
	browserLogin := getAnswer(getOneItem(buildPrompt(browserLoginPrompt, ""), false), "Y")
	trimLowerBrowserLogin := strings.ToLower(strings.TrimSpace(browserLogin))
	return trimLowerBrowserLogin == "no" || trimLowerBrowserLogin == "n"
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

func retrieveTokenFromLocalStorageChrome(url string) (string, error) {
	// Retain the long grace period for maximum chance of natural cleanup
	// 3 seconds was chosen because with shorter times it seemed that the token
	// was NOT persisting in Local Storage. Strangely the CURRENT login WOULD work
	// (i.e. the token must have been found in Local Storage), but upon SUBSEQUENT
	// login attempts the token was NOT found. Not sure if there is a process that
	// needs to complete to make sure the Local Storage is "locked". With 3 seconds
	// I found the token was ALWAYS present upon subsequent login attempts.
	shutdownGracePeriod := 3 * time.Second
	tempDir := os.TempDir()
	tempDataDir := filepath.Join(tempDir, "cb-cli-chprof")
	const customProfileDir = "Default"

	// Pre-flight cleanup for common lock files
	singletonLock := filepath.Join(tempDataDir, "SingletonLock")
	if _, err := os.Stat(singletonLock); err == nil {
		os.Remove(singletonLock)
	}

	// Create a browser context with HEADLESS set to FALSE
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.UserDataDir(tempDataDir),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("profile-directory", customProfileDir),
	)

	// Custom deferred function for graceful shutdown
	defer func() {

		// Signal a shutdown to the ExecAllocator
		cancel()

		// Local State (at the root of the user data directory)
		// Writing empty object to prevent "Chrome didn't shut down correctly" dialog.
		// Getting rid of that dialog leads to a better user experience.
		// Also "corrupting" this file is not too risky since the profile is in a
		// temporary directory separate from the main profile AND a new browser instance
		// is launched with every login attempt.
		localStateFile := filepath.Join(tempDataDir, "Local State")
		if err := os.WriteFile(localStateFile, []byte("{}"), 0644); err != nil {
			fmt.Printf("Warning: Failed to clear Local State file: %v\n", err)
		}

		// Preferences file (inside the custom profile directory)
		// Writing empty object to prevent "Chrome didn't shut down correctly" dialog.
		// Getting rid of that dialog leads to a better user experience.
		// Also "corrupting" this file is not too risky since the profile is in a
		// temporary directory separate from the main profile AND a new browser instance
		// is launched with every login attempt.
		preferencesFile := filepath.Join(tempDataDir, customProfileDir, "Preferences")
		if err := os.WriteFile(preferencesFile, []byte("{}"), 0644); err != nil {
			fmt.Printf("Warning: Failed to clear Preferences file: %v\n", err)
		}
	}()

	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()

	// Set long enough timeout for the entire manual login process
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Minute)
	defer timeoutCancel()

	var token string
	var loginURL = url + "/login"
	const tokenKey = "ngStorage-cb_platform_dev_token" // This is the key with which the browser stores the ClearBlade auth token in its Local Storage

	// JS function to read the token
	jsGetToken := fmt.Sprintf(`localStorage.getItem("%s");`, tokenKey)

	// Launch the browser and navigate
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(loginURL),
	)
	if err != nil {
		return "", fmt.Errorf("failed to launch Chrome browser (ensure Chrome is installed): %w", err)
	}

	browserLoginStarted := false
	tokenRetrieved := false

	// Poll the local storage until the token is found
	for {
		select {
		case <-timeoutCtx.Done():
			// The overall 5-minute timeout was hit
			return "", fmt.Errorf("login timeout reached before token was set")
		default:
			// Execute JS to read the local storage item
			err := chromedp.Run(timeoutCtx,
				chromedp.Evaluate(jsGetToken, &token),
			)
			if err != nil {
				fmt.Printf("Error during token check: %v. Retrying...\n", err)
			}

			if !isBlankOrNull(token) {
				tokenRetrieved = true
				fmt.Printf("Logged into %s.\n", URL)

				if browserLoginStarted {
					fmt.Printf("Complete any activity in the browser. Then click ENTER to close the browser and continue.\n")
					// Wait for user input
					reader := bufio.NewReader(os.Stdin)
					_, _ = reader.ReadString('\n')

					fmt.Printf("Closing browser. Please wait.\n")
					// A few seconds wait to let Chrome complete internal processes
					time.Sleep(shutdownGracePeriod)
				}

				return token, nil

			}

			if !tokenRetrieved && !browserLoginStarted {
				fmt.Printf("Login manually in the browser.\n")
				browserLoginStarted = true
			}

			// Wait 1 second before checking again
			time.Sleep(1 * time.Second)
		}
	}
}

func promptAndFillMissingAuth(defaults *DefaultInfo, promptSet PromptSet) {
	// var defaultURL, defaultMsgURL, defaultEmail, defaultSystemKey string
	var defaultURL, defaultEmail, defaultSystemKey, token string
	var err error
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

	if isBlankOrNull(DevToken) && (isBlankOrNull(Email) || isBlankOrNull(Password)) {
		SkipBrowserLogin := promptIfSkipBrowserLogin()

		if SkipBrowserLogin {
			if !promptSet.Has(PromptSkipEmail) {
				promptAndFillMissingEmail(defaultEmail)
			}

			if !promptSet.Has(PromptSkipPassword) {
				promptAndFillMissingPassword()
			}
			// Browser login was never initiated, continue to prompt for system key
		} else {
			// Browser login was initiated
			token, err = retrieveTokenFromLocalStorageChrome(URL)

			if err == nil {
				DevToken = strings.Trim(token, "\"") // remove double-quotes from returned token
				// Browser login succeeded, continue to prompt for system key
			} else {
				// Browser login failed, abort and don't prompt for system key
				fmt.Printf("Browser login was not completed: %v\n", err)
				return // Exit the function early without prompting for system key
			}
		}
	}

	if !promptSet.Has(PromptSkipSystemKey) {
		promptAndFillMissingSystemKey(defaultSystemKey)
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
		errmsg := fmt.Errorf("must either provide password / token or login using browser")
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
