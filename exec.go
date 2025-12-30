package cblib

import (
	"fmt"

	cb "github.com/clearblade/Go-SDK"
)

var (
	loggingEnabled bool
)

func init() {

	usage :=
		`
	Execute a ClearBlade code service on the ClearBlade Platform
	`

	example :=
		`
	cb-cli exec -service=Service1				# Execute the code service Service1 on the Platform
	`

	execCommand := &SubCommand{
		name:      "exec",
		usage:     usage,
		needsAuth: true,
		run:       doExec,
		example:   example,
	}

	execCommand.flags.BoolVar(&loggingEnabled, "enable-logs", false, "Enable logs for the service")
	execCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to execute")

	AddCommand("exec", execCommand)
}

func doExec(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}

	resp, err := client.CallService(systemInfo.Key, ServiceName, nil, loggingEnabled)
	if err != nil {
		return err
	}

	if success, ok := resp["success"].(bool); ok && success == false {
		return fmt.Errorf("%s", resp["results"].(string))
	}

	fmt.Println(resp)

	return nil
}
