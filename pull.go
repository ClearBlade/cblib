package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

func init() {
	pullCommand := &SubCommand{
		name:         "pull",
		usage:        "pull a specified resource from a system",
		needsAuth:    true,
		mustBeInRepo: true,
		run:          doPull,
	}

	pullCommand.flags.BoolVar(&AllServices, "all-services", false, "pull all services from system")
	pullCommand.flags.BoolVar(&AllLibraries, "all-libraries", false, "pull all libraries from system")
	pullCommand.flags.BoolVar(&AllEdges, "all-edges", false, "pull all edges from system")
	pullCommand.flags.BoolVar(&AllDevices, "all-devices", false, "pull all devices from system")
	pullCommand.flags.BoolVar(&AllDashboards, "all-dashboards", false, "pull all dashboards from system")
	pullCommand.flags.BoolVar(&AllPlugins, "all-plugins", false, "pull all plugins from system")
	pullCommand.flags.BoolVar(&UserSchema, "userschema", false, "pull user table schema")
	
	pullCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to pull")
	pullCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to pull")
	pullCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to pull")
	pullCommand.flags.StringVar(&User, "user", "", "Name of user to pull")
	pullCommand.flags.StringVar(&RoleName, "role", "", "Name of role to pull")
	pullCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to pull")
	pullCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to pull")
	pullCommand.flags.StringVar(&EdgeName, "edge", "", "Name of edge to pull")
	pullCommand.flags.StringVar(&DeviceName, "device", "", "Name of device to pull")
	pullCommand.flags.StringVar(&DashboardName, "dashboard", "", "Name of dashboard to pull")
	pullCommand.flags.StringVar(&PluginName, "plugin", "", "Name of plugin to pull")

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

	didSomething := false

	if AllServices {
		didSomething = true
		fmt.Printf("Pulling all services:")
		if _, err := pullServices(systemInfo.Key, cli); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if ServiceName != "" {
		didSomething = true
		fmt.Printf("Pulling service %+s\n", ServiceName)
		if svc, err := pullService(systemInfo.Key, ServiceName, cli); err != nil {
			return err
		} else {
			writeService(ServiceName, svc)
		}
	}

	if AllLibraries {
		didSomething = true
		fmt.Printf("Pulling all libraries:")
		if _, err := pullLibraries(systemInfo, cli); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if LibraryName != "" {
		didSomething = true
		fmt.Printf("Pulling library %s\n", LibraryName)
		if lib, err := pullLibrary(systemInfo.Key, LibraryName, cli); err != nil {
			return err
		} else {
			writeLibrary(lib["name"].(string), lib)
		}
	}

	if CollectionName != "" {
		didSomething = true
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
		didSomething = true
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
		didSomething = true
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
		didSomething = true
		fmt.Printf("Pulling trigger %+s\n", TriggerName)
		if trigg, err := pullTrigger(systemInfo.Key, TriggerName, cli); err != nil {
			return err
		} else {
			writeTrigger(TriggerName, trigg)
		}
	}

	if TimerName != "" {
		didSomething = true
		fmt.Printf("Pulling timer %+s\n", TimerName)
		if timer, err := pullTimer(systemInfo.Key, TimerName, cli); err != nil {
			return err
		} else {
			writeTimer(TimerName, timer)
		}
	}

	if AllDevices {
		didSomething = true
		fmt.Printf("Pulling all devices:")
		if _, err := pullDevices(systemInfo, cli); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if DeviceName != ""{
		didSomething = true
		fmt.Printf("Pulling device %+s\n", DeviceName)
		if device, err := pullDevice(systemInfo.Key, DeviceName, cli); err != nil {
			return err
		} else {
			writeDevice(DeviceName, device)
		}
	}

	if AllEdges {
		didSomething = true
		fmt.Printf("Pulling all edges:")
		if _, err := pullEdges(systemInfo, cli); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if EdgeName != ""{
		didSomething = true
		fmt.Printf("Pulling edge %+s\n", EdgeName)
		if edge, err := pullEdge(systemInfo.Key, EdgeName, cli); err != nil {
			return err
		} else {
			writeEdge(EdgeName, edge)
		}
	}

	if AllDashboards {
		didSomething = true
		fmt.Printf("Pulling all dashboards:")
		if _, err := pullDashboards(systemInfo, cli); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if DashboardName != ""{
		didSomething = true
		fmt.Printf("Pulling dashboard %+s\n", DashboardName)
		if dashboard, err := pullDashboard(systemInfo.Key, DashboardName, cli); err != nil {
			return err
		} else {
			writeDashboard(DashboardName, dashboard)
		}
	}

	if AllPlugins {
		didSomething = true
		fmt.Printf("Pulling all plugins:")
		if _, err := pullPlugins(systemInfo, cli); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if PluginName != ""{
		didSomething = true
		fmt.Printf("Pulling plugin %+s\n", PluginName)
		if plugin, err := pullPlugin(systemInfo.Key, PluginName, cli); err != nil {
			return err
		} else {
			writePlugin(PluginName, plugin)
		}
	}

	if !didSomething {
		fmt.Printf("Nothing to pull -- you must specify something to pull (ie, -service=<svc_name>)\n")
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

func pullDevice(systemKey string, deviceName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetDevice(systemKey, deviceName)
}

func pullEdge(systemKey string, edgeName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetEdge(systemKey, edgeName)
}

func pullDashboard(systemKey string, dashboardName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetDashboard(systemKey, dashboardName)
}

func pullPlugin(systemKey string, pluginName string, cli *cb.DevClient) (map[string]interface{}, error) {
	return cli.GetPlugin(systemKey, pluginName)
}


