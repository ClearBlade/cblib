package roles

import (
	"fmt"

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

	if convertedPermissions, err := ConvertPermissionsStructure(permissions, fetcher); err != nil {
		return nil, err
	} else {
		return map[string]interface{}{"ID": roleID, "Permissions": convertedPermissions}, nil
	}
}

// The roles structure we get back when we retrieve roles is different from
// the format accepted for updating a role. Thus, we have this beauty of a
// conversion function. -swm
//
// THis is a gigantic cluster. We need to fix and learn from this. -swm
func ConvertPermissionsStructure(in map[string]interface{}, fetcher CollectionIdFetcher) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for key, valIF := range in {
		if valIF == nil {
			continue
		}

		switch key {
		case "CodeServices":
			if err := writeAssetLevelPermission("services", out, valIF); err != nil {
				return nil, err
			}
		case "Collections":
			collections, err := maputil.GetASliceOfMaps(valIF)
			if err != nil {
				return nil, fmt.Errorf("bad format for collections permissions, not a slice of maps: %T\n", valIF)
			}
			cols := make([]map[string]interface{}, 0)
			for _, mapVal := range collections {
				collName := mapVal["Name"].(string)
				id, err := fetcher.GetCollectionIdByName(collName)
				if err != nil {
					return nil, fmt.Errorf("could not get id for collection named %q: %s", collName, err)
				}
				cols = append(cols, map[string]interface{}{
					"itemInfo":    map[string]interface{}{"id": id, "name": collName},
					"permissions": mapVal["Level"],
				})
			}
			out["collections"] = removeDuplicatePermissions(cols, "id")
		case "DevicesList":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["devices"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "MsgHistory":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["msgHistory"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "SystemServices":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["system_services"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Portals":
			if err := writeAssetLevelPermission("portals", out, valIF); err != nil {
				return nil, err
			}
		case "EdgeRemoteAdmin":
			if err := writeAssetLevelPermission("edgeremoteadmin", out, valIF); err != nil {
				return nil, err
			}
		case "ExternalDatabases":
			if err := writeAssetLevelPermission("externaldatabases", out, valIF); err != nil {
				return nil, err
			}
		case "ServiceCaches":
			if err := writeAssetLevelPermission("servicecaches", out, valIF); err != nil {
				return nil, err
			}
		case "Push":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["push"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Topics":
			if err := writeAssetLevelPermission("topics", out, valIF); err != nil {
				return nil, err
			}
		case "UsersList":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["users"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "EdgesList":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["edges"] = map[string]interface{}{"permissions": val["Level"]}
			}

		case "Triggers":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["triggers"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Timers":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["timers"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Deployments":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["deployments"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Roles":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["roles"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "AllCollections":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["allcollections"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "AllServices":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["allservices"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "ManageUsers":
			if valIF != nil {
				if val, err := maputil.GetMap(valIF); err != nil {
					return nil, err
				} else {
					out["manageusers"] = map[string]interface{}{"permissions": val["Level"]}
				}
			}
		case "AllExternalDatabases":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["allexternaldatabases"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "Files":
			if err := writeAssetLevelPermission("files", out, valIF); err != nil {
				return nil, err
			}
		case "Filestores":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["filestores"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "FileStoreFiles":
			if err := writeAssetLevelPermission("filestorefiles", out, valIF); err != nil {
				return nil, err
			}
		case "usersecrets":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["usersecrets"] = map[string]interface{}{"permissions": val["Level"]}
			}
		case "adapters":
			if val, err := maputil.GetMap(valIF); err != nil {
				return nil, err
			} else {
				out["adapters"] = map[string]interface{}{"permissions": val["Level"]}
			}
		default:

		}
	}
	return out, nil
}

func writeAssetLevelPermission(permName string, out map[string]any, valIF any) error {
	assetPerms, err := maputil.GetASliceOfMaps(valIF)
	if err != nil {
		return fmt.Errorf("bad format for %s permissions, not a slice of maps: %T", permName, valIF)
	}
	formattedPerms := make([]map[string]interface{}, len(assetPerms))
	for idx, mapVal := range assetPerms {
		formattedPerms[idx] = map[string]any{
			"itemInfo":    map[string]any{"name": mapVal["Name"]},
			"permissions": mapVal["Level"],
		}
	}
	out[permName] = removeDuplicatePermissions(formattedPerms, "name")
	return nil
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

func DiffRoles(local, backend []string) *diff.UnsafeDiff[string] {
	return listutil.CompareLists(local, backend, roleExists)
}

func roleExists(a, b string) bool {
	return a == b
}
