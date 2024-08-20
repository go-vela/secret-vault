// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-vela/secret-vault/vault"
	"github.com/go-vela/types/raw"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	// ErrNoKeysProvided defines the error type when a
	// no items were provided for a Vault read.
	ErrNoItemsProvided = errors.New("no items provided")

	// ErrNoPathProvided defines the error type when a
	// no path was provided for a Vault read.
	ErrNoPathProvided = errors.New("no path provided")

	// ErrNoKeysProvided defines the error type when a
	// no keys was provided for a Vault read.
	ErrNoSourceProvided = errors.New("no source provided")

	// appFS is a new os filesystem implementation for
	// interacting with modifications to the filesystem.
	appFS = afero.NewOsFs()

	// SecretVolume defines volume that stores secrets during a build execution
	//nolint: gosec // false pos
	SecretVolume = "/vela/secrets/%s/"
)

type (
	// Read represents the plugin configuration reading secrets to the environment.
	Read struct {
		// is a list of items that are in a Vault instance
		Items []*Item
		// raw input of items provided for plugin
		RawItems string
	}

	// Item represents how to read an item from a location and where to write it to.
	Item struct {
		// is the path to where the secret is stored in Vault
		Source string
		// are the paths to store the key in Vela
		Path raw.StringSlice
		// key overwrite option
		Keys map[string]string
	}
)

// Exec runs the read for collecting secrets.
func (r *Read) Exec(v *vault.Client) error {
	logrus.Debug("running plugin with provided configuration")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	for _, item := range r.Items {
		for _, pth := range item.Path {
			// remove any leading slashes from path
			path := strings.TrimPrefix(pth, "/")

			// remove any trailing slashes from path
			path = strings.TrimSuffix(path, "/")

			// read data from the vault provider
			logrus.Tracef("reading data from path %s", item.Source)

			secret, err := v.Read(item.Source)
			if err != nil {
				return err
			}

			// set the location of where to write the secret
			target := fmt.Sprintf(SecretVolume, path)

			// send Filesystem call to create directory path for .netrc file
			logrus.Tracef("creating directories in path %s", path)

			err = a.Fs.MkdirAll(filepath.Dir(target), 0777)
			if err != nil {
				return err
			}

			// loop through keys in vault secret
			for k, v := range secret.Data {
				path = target + k

				// if there is a key override, set the new path
				if overrideKey, ok := item.Keys[k]; ok {
					path = target + overrideKey
				}

				// set the secret in the Vela temp build volume
				logrus.Tracef("write data to file %s", path)

				err = a.WriteFile(path, []byte(v.(string)), 0600)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

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

// Validate verifies the Copy is properly configured.
func (r *Read) Validate() error {
	logrus.Trace("validating read plugin configuration")

	if len(r.Items) == 0 {
		return ErrNoItemsProvided
	}

	for i, item := range r.Items {
		// verify that at least one path was provided
		if len(item.Path) == 0 {
			return fmt.Errorf("%w for item %d %s", ErrNoPathProvided, i, r.RawItems)
		}

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

		// verify source is provided
		if len(item.Source) == 0 {
			return fmt.Errorf("%w for item %d", ErrNoSourceProvided, i)
		}
	}

	return nil
}
