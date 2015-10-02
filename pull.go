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
	AddCommand("pull", pullCommand)
}

func doPull(cmd *SubCommand, cli *cb.DevClient, args ...string) error {
	p := &Pull{
		CLI:    cli,
		SysKey: args[0],
	}
	return p.Cmd(args[1:])
}

type Pull struct {
	SysKey     string
	Service    string
	Library    string
	Collection string
	User       string
	Roles      []string
	Trigger    string
	Timer      string
	CLI        *cb.DevClient
	SysMeta    *System_meta
}

func (p Pull) Cmd(args []string) error {
	for _, arg := range args {
		s := strings.Split(arg, "=")
		if len(s) != 2 {
			return fmt.Errorf("invalid argument for %+v\n", s[0])
		}
		switch s[0] {
		case "service":
			p.Service = s[1]
		case "library":
			p.Library = s[1]
		case "collection":
			p.Collection = s[1]
		case "user":
			p.User = s[1]
		case "roles":
			p.Roles = strings.Split(s[1], ",")
		case "trigger":
			p.Trigger = s[1]
		case "timer":
			p.Timer = s[1]
		default:
			return fmt.Errorf("option \"%+v\" not supported", s[0])
		}
	}
	if sysMeta, err := pullSystemMeta(p.SysKey, p.CLI); err != nil {
		return err
	} else {
		p.SysMeta = sysMeta
	}
	setRootDir(strings.Replace(p.SysMeta.Name, " ", "_", -1))
	if err := setupDirectoryStructure(p.SysMeta); err != nil {
		return err
	}
	storeMeta(p.SysMeta)
	storeSystemDotJSON(systemDotJSON)

	if val := p.Service; len(val) > 0 {
		fmt.Printf("Pulling service %+s\n", val)
		if svc, err := pullService(p.SysKey, val, p.CLI); err != nil {
			return err
		} else {
			writeService(val, svc)
		}
	}
	if val := p.Library; len(val) > 0 {
		fmt.Printf("Pulling library %s\n", val)
		if lib, err := pullLibrary(p.SysKey, val, p.CLI); err != nil {
			return err
		} else {
			writeLibrary(lib["name"].(string), lib)
		}
	}
	if val := p.Collection; len(val) > 0 {
		fmt.Printf("Pulling collection %+s\n", val)
		if co, err := p.CLI.GetCollectionInfo(val); err != nil {
			return err
		} else {
			if data, err := pullCollection(p.SysMeta, co, p.CLI); err != nil {
				return err
			} else {
				writeCollection(data["name"].(string), data)
			}
		}
	}
	if val := p.User; len(val) > 0 {
		fmt.Printf("Pulling user %+s\n", val)
		if users, err := p.CLI.GetAllUsers(p.SysKey); err != nil {
			return err
		} else {
			ok := false
			for _, user := range users {
				if user["email"] == val {
					ok = true
					userId := user["user_id"].(string)
					if roles, err := p.CLI.GetUserRoles(p.SysKey, userId); err != nil {
						return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
					} else {
						user["roles"] = roles
					}
					writeUser(val, user)
				}
			}
			if !ok {
				return fmt.Errorf("User %+s not found\n", val)
			}
		}
		if col, err := pullUserSchemaInfo(p.SysKey, p.CLI, true); err != nil {
			return err
		} else {
			writeUserSchema(col)
		}
	}
	if val := p.Roles; len(val) > 0 {
		roles := make([]map[string]interface{}, 0)
		for _, role := range val {
			fmt.Printf("Pulling role %+s\n", role)
			if r, err := pullRole(p.SysKey, role, p.CLI); err != nil {
				return err
			} else {
				roles = append(roles, r)
				writeRole(role, r)
			}
		}
		storeRoles(roles)
	}
	if val := p.Trigger; len(val) > 0 {
		fmt.Printf("Pulling trigger %+s\n", val)
		if trigg, err := pullTrigger(p.SysKey, val, p.CLI); err != nil {
			return err
		} else {
			writeTrigger(val, trigg)
		}
	}
	if val := p.Timer; len(val) > 0 {
		fmt.Printf("Pulling timer %+s\n", val)
		if timer, err := pullTimer(p.SysKey, val, p.CLI); err != nil {
			return err
		} else {
			writeTimer(val, timer)
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
