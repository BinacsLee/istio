// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"

	"istio.io/istio/pkg/test"
)

// AsBytes is a simple wrapper around ioutil.ReadFile provided for completeness.
func AsBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// AsBytesOrFail calls AsBytes and fails the test if any errors occurred.
func AsBytesOrFail(t test.Failer, filename string) []byte {
	t.Helper()
	content, err := AsBytes(filename)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

// AsString is a convenience wrapper around ioutil.ReadFile that converts the content to a string.
func AsString(filename string) (string, error) {
	bytes, err := AsBytes(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// AsStringOrFail calls AsBytesOrFail and then converts to string.
func AsStringOrFail(t test.Failer, filename string) string {
	t.Helper()
	return string(AsBytesOrFail(t, filename))
}

// NormalizePath expands the homedir (~) and returns an error if the file doesn't exist.
func NormalizePath(originalPath string) (string, error) {
	if originalPath == "" {
		return "", nil
	}
	// trim leading/trailing spaces from the path and if it uses the homedir ~, expand it.
	var err error
	out := strings.TrimSpace(originalPath)
	out, err = homedir.Expand(out)
	if err != nil {
		return "", err
	}

	// Verify that the file exists.
	if _, err := os.Stat(out); os.IsNotExist(err) {
		return "", fmt.Errorf("failed normalizing file %s: %v", originalPath, err)
	}

	return out, nil
}
