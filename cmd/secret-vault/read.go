// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-vela/secret-vault/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	// ErrNoPathProvided defines the error type when a
	// no path was provided for a Vault read
	ErrNoPathProvided = errors.New("no path provided")

	// ErrNoKeysProvided defines the error type when a
	// no keys was provided for a Vault read
	ErrNoKeysProvided = errors.New("no keys provided")

	// appFS is a new os filesystem implementation for
	// interacting with modifications to the filesystem
	appFS = afero.NewOsFs()
)

// Read represents the plugin configuration reading secrets to the environment.
type Read struct {
	// is the path to where the secret is stored
	Path string
	// is the keys to extract from the secret store
	Keys []string
}

// Exec runs the read for collecting secrets
func (r Read) Exec(v *vault.Client) error {
	logrus.Debug("running plugin with provided configuration")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	paths := strings.Split(r.Path, "/")
	name := paths[len(paths)-1]

	// read data from the vault provider
	secret, err := v.Read(r.Path)
	if err != nil {
		return err
	}

	// write data to environment
	for _, key := range r.Keys {
		// TODO support none key=value secrets in vault
		// m, ok := secret.Data["foo"].(map[string]interface{})
		// if !ok {
		// 	return fmt.Errorf("unable to extract secret data")
		// }

		// set the location of where to write the secret
		path := fmt.Sprintf("/vela/secrets/%s", strings.ToLower(name))

		// send Filesystem call to create directory path for .netrc file
		err = a.Fs.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return err
		}

		// set the secret in the Vela temp build volume
		//TODO consider making a const for the Vela secret path
		err = a.WriteFile(path, []byte(secret.Data[key].(string)), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate verifies the Copy is properly configured.
func (r Read) Validate() error {
	logrus.Trace("validating read plugin configuration")

	// verify path is provided
	if len(r.Path) == 0 {
		return ErrNoPathProvided
	}

	// verify keys is provided
	if len(r.Keys) == 0 {
		return ErrNoKeysProvided
	}

	return nil
}