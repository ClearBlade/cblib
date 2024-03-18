package cblib

import (
	//"fmt"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/types"
)

const BACKUP_DIRECTORY_SUFFIX = "_cb_bak"

type compare func(sliceOfSystemResources *[]interface{}, i int, j int) bool

func setupAddrs(paddr string, maddr string) {
	cb.CB_ADDR = paddr

	preIdx := strings.Index(paddr, "://")
	if maddr == "" {
		if preIdx != -1 {
			maddr = paddr[preIdx+3:]
		} else {
			maddr = paddr
		}
	}
	postIdx := strings.Index(maddr, ":")
	if postIdx != -1 {
		cb.CB_MSG_ADDR = maddr[:postIdx] + ":1883"
	} else {
		cb.CB_MSG_ADDR = maddr + ":1883"
	}
}

// processURLs processes the given platform and messaging URL(s) for correctness.
// If the messaging URL is not provided it is derived from the platform URL.
func processURLs(platformURL, messagingURL string) (string, string, error) {

	platformURL = strings.TrimSpace(platformURL)
	messagingURL = strings.TrimSpace(messagingURL)

	purl, err := url.Parse(platformURL)
	if err != nil {
		return "", "", fmt.Errorf("error parsing plaform URL: %s", err)
	}

	if !purl.IsAbs() {
		return "", "", fmt.Errorf("platform URL must specify a scheme (http, https, etc)")
	}

	var mhost, mport string

	if len(messagingURL) <= 0 {
		mhost = purl.Hostname()
		mport = "1883"

	} else if strings.Contains(messagingURL, ":") {
		mhost, mport, err = net.SplitHostPort(messagingURL)
		if err != nil {
			return "", "", fmt.Errorf("error parsing host and port from messaging URL: %s", err)
		}

	} else {
		mhost = messagingURL
		mport = "1883"
	}

	finalPlatformURL := fmt.Sprintf("%s://%s", purl.Scheme, purl.Host)
	finalMesagingURL := fmt.Sprintf("%s:%s", mhost, mport)

	return finalPlatformURL, finalMesagingURL, nil
}

// Bubble sort, compare by map key
func sortByMapKey(arrayPointer *[]interface{}, sortKey string) {
	if arrayPointer == nil {
		return
	}
	array := *arrayPointer
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(array)-1; i++ {
			needToSwap := compareWithKey(sortKey, arrayPointer, i+1, i)
			if needToSwap {
				swap(arrayPointer, i, i+1)
				swapped = true
			}
		}
	}
}

// Bubble sort, compare by function
func sortByFunction(arrayPointer *[]interface{}, compareFn compare) {
	if arrayPointer == nil {
		return
	}
	array := *arrayPointer
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(array)-1; i++ {
			needToSwap := compareFn(arrayPointer, i+1, i)
			if needToSwap {
				swap(arrayPointer, i, i+1)
				swapped = true
			}
		}
	}
}

func swap(array *[]interface{}, i, j int) {
	tmp := (*array)[j]
	(*array)[j] = (*array)[i]
	(*array)[i] = tmp
}

func isString(input interface{}) bool {
	return input != nil && reflect.TypeOf(input).Name() == "string"
}

