package cblib

import (
	"bufio"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
	"time"
)

func init() {
	pushCommand := &SubCommand{
		name:  "push",
		usage: "push a specified resource to a system",
		run:   doPush,
	}
	AddCommand("push", pushCommand)
}

func doPush(cmd *SubCommand, cli *cb.DevClient, args ...string) error {
	p := &Push{
		SysKey:  args[0],
		CLI:     cli,
		SysInfo: map[string]interface{}{},
	}
	return p.Cmd(args[1:])
}

type Push struct {
	SysKey     string
	DevToken   string
	Service    string
	Collection string
	User       string
	Roles      []string
	Trigger    string
	Timer      string
	CLI        *cb.DevClient
	SysInfo    map[string]interface{}
}

func (p Push) Cmd(args []string) error {
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

	if systemInfo, err := getDict("system.json"); err != nil {
		return err
	} else {
		p.SysInfo = systemInfo
	}
	setRootDir(strings.Replace(p.SysInfo["name"].(string), " ", "_", -1))

	if val := p.Service; len(val) > 0 {
		fmt.Printf("Pushing service %+v\n", val)
		ok := false
		for _, svc := range p.SysInfo["services"].([]interface{}) {
			service := svc.(map[string]interface{})
			if service["name"] == val {
				ok = true
				if err := createService(p.SysKey, service, p.CLI); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Service %+v not found\n", val)
		}
	}
	if val := p.Collection; len(val) > 0 {
		fmt.Printf("Pushing collection %+v\n", val)
		ok := false
		coll := p.SysInfo["data"].(map[string]interface{})
		if coll["collectionID"] == val {
			ok = true
			if err := createCollection(p.SysKey, coll, p.CLI); err != nil {
				return err
			}
		}
		if !ok {
			return fmt.Errorf("Collection %+v not found\n", val)
		}
	}
	if val := p.User; len(val) > 0 {
		fmt.Printf("Pushing user %+v\n", val)
		meta, err := pullSystemMeta(p.SysKey, p.CLI)
		if err != nil {
			return err
		}
		sysSecret := meta.Secret
		if users, err := getArray("users.json"); err != nil {
			return err
		} else {
			ok := false
			for _, user := range users {
				userMap := user.(map[string]interface{})
				if userMap["email"] == val {
					ok = true
					for _, userCol := range p.SysInfo["users"].([]interface{}) {
						column := userCol.(map[string]interface{})
						columnName := column["ColumnName"].(string)
						columnType := column["ColumnType"].(string)
						if err := p.CLI.CreateUserColumn(p.SysKey, columnName, columnType); err != nil {
							return fmt.Errorf("Could not create user column %s: %s", columnName, err.Error())
						}
					}
					userId := userMap["user_id"].(string)
					if roles, err := p.CLI.GetUserRoles(p.SysKey, userId); err != nil {
						return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
					} else {
						userMap["roles"] = roles
					}
					if _, err := createUser(p.SysKey, sysSecret, userMap, p.CLI); err != nil {
						return fmt.Errorf("Could not create user %s: %s", val, err.Error())
					}
				}
			}
			if !ok {
				return fmt.Errorf("User %+v not found\n", val)
			}
		}
	}
	if val := p.Roles; len(val) > 0 {
		ok := false
		roles, ok := p.SysInfo["roles"]
		if !ok {
			return fmt.Errorf("No roles found locally.\n")
		}
		for _, roleIF := range roles.([]interface{}) {
			role := roleIF.(map[string]interface{})
			for _, roleName := range p.Roles {
				if role["Name"] == roleName {
					ok = true
					fmt.Printf("Pushing role %+v\n", roleName)
					if err := updateRole(p.SysKey, role, p.CLI); err != nil {
						return err
					}
				}
			}
		}
		if !ok {
			return fmt.Errorf("Role %+v not found\n", val)
		}
	}
	if val := p.Trigger; len(val) > 0 {
		fmt.Printf("Pushing trigger %+v\n", val)
		triggers, ok := p.SysInfo["triggers"]
		if !ok {
			return fmt.Errorf("No triggers found locally.\n")
		}
		ok = false
		for _, triggIF := range triggers.([]interface{}) {
			trigger := triggIF.(map[string]interface{})
			if trigger["name"] == val {
				ok = true
				if err := updateTrigger(p.SysKey, trigger, p.CLI); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Trigger %+v not found\n", val)
		}
	}
	if val := p.Timer; len(val) > 0 {
		fmt.Printf("Pushing timer %+v\n", val)
		timers, ok := p.SysInfo["timers"]
		if !ok {
			return fmt.Errorf("No timers found locally.\n")
		}
		ok = false
		for _, timerIF := range timers.([]interface{}) {
			timer := timerIF.(map[string]interface{})
			if timer["name"] == val {
				ok = true
				if err := updateTimer(p.SysKey, timer, p.CLI); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Timer %+v not found\n", val)
		}
	}

	return nil
}

func createUser(systemKey string, systemSecret string, user map[string]interface{}, client *cb.DevClient) (string, error) {
	email := user["email"].(string)
	password := "password"
	if pwd, ok := user["password"]; ok {
		password = pwd.(string)
	}
	newUser, err := client.RegisterNewUser(email, password, systemKey, systemSecret)
	if err != nil {
		return "", fmt.Errorf("Could not create user %s: %s", email, err.Error())
	}
	userId := newUser["user_id"].(string)
	niceRoles := mungeRoles(user["roles"].([]interface{}))
	if len(niceRoles) > 0 {
		if err := client.AddUserToRoles(systemKey, userId, niceRoles); err != nil {
			return "", err
		}
	}
	return userId, nil
}

func createTrigger(sysKey string, trigger map[string]interface{}, client *cb.DevClient) error {
	triggerName := trigger["name"].(string)
	triggerDef := trigger["event_definition"].(map[string]interface{})
	trigger["def_module"] = triggerDef["def_module"]
	trigger["def_name"] = triggerDef["def_name"]
	trigger["system_key"] = sysKey
	delete(trigger, "name")
	delete(trigger, "event_definition")
	if _, err := client.CreateEventHandler(sysKey, triggerName, trigger); err != nil {
		return fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
	}
	return nil
}

func updateTrigger(systemKey string, trigger map[string]interface{}, client *cb.DevClient) error {
	triggerName := trigger["name"].(string)
	triggerDef := trigger["event_definition"].(map[string]interface{})
	trigger["def_module"] = triggerDef["def_module"]
	trigger["def_name"] = triggerDef["def_name"]
	trigger["system_key"] = systemKey
	delete(trigger, "name")
	delete(trigger, "event_definition")
	if _, err := client.UpdateEventHandler(systemKey, triggerName, trigger); err != nil {
		fmt.Printf("Could not update trigger %s\n", triggerName)
		fmt.Printf("Would you like to create a new trigger named %s? (Y/n)", triggerName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreateEventHandler(systemKey, triggerName, trigger); err != nil {
					return fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
				} else {
					fmt.Printf("Successfully created new trigger %s\n", triggerName)
				}
			} else {
				fmt.Printf("Trigger will not be created.\n")
			}
		}
	}
	return nil
}

func createTimer(systemKey string, timer map[string]interface{}, client *cb.DevClient) error {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		timer["start_time"] = time.Now().Format(time.RFC3339)
	}
	if _, err := client.CreateTimer(systemKey, timerName, timer); err != nil {
		return fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
	}
	return nil
}

func updateTimer(systemKey string, timer map[string]interface{}, client *cb.DevClient) error {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		timer["start_time"] = time.Now().Format(time.RFC3339)
	}
	if _, err := client.UpdateTimer(systemKey, timerName, timer); err != nil {
		fmt.Printf("Could not update timer %s\n", timerName)
		fmt.Printf("Would you like to create a new timer named %s? (Y/n)", timerName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreateEventHandler(systemKey, timerName, timer); err != nil {
					return fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
				} else {
					fmt.Printf("Successfully created new timer %s\n", timerName)
				}
			} else {
				fmt.Printf("Timer will not be created.\n")
			}
		}
	}
	return nil
}

