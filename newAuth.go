package cblib

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
	browserLoginPrompt = "Login using Browser?"
	emailPrompt        = "Developer Email"
	passwordPrompt     = "Developer Password"
	callbackPort       = ":8080"
)

// AuthRequest mirrors the payload of the ClearBlade login POST request.
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse mirrors the successful response body from the ClearBlade login.
type AuthResponse struct {
	DevToken     string `json:"dev_token"`
	ExpiresAt    int64  `json:"expires_at"`
	IsTwoFactor  bool   `json:"is_two_factor"`
	RefreshToken string `json:"refresh_token"`
}

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

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin": // macOS
		cmd = "open"
		args = []string{url}
	default: // Linux and others
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
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

func promptAndFillMissingBrowserLogin(defaultBrowserLogin string) string {
	BrowserLogin := getAnswer(getOneItem(buildPrompt(browserLoginPrompt, defaultBrowserLogin), false), defaultBrowserLogin)
	TrimLowerBrowserLogin := strings.ToLower(strings.TrimSpace(BrowserLogin))
	if TrimLowerBrowserLogin == "no" || TrimLowerBrowserLogin == "n" {
		return "n"
	}
	return "y"
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

func retrieveTokenViaRedirect(url string) {
	var finalToken string
	var choice string
	// A channel to receive the token and an error.
	tokenChan := make(chan string)
	errChan := make(chan error)
	fmt.Printf("Akash URL: %s\n", URL)
	loginURL := url + "/login"
	redirectURI := "http://localhost" + callbackPort + "/callback"

	// Set up a local web server to handle the redirect (for future OIDC)
	server := &http.Server{Addr: callbackPort}

	// The /callback handler remains, but we know it won't run for ClearBlade's current page.
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// THIS LOGIC IS FOR FUTURE OIDC PROVIDERS.
		token := r.URL.Query().Get("token")
		if token == "" {
			errChan <- fmt.Errorf("no token found in callback URL")
			fmt.Fprintf(w, "<html><body><h1>Login Failed</h1><p>No token received.</p></body></html>")
			return
		}

		fmt.Fprintf(w, "<html><body><h1>Login Successful!</h1><p>You can now close this tab.</p></body></html>")
		tokenChan <- token
	})

	// Start the local web server in a new goroutine.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
			errChan <- err
		}
	}()
	fmt.Printf("Akash loginURL: %s\n", loginURL)
	fmt.Println("Launching browser for user login (ready for OIDC)...")
	fmt.Printf("Login URL: %s?redirect_uri=%s\n", loginURL, redirectURI)

	// Open the default browser tab.
	// NOTE: We pass the redirect_uri to signal our intent, even if the current page ignores it.
	if err := openBrowser(loginURL + "?redirect_uri=" + redirectURI); err != nil {
		log.Fatalf("could not open browser: %v", err)
	}

	fmt.Println("Waiting for user to complete login in the browser (OIDC check)...")
	fmt.Println("Waiting for up to 1 minute for an automated redirect.")

	// Set up a context with a shorter timeout for the *OIDC redirect check*.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Wait for the token via the OIDC redirect, an error, or the timeout.
	select {
	case token := <-tokenChan:
		// Success! This will happen when you switch to an OIDC provider.
		fmt.Printf("\nLogin successful via automated redirect. Received token:\n%s\n", token)
		// Go to shutdown section below.

	case err := <-errChan:
		fmt.Printf("\nAutomated login failed: %v\n", err)
		goto manualLogin // Jump to manual input/API fallback

	case <-ctx.Done():
		fmt.Println("\nAutomated redirect timed out after 1 minute.")
		// The current ClearBlade page doesn't redirect, so we proceed to manual input.
		goto manualLogin
	}

	// If we successfully received the token via the channel, we skip the manual step.
	// We also don't need to shutdown the server here because we hit the end of the function.
	// The server shutdown is handled below.

	// If we successfully received the token, we exit the main logic block and skip the manual section.
	goto cleanup

	// 2. Manual Token Input / API Fallback
