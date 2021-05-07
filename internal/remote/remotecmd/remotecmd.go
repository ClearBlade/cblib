// package remotecmd contains the CLI for the remote subcommand.
package remotecmd

import (
	"fmt"

	"github.com/clearblade/cblib/internal/auth"
	"github.com/clearblade/cblib/internal/remote"
	"github.com/urfave/cli/v2"
)

const (
	Name  = "remote"
	Usage = "manage remotes"
)

var (
	flagName         string
	flagPlatformURL  string
	flagMessagingURL string
	flagDevEmail     string
	flagDevPassword  string
	flagDevToken     string
	flagSystemKey    string
)

const (
	errRemoteNotFound = "not a remote"
)

// remoteCommand implements all the actions that are trigged by the CLI interactions.
// It uses dependency inversion to receive the remotes that are gonna be managed.
type remoteCommand struct {
	remotes *remote.Remotes
}

// New returns a new *cli.Command instance that implements the CLI interactions
// against the given *remote.Remotes instance.
func New(remotes *remote.Remotes) *cli.Command {
	cmd := &remoteCommand{remotes}
	return &cli.Command{
		Name:  Name,
		Usage: Usage,
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list remotes",
				Action:  cmd.list,
			},
			{
				Name:  "put",
				Usage: "create or update remotes",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "name of the remote",
						Required:    true,
						Destination: &flagName,
					},
					&cli.StringFlag{
						Name:        "platform-url",
						Usage:       "Platform URL to use",
						Required:    true,
						Destination: &flagPlatformURL,
					},
					&cli.StringFlag{
						Name:        "messaging-url",
						Usage:       "Messaging URL to use",
						Required:    true,
						Destination: &flagMessagingURL,
					},
					&cli.StringFlag{
						Name:        "dev-email",
						Usage:       "Developer email with access to the platform",
						Required:    true,
						Destination: &flagDevEmail,
					},
					&cli.StringFlag{
						Name:        "dev-password",
						Usage:       "Developer password",
						Required:    false,
						Destination: &flagDevPassword,
					},
					&cli.StringFlag{
						Name:        "dev-token",
						Usage:       "Developer token",
						Required:    false,
						Destination: &flagDevToken,
					},
					&cli.StringFlag{
						Name:        "system-key",
						Usage:       "System to use for the remote",
						Required:    true,
						Destination: &flagSystemKey,
					},
				},
				Action: cmd.put,
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "remove remotes",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "name of the remote",
						Required:    true,
						Destination: &flagName,
					},
				},
				Action: cmd.remove,
			},
			{
				Name:  "set-current",
				Usage: "set current remote",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "name of the remote",
						Required:    true,
						Destination: &flagName,
					},
				},

				Action: cmd.setCurrent,
			},
		},
	}
}

func (rc *remoteCommand) list(c *cli.Context) error {
	curr, _ := rc.remotes.Current()

	prefixIsCurrent := "* "
	prefixNoCurrent := "  "

	fmt.Printf("%s%s %s %s\n", prefixNoCurrent, "Name", "URL", "SystemKey")

	for _, r := range rc.remotes.List() {
		prefix := prefixNoCurrent
		if r == curr {
			prefix = prefixIsCurrent
		}
		fmt.Printf("%s%s %s %s\n", prefix, r.Name, r.PlatformURL, r.SystemKey)
	}

	return nil
}

func (rc *remoteCommand) put(c *cli.Context) error {
	if flagDevPassword == "" && flagDevToken == "" {
		return fmt.Errorf("Must provide dev-password or dev-token")
	}

	client, err := auth.AuthorizeUsing(flagPlatformURL, flagMessagingURL, flagDevEmail, flagDevPassword, flagDevToken)
	if err != nil {
		return err
	}

	sys, err := client.GetSystem(flagSystemKey)
	if err != nil {
		return err
	}

	remote := &remote.Remote{
		Name:         flagName,
		PlatformURL:  flagPlatformURL,
		MessagingURL: flagMessagingURL,
		SystemKey:    sys.Key,
		SystemSecret: sys.Secret,
		Token:        client.DevToken,
	}

	err = rc.remotes.Put(remote)
	if err != nil {
		return err
	}

	fmt.Println("Remote created")
	return nil
}

func (rc *remoteCommand) remove(c *cli.Context) error {
	r, ok := rc.remotes.FindByName(flagName)
	if !ok {
		return fmt.Errorf("%s: %s", errRemoteNotFound, flagName)
	}

	err := rc.remotes.Remove(r)
	if err != nil {
		return err
	}

	fmt.Println("Remote removed")
	return nil
}

func (rc *remoteCommand) setCurrent(c *cli.Context) error {
	r, ok := rc.remotes.FindByName(flagName)
	if !ok {
		return fmt.Errorf("%s: %s", errRemoteNotFound, flagName)
	}

	err := rc.remotes.SetCurrent(r)
	if err != nil {
		return err
	}

	fmt.Println("Remote changed")
	return nil
}
