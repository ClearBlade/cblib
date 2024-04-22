package roles

import (
	"fmt"
	"os"

	"github.com/clearblade/cblib/diff"
	"github.com/clearblade/cblib/listutil"
	"github.com/clearblade/cblib/maputil"
)

type CollectionIdFetcher interface {
	GetCollectionIdByName(theNameWeWant string) (string, error)
}

func PackageRoleForUpdate(roleID string, role map[string]interface{}, fetcher CollectionIdFetcher) (map[string]interface{}, error) {
	permissions, ok := role["Permissions"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, fmt.Errorf("permissions for role do not exist or is not a map")
	}
	convertedPermissions := ConvertPermissionsStructure(permissions, fetcher)
	return map[string]interface{}{"ID": roleID, "Permissions": convertedPermissions}, nil
}

// The roles structure we get back when we retrieve roles is different from
// the format accepted for updating a role. Thus, we have this beauty of a
// conversion function. -swm
//
// THis is a gigantic cluster. We need to fix and learn from this. -swm
func ConvertPermissionsStructure(in map[string]interface{}, fetcher CollectionIdFetcher) map[string]interface{} {
	out := map[string]interface{}{}
	for key, valIF := range in {
		switch key {
		case "CodeServices":
			if valIF != nil {
				services, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for services permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				svcs := make([]map[string]interface{}, len(services))
				for idx, mapVal := range services {
					svcs[idx] = map[string]interface{}{
						"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
						"permissions": mapVal["Level"],
					}
				}
				out["services"] = removeDuplicatePermissions(svcs, "name")
			}
		case "Collections":
			if valIF != nil {
				collections, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for collections permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				cols := make([]map[string]interface{}, 0)
				for _, mapVal := range collections {
					collName := mapVal["Name"].(string)
					id, err := fetcher.GetCollectionIdByName(collName)
					if err != nil {
						fmt.Printf("Skipping permissions for collection '%s'; Error is - %s", collName, err.Error())
						continue
					}
					cols = append(cols, map[string]interface{}{
						"itemInfo":    map[string]interface{}{"id": id, "name": collName},
						"permissions": mapVal["Level"],
					})
				}
				out["collections"] = removeDuplicatePermissions(cols, "id")
			}
		case "DevicesList":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["devices"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "MsgHistory":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["msgHistory"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "SystemServices":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["system_services"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Portals":
			if valIF != nil {
				portals, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for portals permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				ptls := make([]map[string]interface{}, len(portals))
				for idx, mapVal := range portals {
					ptls[idx] = map[string]interface{}{
						"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
						"permissions": mapVal["Level"],
					}
				}
				out["portals"] = removeDuplicatePermissions(ptls, "name")
			}
		case "ExternalDatabases":
			if valIF != nil {
				externalDatabases, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for externalDatabases permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				extDbs := make([]map[string]interface{}, len(externalDatabases))
				for idx, mapVal := range externalDatabases {
					extDbs[idx] = map[string]interface{}{
						"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
						"permissions": mapVal["Level"],
					}
				}
				out["externaldatabases"] = removeDuplicatePermissions(extDbs, "name")
			}
		case "ServiceCaches":
			if valIF != nil {
				serviceCaches, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for serviceCaches permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				svcCaches := make([]map[string]interface{}, len(serviceCaches))
				for idx, mapVal := range serviceCaches {
					svcCaches[idx] = map[string]interface{}{
						"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
						"permissions": mapVal["Level"],
					}
				}
				out["servicecaches"] = removeDuplicatePermissions(svcCaches, "name")
			}
		case "Push":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["push"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Topics":
			if valIF != nil {
				topics, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for topics permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				tpcs := make([]map[string]interface{}, len(topics))
				for idx, mapVal := range topics {
					tpcs[idx] = map[string]interface{}{
						"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
						"permissions": mapVal["Level"],
					}
				}
				out["topics"] = tpcs
			}
		case "UsersList":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["users"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "EdgesList":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["edges"] = map[string]interface{}{"permissions": val["Level"]}
			}

		case "Triggers":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["triggers"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Timers":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["timers"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Deployments":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["deployments"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Roles":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["roles"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "AllCollections":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["allcollections"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "AllServices":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["allservices"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "ManageUsers":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["manageusers"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "AllExternalDatabases":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["allexternaldatabases"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Files":
			if valIF != nil {
				files, err := maputil.GetASliceOfMaps(valIF)
				if err != nil {
					fmt.Printf("Bad format for files permissions, not a slice of maps: %T\n", valIF)
					os.Exit(1)
				}
				theFiles := make([]map[string]interface{}, len(files))
				for idx, mapVal := range files {
					theFiles[idx] = map[string]interface{}{
						"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
						"permissions": mapVal["Level"],
					}
				}
				out["files"] = removeDuplicatePermissions(theFiles, "name")
			}
		case "usersecrets":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["usersecrets"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "adapters":
			if valIF != nil {
				val := maputil.GetMap(valIF)
				out["adapters"] = map[string]interface{}{"permissions": val["Level"]}
			}
		default:

		}
	}
	return out
}

// it's possible that there are duplicate permissions
// we need to remove any duplicates so that a role create/update succeeds
func removeDuplicatePermissions(perms []map[string]interface{}, idKey string) []map[string]interface{} {
	rtn := make([]map[string]interface{}, 0)
	foundIds := make(map[string]bool)

	for i := 0; i < len(perms); i++ {
		id := perms[i]["itemInfo"].(map[string]interface{})[idKey].(string)
		if _, found := foundIds[id]; !found {
			foundIds[id] = true
			rtn = append(rtn, perms[i])
		}
	}

	return rtn
}

func DiffRoles(local, backend []interface{}) *diff.UnsafeDiff {
	return listutil.CompareLists(local, backend, roleExists)
}

func roleExists(a interface{}, b interface{}) bool {
	return a == b
}
