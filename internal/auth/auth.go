package auth

import cb "github.com/clearblade/Go-SDK"

// AuthorizeUsing returns a *cb.DevClient authenticated and ready to use. The
// password and the token are mutually exclusive, if the token is given, it is
// gonna be used instead of the password for creating the client.
func AuthorizeUsing(platformURL, messagingURL, devEmail, devPassword, devToken string) (*cb.DevClient, error) {
	var client *cb.DevClient

	if devToken != "" {
		client = cb.NewDevClientWithTokenAndAddrs(platformURL, messagingURL, devToken, devEmail)
	} else if devPassword != "" {
		client = cb.NewDevClientWithAddrs(platformURL, messagingURL, devEmail, devPassword)
	}

	_, err := client.Authenticate()
	if err != nil {
		return nil, err
	}

	return client, nil
}
