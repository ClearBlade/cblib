package syspath

import (
	"reflect"
	"testing"
)

func TestServiceNameParsing(t *testing.T) {
	tests := []struct {
		path         string
		expectedName string
		shouldError  bool
	}{
		{"code/services/foo/foo.js", "foo", false},
		{"code/services/foo/bar.js", "", true},
		{"code/services/bar/bar.json", "bar", false},
		{"code/libraries/bar/bar.json", "", true},
		{".cb-cli/map-name-to-id/roles.json", "", true},
	}

	for _, test := range tests {
		name, err := GetServiceNameFromPath(test.path)
		didErr := err != nil
		if didErr != test.shouldError {
			t.Errorf("Expected error: %v, got error: %v", test.shouldError, didErr)
		}

		if name != test.expectedName {
			t.Errorf("Expected name %q for path %q, got %q", test.expectedName, test.path, name)
		}
	}
}

func TestLibraryNameParsing(t *testing.T) {
	tests := []struct {
		path         string
		expectedName string
		shouldError  bool
	}{
		{"code/libraries/foo/foo.js", "foo", false},
		{"code/libraries/foo/bar.js", "", true},
		{"code/libraries/bar/bar.json", "bar", false},
		{"code/services/bar/bar.json", "", true},
	}

	for _, test := range tests {
		name, err := GetLibraryNameFromPath(test.path)
		didErr := err != nil
		if didErr != test.shouldError {
			t.Errorf("Expected error: %v, got error: %v", test.shouldError, didErr)
		}

		if name != test.expectedName {
			t.Errorf("Expected name %q for path %q, got %q", test.expectedName, test.path, name)
		}
	}
}

func TestRolePathParsing(t *testing.T) {
	tests := []struct {
		path   string
		isRole bool
	}{
		{"roles/Authenticated.json", true},
		{"code/MyRole.json", false},
	}

	for _, test := range tests {
		isRole := IsRolePath(test.path)
		if isRole != test.isRole {
			t.Errorf("Expected isRolePath to return %v for %q, got: %v", test.isRole, test.path, isRole)
		}
	}
}

func TestUserDataPathParsing(t *testing.T) {
	tests := []struct {
		path        string
		email       string
		shouldError bool
	}{
		{"users/mwalowski@clearblade.com", "", true},
		{"users/mwalowski@clearblade.com.json", "mwalowski@clearblade.com", false},
		{"users/invalidemail.json", "invalidemail", false},
		{"code/email@gmail.com", "", true},
		{"users/roles/email@gmail.com.json", "", true},
	}

	for _, test := range tests {
		email, err := GetUserEmailFromDataPath(test.path)
		didErr := err != nil
		if didErr != test.shouldError {
			t.Errorf("Expected error: %v, got error: %v", test.shouldError, didErr)
		}

		if email != test.email {
			t.Errorf("Expected email %q for path %q, got %q", test.email, test.path, email)
		}
	}
}

func TestUserRolePathParsing(t *testing.T) {
	tests := []struct {
		path        string
		email       string
		shouldError bool
	}{
		{"users/roles/mwalowski@clearblade.com", "", true},
		{"users/roles/mwalowski@clearblade.com.json", "mwalowski@clearblade.com", false},
		{"users/roles/invalidemail.json", "invalidemail", false},
		{"code/roles/email@gmail.com", "", true},
		{"users/roles/email@gmail.com.json", "email@gmail.com", false},
	}

	for _, test := range tests {
		email, err := GetUserEmailFromRolePath(test.path)
		didErr := err != nil
		if didErr != test.shouldError {
			t.Errorf("Expected error: %v, got error: %v", test.shouldError, didErr)
		}

		if email != test.email {
			t.Errorf("Expected email %q for path %q, got %q", test.email, test.path, email)
		}
	}
}

func TestBucketSetPathParsing(t *testing.T) {
	tests := []struct {
		path        string
		name        string
		shouldError bool
	}{
		{"bucket-sets/Files.json", "Files", false},
		{"bucket-sets/myfiles.readme", "", true},
		{"bucket-sets/My Bucket Set.json", "My Bucket Set", false},
		{"randompath", "", true},
	}

	for _, test := range tests {
		name, err := GetBucketSetNameFromPath(test.path)
		didErr := err != nil
		if didErr != test.shouldError {
			t.Errorf("Expected error: %v, got error: %v", test.shouldError, didErr)
		}

		if name != test.name {
			t.Errorf("Expected name %q for path %q, got %q", test.name, test.path, name)
		}
	}
}

func TestBucketSetFilePathParsing(t *testing.T) {
	tests := []struct {
		path        string
		shouldError bool
		expected    FullBucketPath
	}{
		{"bucket-set-files/mybucket/mybox/File.zip", false, FullBucketPath{"mybucket", "mybox", "File.zip"}},
		{"bucket-set-files/mybucket/inbox/img.jpg", false, FullBucketPath{"mybucket", "inbox", "img.jpg"}},
		{"bucket-sets/another bucket/another box/anotherfile.bin", true, FullBucketPath{}},
		{"bucket-set-files/bucket/", true, FullBucketPath{}},
		{"bucket-set-files/my bucket/box/myfile/with/a/very/long/path.bin", false, FullBucketPath{"my bucket", "box", "myfile/with/a/very/long/path.bin"}},
	}

	for _, test := range tests {
		actual, err := ParseBucketPath(test.path)
		didErr := err != nil
		if didErr != test.shouldError {
			t.Fatalf("Expected error: %v, got error: %v", test.shouldError, didErr)
		}

		if !test.shouldError && !reflect.DeepEqual(&test.expected, actual) {
			t.Errorf("Expected %v for path %q, got %v", test.expected, test.path, actual)
		}
	}
}