func compareWithKey(sortKey string, sliceOfCodeServices *[]interface{}, i, j int) bool {
	slice := *sliceOfCodeServices

	map1, castSuccess1 := slice[i].(map[string]interface{})
	map2, castSuccess2 := slice[j].(map[string]interface{})

	if !castSuccess1 || !castSuccess2 {
		return false
	}

	name1 := map1[sortKey]
	name2 := map2[sortKey]

	if !isString(name1) || !isString(name2) {
		return false
	}

	return name1.(string) < name2.(string)
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createFilePath(args ...string) string {
	return strings.Join(args, "/")
}

func copyMap(daMap map[string]interface{}) map[string]interface{} {
	rtn := make(map[string]interface{})
	for k, v := range daMap {
		rtn[k] = v
	}
	return rtn
}

func getBackupDirectoryName(directoryName string) string {
	return directoryName + BACKUP_DIRECTORY_SUFFIX
}

func removeBackupDirectory(directoryName string) error {
	return removeDirectory(getBackupDirectoryName(directoryName))
}

func backupAndCleanDirectory(directoryName string) error {
	if err := backupDirectory(directoryName); err != nil {
		return err
	}
	return removeDirectoryContents(directoryName)
}

func removeDirectoryContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func backupDirectory(directoryName string) error {
	return copyDir(directoryName, getBackupDirectoryName(directoryName))
}

func removeDirectory(directoryName string) error {
	if err := os.RemoveAll(directoryName); err != nil && err != os.ErrNotExist {
		// if we have an error that doesn't relate to the directory not existing, let's return the error
		return err
	}
	return nil
}

func restoreBackupDirectory(directoryName string) error {
	if err := removeDirectory(directoryName); err != nil && err != os.ErrNotExist {
		fmt.Printf("Error while restoring backup directory for '%s'; Unable to remove destination directory", directoryName)
		return err
	}
	if err := copyDir(getBackupDirectoryName(directoryName), directoryName); err != nil {
		fmt.Printf("Error while restoring backup directory for '%s'; Unable to copy backup directory", directoryName)
		return err
	}
	if err := removeDirectory(getBackupDirectoryName(directoryName)); err != nil {
		fmt.Printf("Error while restoring backup directory for '%s'; Unable to remove backup directory", directoryName)
		return err
	}
	return nil
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// These keys are generated upon GET, and not representative of the data model
// If we store to filesystem with these keys, the corresponding PUT/POST for portal fails
func removeBlacklistedPortalKeys(portal map[string]interface{}) map[string]interface{} {
	var blacklist = []string{"permissions", "plugins"}
	for _, key := range blacklist {
		delete(portal, key)
	}
	return portal
}

type ListDiff struct {
	add    []interface{}
	remove []interface{}
}

func isDefaultColumn(defaultColumns []string, colName string) bool {
	for i := 0; i < len(defaultColumns); i++ {
		if defaultColumns[i] == colName {
			return true
		}
	}
	return false
}

// FilterSlice returns the items of the slice `s` for which `predicate` returns true.
func FilterSlice(s []interface{}, predicate func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0, len(s))

	for _, item := range s {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func compareListsAndFilter(after []interface{}, before []interface{}, compare func(interface{}, interface{}) bool, filter func(interface{}) bool) *UnsafeDiff {
	diff := UnsafeDiff{after, before, nil, nil, compare}
	Diff(&diff)
	diff.Added = FilterSlice(diff.Added, filter)
	diff.Removed = FilterSlice(diff.Removed, filter)
	return &diff
}

func compareLists(after []interface{}, before []interface{}, compare func(interface{}, interface{}) bool) *UnsafeDiff {
	diff := UnsafeDiff{after, before, nil, nil, compare}
	Diff(&diff)
	return &diff
}

func convertStringSliceToInterfaceSlice(strs []string) []interface{} {
	rtn := make([]interface{}, len(strs))
	for i, s := range strs {
		rtn[i] = s
	}
	return rtn
}

func convertInterfaceSliceToStringSlice(ifaces []interface{}) []string {
	rtn := make([]string, len(ifaces))
	for i, s := range ifaces {
		rtn[i] = s.(string)
	}
	return rtn
}

func myLogger(str string) {
	fmt.Printf("\n\n%s\n\n", str)
}

func logError(err string) {
	myLogger(fmt.Sprintf("[ERROR] %s", err))

}

func logInfo(info string) {
	myLogger(fmt.Sprintf("[INFO] %s", info))
}

func logWarning(info string) {
	myLogger(fmt.Sprintf("[WARNING] %s", info))
}

func logErrorForUpdatingMapFile(fileName string, err error) {
	logError(fmt.Sprintf("Failed to update %s - subsequent operations may fail. Error is - %s", fileName, err.Error()))
}

func confirmPrompt(question string) (bool, error) {
	if AutoApprove {
		fmt.Println("-auto-approve is true. Creating entity...")
		return true, nil
	}
	fmt.Printf("\n%s (Y/n)", question)
	reader := bufio.NewReader(os.Stdin)
	if text, err := reader.ReadString('\n'); err != nil {
		return false, err
	} else {
		if strings.Contains(strings.ToUpper(text), "Y") {
			return true, nil
		} else {
			return false, nil
		}
	}
}

type countRequestFunc = func(systemKey string, query *cb.Query) (cb.CountResp, error)
type dataRequestFunc = func(systemKey string, query *cb.Query) ([]interface{}, error)

func paginateRequests(systemKey string, pageSize int, cf countRequestFunc, df dataRequestFunc) ([]interface{}, error) {
	u, err := cf(systemKey, nil)
	if err != nil {
		return nil, err
	}

	rtn := make([]interface{}, 0)
	for i := 0; i*pageSize < int(u.Count); i++ {
		pageQuery := cb.NewQuery()
		pageQuery.PageNumber = i + 1
		pageQuery.PageSize = pageSize
		data, err := df(systemKey, pageQuery)
		if err != nil {
			return nil, err
		}
		rtn = append(rtn, data...)
	}
	return rtn, nil
}

func getUserEmailByID(id string) (string, error) {
	u, err := getUserEmailToId()
	if err != nil {
		return id, err
	}
	for email, userID := range u {
		if userID == id {
			return email, nil
		}
	}
	// couldn't find a match, just return the id
	return id, nil
}

type requestFunc = func() (interface{}, error)

func retryRequest(funk requestFunc, maxRetries int) (interface{}, error) {
	numOfRetries := 0

	var recur func() (interface{}, error)
	recur = func() (interface{}, error) {
		data, err := funk()
		if err != nil {
			retryNumber := numOfRetries + 1
			logError(err.Error())
			if numOfRetries < maxRetries {
				logInfo(fmt.Sprintf("Retrying request number %d out of %d", retryNumber, maxRetries))
				numOfRetries++
				return recur()
			}
			return nil, err
		}
		return data, nil
	}
	return recur()
}

func replaceUserIdWithEmailInTriggerKeyValuePairs(trig map[string]interface{}, userEmailToId map[string]interface{}) {
	// check to see
	if kv, ok := trig["key_value_pairs"].(map[string]interface{}); ok {
		if thisUserID, ok := kv["userId"]; ok {
			for email, userID := range userEmailToId {
				if thisUserID == userID {
					delete(kv, "userId")
					kv["email"] = email
				}
			}

		}
	}
}

func isTriggerForSpecificUser(trig map[string]interface{}) (string, map[string]interface{}, bool) {
	if kv, ok := trig["key_value_pairs"]; ok {
		if userEmail, ok := kv.(map[string]interface{})["email"]; ok {
			return userEmail.(string), kv.(map[string]interface{}), ok
		}
	}
	return "", nil, false
}

func replaceEmailWithUserIdInTriggerKeyValuePairs(trig map[string]interface{}, usersInfo []UserInfo) {
	if userEmail, kv, ok := isTriggerForSpecificUser(trig); ok {
		// found an email that we stored on the FS. need to remove it and replace with the users new user_id
		delete(kv, "email")
		if usersInfo != nil {
			for i := 0; i < len(usersInfo); i++ {
				if usersInfo[i].Email == userEmail {
					kv["userId"] = usersInfo[i].UserID
				}
			}
		}
	}
}

func getCollectionName(collection map[string]interface{}) (string, error) {
	collection_name, ok := collection["name"].(string)
	if !ok {
		return "", fmt.Errorf("No name in collection json file: %+v\n", collection)
	}
	return collection_name, nil
}

type CreateCollectionIfNecessaryOptions struct {
	pushItems bool
	pullItems bool
}

type CreateCollectionIfNecessaryOutput struct {
	collectionExistsOrWasCreated bool
}

func createCollectionIfNecessary(meta *types.System_meta, collection map[string]interface{}, client *cb.DevClient, options CreateCollectionIfNecessaryOptions) (CreateCollectionIfNecessaryOutput, error) {
	collection_name, err := getCollectionName(collection)
	if err != nil {
		return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: false}, err
	}

	_, err = client.GetDataTotalByName(meta.Key, collection_name, cb.NewQuery())
	if err != nil {
		fmt.Printf("Could not find collection '%s'. Error is - %s\n", collection_name, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new collection named %s?", collection_name))
		if err != nil {
			return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: false}, err
		} else {
			if c {
				if _, err := CreateCollection(meta.Key, collection, options.pushItems, client); err != nil {
					return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: false}, fmt.Errorf("Could not create collection %s: %s", collection_name, err.Error())
				} else {
					fmt.Printf("Successfully created new collection %s\n", collection_name)
				}
				if options.pullItems {
					fmt.Printf("Updating local copy... %s\n", collection_name)
					return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: true}, PullAndWriteCollection(meta, collection_name, client, true, true)
				}
				return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: true}, nil
			} else {
				fmt.Printf("Collection will not be created.\n")
				return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: false}, nil
			}
		}
	}
	return CreateCollectionIfNecessaryOutput{collectionExistsOrWasCreated: true}, nil
}

func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}
