package auth

import cb "github.com/clearblade/Go-SDK"

func AuthorizeUsing(platformURL, messagingURL, devEmail, devPassword, devToken string) (*cb.DevClient, error) {
	var client *cb.DevClient

	if devPassword != "" {
		client = cb.NewDevClientWithAddrs(platformURL, messagingURL, devEmail, devPassword)
	} else if devToken != "" {
		client = cb.NewDevClientWithTokenAndAddrs(platformURL, messagingURL, devToken, devEmail)
	}

	_, err := client.Authenticate()
	if err != nil {
		return nil, err
	}

	return client, nil
}
