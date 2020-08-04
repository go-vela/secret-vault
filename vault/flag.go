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
		Name:    "log.level",
		Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
		Value:   "info",
		EnvVars: []string{"PARAMETER_LOG_LEVEL", "VAULT_LOG_LEVEL", "VELA_LOG_LEVEL", "LOG_LEVEL"},
	},

	// Config Flags
	&cli.StringFlag{
		Name:  "config.addr",
		Usage: "address to the instance",
	},
	&cli.StringFlag{
		Name:  "config.auth-method",
		Usage: "authentication method for interfacing instance - options: (token|ldap)",
	},
	&cli.StringFlag{
		Name:  "config.password",
		Usage: "password for server authentication with LDAP",
	},
	&cli.StringFlag{
		Name:  "config.token",
		Usage: "token for server authentication",
	},
	&cli.StringFlag{
		Name:  "config.username",
		Usage: "username for server authentication with LDAP",
	},
}
