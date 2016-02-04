package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

func init() {
	pullCommand := &SubCommand{
		name:  "pull",
		usage: "pull a specified resource from a system",
		run:   doPull,
	}
	pullCommand.flags.BoolVar(&UserSchema, "userschema", false, "diff user table schema")
	pullCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to diff")
	pullCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to diff")
	pullCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to diff")
	pullCommand.flags.StringVar(&User, "user", "", "Name of user to diff")
	pullCommand.flags.StringVar(&RoleName, "role", "", "Name of role to diff")
	pullCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to diff")
	pullCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to diff")
	AddCommand("pull", pullCommand)
}

func doPull(cmd *SubCommand, cli *cb.DevClient, args ...string) error {
	setRootDir(".")
	systemInfo, err := getSysMeta()
	setupDirectoryStructure(systemInfo)
	if err != nil {
		return err
	}

	// ??? we already have them locally
	if r, err := pullRoles(systemInfo.Key, cli, false); err != nil {
		return err
	} else {
		rolesInfo = r
	}

	if ServiceName != "" {
		fmt.Printf("Pulling service %+s\n", ServiceName)
		if svc, err := pullService(systemInfo.Key, ServiceName, cli); err != nil {
			return err
		} else {
			writeService(ServiceName, svc)
		}
	}

	if LibraryName != "" {
		fmt.Printf("Pulling library %s\n", LibraryName)
		if lib, err := pullLibrary(systemInfo.Key, LibraryName, cli); err != nil {
			return err
		} else {
			writeLibrary(lib["name"].(string), lib)
		}
	}

	if CollectionName != "" {
		exportRows = true
		fmt.Printf("Pulling collection %+s\n", CollectionName)
		if allColls, err := cli.GetAllCollections(systemInfo.Key); err != nil {
			return err
		} else {
			var collID string
			// iterate over allColls and find one with matching name
			for _, c := range allColls {
				coll := c.(map[string]interface{})
				if CollectionName == coll["name"] {
					collID = coll["collectionID"].(string)
				}
			}
			if len(collID) < 1 {
				return fmt.Errorf("Collection %s not found.", CollectionName)
			}
			if coll, err := cli.GetCollectionInfo(collID); err != nil {
				return err
			} else {
				if data, err := pullCollection(systemInfo, coll, cli); err != nil {
					return err
				} else {
					writeCollection(data["name"].(string), data)
				}
			}
		}
	}

	if User != "" {
		fmt.Printf("Pulling user %+s\n", User)
		if users, err := cli.GetAllUsers(systemInfo.Key); err != nil {
			return err
		} else {
			ok := false
			for _, user := range users {
				if user["email"] == User {
					ok = true
					userId := user["user_id"].(string)
					if roles, err := cli.GetUserRoles(systemInfo.Key, userId); err != nil {
						return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
					} else {
						user["roles"] = roles
					}
					writeUser(User, user)
				}
			}
			if !ok {
				return fmt.Errorf("User %+s not found\n", User)
			}
		}
		if col, err := pullUserSchemaInfo(systemInfo.Key, cli, true); err != nil {
			return err
		} else {
			writeUserSchema(col)
		}
	}

	if RoleName != "" {
		roles := make([]map[string]interface{}, 0)
		splitRoles := strings.Split(RoleName, ",")
		for _, role := range splitRoles {
			fmt.Printf("Pulling role %+s\n", role)
			if r, err := pullRole(systemInfo.Key, role, cli); err != nil {
				return err
			} else {
				roles = append(roles, r)
				writeRole(role, r)
			}
		}
		storeRoles(roles)
	}

	if TriggerName != "" {
		fmt.Printf("Pulling trigger %+s\n", TriggerName)
		if trigg, err := pullTrigger(systemInfo.Key, TriggerName, cli); err != nil {
			return err
		} else {
			writeTrigger(TriggerName, trigg)
		}
	}

	if TimerName != "" {
		fmt.Printf("Pulling timer %+s\n", TimerName)
		if timer, err := pullTimer(systemInfo.Key, TimerName, cli); err != nil {
			return err
		} else {
			writeTimer(TimerName, timer)
		}
	}
	return nil
}

func pullRole(systemKey string, roleName string, cli *cb.DevClient) (map[string]interface{}, error) {
	r, err := cli.GetAllRoles(systemKey)
	if err != nil {
		return nil, err
	}
	ok := false
	var rval map[string]interface{}
	for _, rIF := range r {
		r := rIF.(map[string]interface{})
		if r["Name"].(string) == roleName {
			ok = true
			rval = r
		}
	}
	if !ok {
		return nil, fmt.Errorf("Role %s not found\n", roleName)
	}
	return rval, nil
}

func pullService(systemKey string, serviceName string, cli *cb.DevClient) (map[string]interface{}, error) {
	if service, err := cli.GetServiceRaw(systemKey, serviceName); err != nil {
		return nil, err
	} else {
		service["code"] = strings.Replace(service["code"].(string), "\\n", "\n", -1)
		return service, nil
	}
}

func pullLibrary(systemKey string, libraryName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetLibrary(systemKey, libraryName)
}

func pullTrigger(systemKey string, triggerName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetEventHandler(systemKey, triggerName)
}

func pullTimer(systemKey string, timerName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetTimer(systemKey, timerName)
}
