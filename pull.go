package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

type Pull struct {
	SysKey        string
	DevToken      string
	Service       string
	Collection    string
	User          string
	Roles         []string
	Trigger       string
	Timer         string
	URL           string
	SysMeta       *System_meta
	SystemDotJSON map[string]interface{}
}

func Pull_cmd(sysKey, devToken string, args []string) error {
	fmt.Printf("Initializing...")
	p := &Pull{
		URL:           URL,
		SysKey:        sysKey,
		DevToken:      devToken,
		SystemDotJSON: map[string]interface{}{},
	}
	fmt.Printf("Done\n")
	return p.Cmd(args[3:])
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

	cli, err := auth(p.DevToken)
	if err != nil {
		return err
	}

	if sysMeta, err := pullSystemMeta(p.SysKey, cli); err != nil {
		return err
	} else {
		p.SysMeta = sysMeta
	}

	if err := setupDirectoryStructure(p.SysMeta); err != nil {
		return err
	}

	if val := p.Service; len(val) > 0 {
		fmt.Printf("Pulling service %+v\n", val)
		if svc, err := pullService(p.SysKey, val, cli); err != nil {
			return err
		} else {
			svcs := []map[string]interface{}{svc}
			if err := storeServices("", svcs, p.SysMeta); err != nil {
				return err
			}
			p.SystemDotJSON["services"] = svcs
		}
	}
	if val := p.Collection; len(val) > 0 {
		fmt.Printf("Pulling collection %+v\n", val)
		if co, err := cli.GetCollectionInfo(val); err != nil {
			return err
		} else {
			if data, err := pullCollection(co, p.SysMeta, cli); err != nil {
				return err
			} else {
				p.SystemDotJSON["data"] = data
			}
		}
	}
	if val := p.User; len(val) > 0 {
		fmt.Printf("Pulling user %+v\n", val)
		if users, err := cli.GetAllUsers(p.SysKey); err != nil {
			return err
		} else {
			ok := false
			for _, user := range users {
				if user["email"] == val {
					ok = true
					userId := user["user_id"].(string)
					if roles, err := cli.GetUserRoles(p.SysKey, userId); err != nil {
						return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
					} else {
						user["roles"] = roles
					}
					writeUsersFile([]map[string]interface{}{user})
				}
			}
			if !ok {
				return fmt.Errorf("User %+v not found\n", val)
			}
		}
		if col, err := pullUserColumns(p.SysKey, cli); err != nil {
			return err
		} else {
			p.SystemDotJSON["users"] = col
		}
	}
	if val := p.Roles; len(val) > 0 {
		roles := make([]map[string]interface{}, 0)
		for _, role := range val {
			fmt.Printf("Pulling role %+v\n", role)
			if r, err := pullRole(p.SysKey, role, cli); err != nil {
				return err
			} else {
				roles = append(roles, r)
			}
		}
		p.SystemDotJSON["roles"] = roles
	}
	if val := p.Trigger; len(val) > 0 {
		fmt.Printf("Pulling trigger %+v\n", val)
		if trigg, err := pullTrigger(p.SysKey, val, cli); err != nil {
			return err
		} else {
			p.SystemDotJSON["triggers"] = []interface{}{trigg}
		}
	}
	if val := p.Timer; len(val) > 0 {
		fmt.Printf("Pulling timer %+v\n", val)
		if timer, err := pullTimer(p.SysKey, val, cli); err != nil {
			return err
		} else {
			p.SystemDotJSON["timer"] = []interface{}{timer}
		}
	}

	if err := storeSystemDotJSON(p.SystemDotJSON); err != nil {
		return err
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
