// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the runtime.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{

	// Logging Flags

	&cli.StringFlag{
		EnvVars: []string{"VAULT_LOG_LEVEL", "VELA_LOG_LEVEL", "LOG_LEVEL"},
		Name:    "log.level",
		Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
		Value:   "info",
	},

	// Config Flags

	&cli.StringFlag{
		EnvVars: []string{"PARAMETERS_VAULT_ADDR", "SECRET_VAULT_ADDR", "VELA_VAULT_ADDR", "VAULT_ADDR"},
		Name:    "config.addr",
		Usage:   "name of runtime driver to use",
	},
	&cli.StringFlag{
		EnvVars: []string{"PARAMETERS_VAULT_TOKEN", "SECRET_VAULT_TOKEN", "VELA_VAULT_TOKEN", "VAULT_TOKEN"},
		Name:    "config.token",
		Usage:   "vault token for storing secrets",
	},
	&cli.StringFlag{
		EnvVars: []string{"PARAMETERS_AUTH_METHOD", "SECRET_AUTH_METHOD", "VAULT_AUTH_METHOD"},
		Name:    "config.auth-method",
		Usage:   "vault token for storing secrets",
	},

	&cli.StringFlag{
		EnvVars: []string{"PARAMETERS_PATH", "PATH"},
		Name:    "path",
		Usage:   "path to a secret stored in vault",
	},
}
