// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/hashicorp/go-envparse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/go-vela/secret-vault/vault"
	"github.com/go-vela/server/compiler/types/raw"
)

var (
	// ErrNoKeysProvided defines the error type when a
	// no items were provided for a Vault read.
	ErrNoItemsProvided = errors.New("no items provided")

	// ErrNoPathProvided defines the error type when a
	// no path was provided for a Vault read.
	ErrNoPathProvided = errors.New("no `path` or `keys` provided")

	// ErrNoPathOrTargetProvided defines the error type when a
	// no path or target was provided for a Vault read key item.
	ErrNoPathOrTargetProvided = errors.New("no `path` or `target` provided for key item")

	// ErrNoKeysProvided defines the error type when a
	// no keys was provided for a Vault read.
	ErrNoSourceProvided = errors.New("no source provided")

	// appFS is a new os filesystem implementation for
	// interacting with modifications to the filesystem.
	appFS = afero.NewOsFs()

	// regexp for environment variable key.
	envVarNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

	// SecretVolumeLegacy defines volume that stores secrets during a build execution
	// in the legacy pattern where the user defines a directory for all keys
	//
	//nolint: gosec // false pos
	SecretVolumeLegacy = "/vela/secrets/%s/"

	// SecretVolume defines volume that stores secrets during a build execution
	//nolint: gosec // false pos
	SecretVolume = "/vela/secrets/%s"
)

type (
	// Read represents the plugin configuration reading secrets to the environment.
	Read struct {
		// is a list of items that are in a Vault instance
		Items []*Item
		// raw input of items provided for plugin
		RawItems string
		// vela masked outputs file location
		OutputsPath string
		// outputs map
		Outputs map[string]string
	}

	// Item represents how to read an item from a location and where to write it to.
	Item struct {
		// is the path to where the secret is stored in Vault
		Source string
		// are the paths to store the key in Vela
		Path raw.StringSlice
		// key overwrite option
		Keys map[string]KeyItem
	}

	KeyItem struct {
		Name   string
		Path   raw.StringSlice
		Target raw.StringSlice
	}
)

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (r *Read) Unmarshal() error {
	logrus.Trace("unmarshaling raw items")

	// cast raw items into bytes
	bytes := []byte(r.RawItems)

	// serialize raw items into expected Items type
	err := json.Unmarshal(bytes, &r.Items)
	if err != nil {
		return err
	}

	return nil
}

// Custom unmarshal for KeyItem to translate slice to map.
func (i *Item) UnmarshalJSON(data []byte) error {
	expectedInput := new(struct {
		Source string          `json:"source"`
		Path   raw.StringSlice `json:"path"`
		Keys   []KeyItem       `json:"keys"`
	})

	err := json.Unmarshal(data, expectedInput)
	if err != nil {
		return err
	}

	i.Source = expectedInput.Source
	i.Path = expectedInput.Path

	if len(expectedInput.Keys) > 0 {
		i.Keys = make(map[string]KeyItem)

		for _, keyItem := range expectedInput.Keys {
			if len(keyItem.Name) == 0 {
				return fmt.Errorf("key item missing name")
			}

			i.Keys[keyItem.Name] = keyItem
		}
	}

	return nil
}

