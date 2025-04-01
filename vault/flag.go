// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"github.com/urfave/cli/v3"
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
		Name:    "config.addr",
		Usage:   "address to the instance",
		EnvVars: []string{"PARAMETER_ADDR", "SECRET_VAULT_ADDR", "VELA_VAULT_ADDR", "VAULT_ADDR"},
	},
	&cli.StringFlag{
		Name:    "config.auth-method",
		Usage:   "authentication method for interfacing instance - options: (token|ldap)",
		EnvVars: []string{"PARAMETER_AUTH_METHOD", "SECRET_AUTH_METHOD", "VAULT_AUTH_METHOD"},
	},
	&cli.StringFlag{
		Name:    "config.password",
		Usage:   "password for server authentication with LDAP",
		EnvVars: []string{"PARAMETER_PASSWORD", "SECRET_VAULT_PASSWORD", "VELA_VAULT_PASSWORD", "VAULT_PASSWORD"},
	},
	&cli.StringFlag{
		Name:    "config.token",
		Usage:   "token for server authentication",
		EnvVars: []string{"PARAMETER_TOKEN", "SECRET_VAULT_TOKEN", "VELA_VAULT_TOKEN", "VAULT_TOKEN"},
	},
	&cli.StringFlag{
		Name:    "config.username",
		Usage:   "username for server authentication with LDAP",
		EnvVars: []string{"PARAMETER_USERNAME", "SECRET_VAULT_USERNAME", "VELA_VAULT_USERNAME", "VAULT_USERNAME"},
	},
}
