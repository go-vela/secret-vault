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
			EnvVars: []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "LOG_LEVEL"},
			Name:    "log.level",
			Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:   "info",
		},
	}

	// Add the Vault specific flags
	f = append(f, vault.Flags...)

	return f
}