manualLogin:
	fmt.Println("\n--- ClearBlade does not support automatic token passing. ---")
	fmt.Println("--- You must now manually enter your token or credentials. ---")

	// Option A: Ask user to manually paste the 'dev_token' from the network traffic.
	// This is difficult for non-technical users, so Option B is better.

	// **Option B: Fallback to the secure CLI/API login.**
	// NOTE: You must insert the `authenticateUser` function here or import it.

	// For now, let's use a simplified manual paste, as it's the closest to the browser experience:
	fmt.Println("\nPLEASE NOTE: Since the browser tab didn't redirect, the only way to get the token is to:")
	fmt.Println("1. Open your browser's Developer Tools (F12) on the login page.")
	fmt.Println("2. Login with email/password.")
	fmt.Println("3. Find the response from the '/admin/auth' request.")
	fmt.Println("4. Copy the 'dev_token' value.")
	fmt.Println("\nAlternatively, you can provide your credentials via the CLI:")
	fmt.Print("Enter 'token' to paste the token, or 'api' to use the API login: ")

	fmt.Fscanln(os.Stdin, &choice)

	if choice == "token" {
		fmt.Print("Paste 'dev_token' here and press Enter: ")
		fmt.Fscanln(os.Stdin, &finalToken)
	} else if choice == "api" {
		// You'll need the authenticateUser function imported/defined here.
		// Assuming you use the CLI authentication function from the previous answer.
		// token, err := authenticateUser(URL)
		// if err != nil { /* handle error */ }
		// finalToken = token

		// For simplicity in this example, we'll keep the manual paste.
		fmt.Println("API login not implemented in this snippet. Please choose 'token'.")
	}

	if finalToken != "" {
		tokenChan <- finalToken // Push the manually acquired token into the channel.
	} else {
		errChan <- fmt.Errorf("manual token input failed or was skipped")
	}

	// Wait for the manual token to be processed by the select statement.
	select {
	case token := <-tokenChan:
		fmt.Printf("\nLogin successful via manual input. Received token:\n%s\n", token)
	case err := <-errChan:
		fmt.Printf("\nLogin failed: %v\n", err)
	}

	// 3. Cleanup (Server Shutdown)
cleanup:
	// Shut down the server gracefully.
	fmt.Println("\nShutting down local server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}
	// NOTE: The browser tab still remains open and must be closed manually.
}

func retrieveTokenFromLocalStorage(url string) (string, error) {
	// 1. Create a browser context with HEADLESS set to FALSE
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),   // <-- THIS IS THE KEY CHANGE
		chromedp.Flag("disable-gpu", true), // Good practice, especially on Windows/Linux
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a generous timeout for the entire manual login process
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var token string
	var loginURL = url + "/login"
	const tokenKey = "ngStorage-cb_platform_dev_token"

	// The JavaScript to execute to read the token
	jsGetToken := fmt.Sprintf(`localStorage.getItem("%s");`, tokenKey)

	// Launch the browser and navigate
	err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
	)
	if err != nil {
		return "", fmt.Errorf("failed to launch browser or navigate: %w", err)
	}

	log.Println("A browser window has opened. Please complete the login manually.")
	log.Printf("Waiting for token '%s' to be set in local storage...", tokenKey)

	// 2. Poll the local storage until the token is found
	// We'll poll every 1 second for the token.
	for {
		select {
		case <-ctx.Done():
			// The overall 5-minute timeout was hit
			return "", fmt.Errorf("login timeout reached before token was set")
		default:
			// Execute JavaScript to read the local storage item
			err := chromedp.Run(ctx,
				chromedp.Evaluate(jsGetToken, &token),
			)
			if err != nil {
				// Log the error but continue polling, as it might be a temporary state
				log.Printf("Error during token check: %v. Retrying...", err)
			}

			if token != "" {
				log.Printf("Token successfully retrieved: %s\n", token)
				return token, nil
			}

			// Wait 1 second before checking again
			time.Sleep(1 * time.Second)
		}
	}
}

func promptAndFillMissingAuth(defaults *DefaultInfo, promptSet PromptSet) {
	// var defaultURL, defaultMsgURL, defaultEmail, defaultSystemKey string
	var defaultURL, defaultEmail, defaultSystemKey, defaultBrowserLogin string
	if defaults != nil {
		defaultURL = defaults.url
		// defaultMsgURL = defaults.msgUrl
		defaultEmail = defaults.email
		defaultSystemKey = defaults.systemKey
		defaultBrowserLogin = defaults.browserLogin
	}

	if !promptSet.Has(PromptSkipURL) {
		promptAndFillMissingURL(defaultURL)
	}

	// // TODO: messaging URL is optional since it can be derived from platform URL
	// // when not present
	// if !promptSet.Has(PromptSkipMsgURL) {
	// 	promptAndFillMissingMsgURL(defaultMsgURL)
	// }

	BrowserLogin := promptAndFillMissingBrowserLogin(defaultBrowserLogin)

	if BrowserLogin == "n" {
		if !promptSet.Has(PromptSkipEmail) {
			promptAndFillMissingEmail(defaultEmail)
		}

		if !promptSet.Has(PromptSkipPassword) {
			promptAndFillMissingPassword()
		}
	} else {
		// retrieveTokenViaRedirect(URL)
		token, err := retrieveTokenFromLocalStorage(URL)
		if err == nil {
			DevToken = strings.Trim(token, "\"") // remove double-quotes from returned token
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
	fmt.Printf("Akash URL: %s\nDevToken: %s\n", URL, DevToken)
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
