// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/urfave/cli/v2"

	"github.com/go-vela/secret-vault/vault"
)

// helper function to load the flags into the plugin configuration.
func flags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "items",
			Usage:   "list of items to extract from a Vault",
			EnvVars: []string{"PARAMETER_ITEMS", "ITEMS"},
		},
	}

	// Add the Vault specific flags
	f = append(f, vault.Flags...)

	return f
}