func createService(systemKey string, service map[string]interface{}, client *cb.DevClient) error {
	svcName := service["name"].(string)
	svcParams := mkSvcParams(service["params"].([]interface{}))
	svcDeps := service["dependencies"].(string)
	svcCode, err := getServiceCode(svcName)
	delete(service, "code")
	if err != nil {
		return err
	}
	if err := client.NewServiceWithLibraries(systemKey, svcName, svcCode, svcDeps, svcParams); err != nil {
		return err
	}
	if enableLogs(service) {
		if err := client.EnableLogsForService(systemKey, svcName); err != nil {
			return err
		}
	}
	permissions := service["permissions"].(map[string]interface{})
	for roleId, level := range permissions {
		if err := client.AddServiceToRole(systemKey, svcName, roleId, int(level.(float64))); err != nil {
			return err
		}
	}
	return nil
}

func createLibrary(systemKey string, library map[string]interface{}, client *cb.DevClient) error {
	libName := library["name"].(string)
	delete(library, "name")
	delete(library, "version")
	if _, err := client.CreateLibrary(systemKey, libName, library); err != nil {
		return fmt.Errorf("Could not create library %s: %s", libName, err.Error())
	}
	return nil
}

func createCollection(systemKey string, collection map[string]interface{}, client *cb.DevClient) error {
	collectionName := collection["name"].(string)
	colId, err := client.NewCollection(systemKey, collectionName)
	if err != nil {
		return err
	}

	permissions := collection["permissions"].(map[string]interface{})
	for roleId, level := range permissions {
		if err := client.AddCollectionToRole(systemKey, colId, roleId, int(level.(float64))); err != nil {
			return err
		}
	}

	columns := collection["schema"].([]interface{})
	for _, columnIF := range columns {
		column := columnIF.(map[string]interface{})
		colName := column["ColumnName"].(string)
		colType := column["ColumnType"].(string)
		if colName == "item_id" {
			continue
		}
		if err := client.AddColumn(colId, colName, colType); err != nil {
			return err
		}
	}

	//  Add the items
	itemsIF, err := getCollectionItems(collectionName)
	if err != nil {
		return err
	}
	items := make([]map[string]interface{}, len(itemsIF))
	for idx, itemIF := range itemsIF {
		items[idx] = itemIF.(map[string]interface{})
	}
	if _, err := client.CreateData(colId, items); err != nil {
		return err
	}
	return nil
}

func updateRole(systemKey string, role map[string]interface{}, cli *cb.DevClient) error {
	roleName := role["Name"].(string)
	if err := cli.UpdateRole(systemKey, roleName, role); err != nil {
		return fmt.Errorf("Role %s not updated\n", roleName)
	}
	return nil
}
