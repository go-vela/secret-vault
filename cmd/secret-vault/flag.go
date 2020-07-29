// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/urfave/cli/v2"

	"github.com/go-vela/secret-vault/vault"
)

// helper function to load the flags into
// the plugin configuration
func flags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "path",
			Usage:   "path to a secret stored in vault",
			EnvVars: []string{"PARAMETER_PATH", "PATH"},
		},
		&cli.StringSliceFlag{
			Name:    "keys",
			Usage:   "the keys to extract out of the item stored in Vault",
			EnvVars: []string{"PARAMETER_KEYS", "KEYS"},
		},
	}

	// Add the Vault specific flags
	f = append(f, vault.Flags...)

	return f
}
