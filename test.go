package cblib

import (
	"clearblade/cbjson"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
)

func init() {
	systemDotJSON = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	rolesInfo = []map[string]interface{}{}
	myTestCommand := &SubCommand{
		name:  "test",
		usage: "<fill this in>",
		run:   doTest,
		//  TODO -- add help, usage, etc.
	}
	myTestCommand.flags.StringVar(&ServiceName, "service", "", "name of service to test")
	myTestCommand.flags.StringVar(&Params, "params", "", "params to service in json stringified format")
	myTestCommand.flags.StringVar(&Topic, "topic", "", "message topic for publishing")
	myTestCommand.flags.StringVar(&Payload, "payload", "", "The message payload")
	myTestCommand.flags.BoolVar(&Push, "push", true, "Push the service prior to running")
	AddCommand("test", myTestCommand)
}

func doTest(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("Extra arguments passed to test command:%s\n", args)
	}

	setRootDir(".")
	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}
	if ServiceName != "" {
		if Push {
			if err = doPushService(systemInfo.Key, client); err != nil {
				return err
			}
			fmt.Printf("Sucessfully pushed service '%s'\n", ServiceName)
		}
		fmt.Printf("Running service '%s'... ", ServiceName)
		resMap, err := doCallService(systemInfo.Key, client)
		if err != nil {
			fmt.Printf("Failed\n")
			return err
		}
		fmt.Printf("Succeeded\nResults: %+v\n", resMap)
	} else if Topic != "" {
		if err = doPublishMessage(systemInfo.Key, client); err != nil {
			return err
		}
	} else if Payload != "" {
		return fmt.Errorf("-payload provided but -topic is missing")
	}

	return nil
}

func doPushService(systemKey string, client *cb.DevClient) error {
	svcMap, err := findService(systemKey, ServiceName)
	if err != nil {
		if err.Error() == NotExistErrorString {
			fmt.Printf("Service '%s' does not exist locally. Not pushing\n")
			return nil
		}
		return err
	}
	return updateService(systemKey, svcMap, client)
}

func doCallService(systemKey string, client *cb.DevClient) (map[string]interface{}, error) {
	parsedParams := map[string]interface{}{}
	var err error
	if Params != "" {
		parsedParams, _, err = cbjson.GetJSONFromString(Params)
		if err != nil {
			return nil, fmt.Errorf("Could not parse parameters string: %s", err.Error())
		}
	}
	resMap, err := client.CallService(systemKey, ServiceName, parsedParams, true /* turn on logging */)
	if err != nil {
		return nil, fmt.Errorf("Call service failed: %s", err.Error())
	}
	return resMap, nil
}

func doPublishMessage(systemKey string, client *cb.DevClient) error {
	return fmt.Errorf("test command: publishing a message not yet implemented")
}