// Exec runs the read for collecting secrets.
func (r *Read) Exec(v *vault.Client) error {
	logrus.Debug("running plugin with provided configuration")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// gather existing encoded outputs
	rawOutputs, err := a.ReadFile(r.OutputsPath)
	if err != nil {
		logrus.Debug("empty masked outputs file. creating one...")
	}

	r.Outputs, err = envparse.Parse(bytes.NewReader(rawOutputs))
	if err != nil {
		return fmt.Errorf("error parsing masked outputs file")
	}

	// decode existing secrets which will be base64 encoded
	for k, v := range r.Outputs {
		// decode the base64 value
		decodedValue, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return fmt.Errorf("unable to decode base64 value for key %s", k)
		}

		r.Outputs[k] = string(decodedValue)
	}

	for _, item := range r.Items {
		if len(item.Keys) > 0 {
			// if keys are defined, use new key based handling
			logrus.Debug("iterating through configured key items")

			err := r.execKeyItem(v, a, item)
			if err != nil {
				return err
			}
		} else {
			// if no keys defined, use legacy path based handling
			logrus.Debug("key items not configured. using legacy path handling")

			err := r.execLegacyPath(v, a, item)
			if err != nil {
				return err
			}
		}
	}

	if len(r.Outputs) > 0 {
		err := writeEnvFiles(a, r.Outputs, r.OutputsPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// execLegacyPath handles the old format of reading `source` and iterating over
// data keys and writing to the /vela/secrets/<path>/<key> file.
//
// it also populates the outputs map with the default key of VELA_SECRETS_<PATH>_<KEY>.
func (r *Read) execLegacyPath(v *vault.Client, a *afero.Afero, item *Item) error {
	secret, err := v.Read(item.Source)
	if err != nil {
		return err
	}

	for _, pth := range item.Path {
		// read data from the vault provider
		logrus.Tracef("reading data from path %s", item.Source)

		// remove any leading slashes from path
		p := strings.TrimPrefix(pth, "/")

		// remove any trailing slashes from path
		p = strings.TrimSuffix(p, "/")

		err = r.writeLegacySecretFiles(a, p, secret.Data)
		if err != nil {
			return err
		}

		if r.OutputsPath != "" {
			for _, v := range secret.Data {
				envKey := sanitizeEnvKey(strings.ToUpper(strings.TrimPrefix(p, "/")))

				r.Outputs[envKey] = v.(string)
			}
		}
	}

	return nil
}

// execKeyItem iterates over the defined keys from the `source` and writes them to their defined paths
// or environment variables.
func (r *Read) execKeyItem(v *vault.Client, a *afero.Afero, item *Item) error {
	secret, err := v.Read(item.Source)
	if err != nil {
		return err
	}

	for _, keyItem := range item.Keys {
		value, ok := secret.Data[keyItem.Name]
		if !ok {
			return fmt.Errorf("key %s not found in vault secret at %s", keyItem.Name, item.Source)
		}

		for _, pth := range keyItem.Path {
			// remove any leading slashes from path
			p := strings.TrimPrefix(pth, "/")

			// remove any trailing slashes from path
			p = strings.TrimSuffix(p, "/")

			// set the location of where to write the secret
			path := fmt.Sprintf(SecretVolume, p)

			// send Filesystem call to create directory path for .netrc file
			logrus.Tracef("creating directories in path %s", p)

			err := a.MkdirAll(filepath.Dir(path), 0777)
			if err != nil {
				return err
			}

			// set the secret in the Vela temp build volume
			logrus.Tracef("write data to file %s", p)

			err = a.WriteFile(path, []byte(value.(string)), 0600)
			if err != nil {
				return err
			}
		}

		for _, target := range keyItem.Target {
			r.Outputs[target] = value.(string)
		}
	}

	return nil
}

func (r *Read) writeLegacySecretFiles(a *afero.Afero, path string, data map[string]interface{}) error {
	// set the location of where to write the secret
	target := fmt.Sprintf(SecretVolumeLegacy, path)

	// send Filesystem call to create directory path for .netrc file
	logrus.Tracef("creating directories in path %s", path)

	err := a.MkdirAll(filepath.Dir(target), 0777)
	if err != nil {
		return err
	}

	// loop through keys in vault secret
	for k, v := range data {
		path = target + k

		// set the secret in the Vela temp build volume
		logrus.Tracef("write data to file %s", path)

		err = a.WriteFile(path, []byte(v.(string)), 0600)
		if err != nil {
			return err
		}

		// default env key
		if r.OutputsPath != "" {
			envKey := sanitizeEnvKey(strings.ToUpper(strings.TrimPrefix(path, "/")))

			r.Outputs[envKey] = v.(string)
		}
	}

	return nil
}

// Validate verifies the Copy is properly configured.
func (r *Read) Validate() error {
	logrus.Trace("validating read plugin configuration")

	if len(r.Items) == 0 {
		return ErrNoItemsProvided
	}

	for i, item := range r.Items {
		// verify that at least one path was provided
		if len(item.Path) == 0 && len(item.Keys) == 0 {
			return fmt.Errorf("%w for item %d %s", ErrNoPathProvided, i, r.RawItems)
		}

		if len(item.Keys) > 0 {
			for k, keyItem := range item.Keys {
				// verify that at least one path was provided for key item
				if len(keyItem.Path) == 0 && len(keyItem.Target) == 0 {
					return fmt.Errorf("%w for key item %s in item %d", ErrNoPathOrTargetProvided, k, i)
				}
			}
		} else {
			noPath := 0

			for _, path := range item.Path {
				// verify that at least one non-nil path was provided
				if len(path) != 0 {
					noPath = 1
					break
				}
			}

			if noPath == 0 {
				return fmt.Errorf("%w for item %d %s", ErrNoPathProvided, i, r.RawItems)
			}
		}

		// verify source is provided
		if len(item.Source) == 0 {
			return fmt.Errorf("%w for item %d", ErrNoSourceProvided, i)
		}
	}

	return nil
}

// writeEnvFiles creates the k=v pairs and writes them to the correct outputs path.
func writeEnvFiles(fs *afero.Afero, outputs map[string]string, path string) error {
	buffer := new(bytes.Buffer)

	keys := make([]string, 0, len(outputs))
	for k := range outputs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := outputs[k]

		if !envVarNamePattern.MatchString(k) {
			return fmt.Errorf("invalid environment variable name: %s", k)
		}

		encoded := base64.StdEncoding.EncodeToString([]byte(v))

		value := "'" + encoded + "'"

		fmt.Fprintf(buffer, "%s=%s\n", k, value)
	}

	if len(buffer.Bytes()) > 0 {
		err := fs.WriteFile(path, buffer.Bytes(), 0600)
		if err != nil {
			logrus.Warn("error writing secret values to outputs file. values will not be masked if accidentally logged, nor will they be available in the environment.")

			//nolint:nilerr // error string can contain sensitive information
			return nil
		}

		logrus.Info("successfully wrote secrets to outputs file")
	}

	return nil
}

// sanitizeEnvKey is a helper function that copies the key locator logic from the godotenv library
// and applies it to the vault keys for outputs.
func sanitizeEnvKey(s string) string {
	bytes := []byte(s)

	// use same logic as godotenv
	for i, c := range bytes {
		if unicode.IsLetter(rune(c)) || unicode.IsNumber(rune(c)) {
			continue
		}

		bytes[i] = '_'
	}

	return string(bytes)
}
