// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/urfave/cli/v3"

	"github.com/go-vela/secret-vault/vault"
)

// helper function to load the flags into the plugin configuration.
func flags() []cli.Flag {
	f := []cli.Flag{
		&cli.StringFlag{
			Name:    "items",
			Usage:   "list of items to extract from a Vault",
			Sources: cli.EnvVars("PARAMETER_ITEMS", "ITEMS"),
		},
		&cli.StringFlag{
			Sources: cli.EnvVars("VELA_MASKED_BASE64_OUTPUTS"),
			Name:    "vela.masked-outputs",
			Usage:   "env file path to store secrets",
		},
	}

	// Add the Vault specific flags
	f = append(f, vault.Flags...)

	return f
}
