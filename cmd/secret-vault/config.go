// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"

	"github.com/go-vela/secret-vault/vault"
	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Vault config information.
type Config struct {
	// enables setting the addr for a Vault instance
	Addr string
	// enables setting the type of authentication method
	AuthMethod string
	// enables setting the password for authentication
	Password string
	// enables setting the token for for authentication
	Token string
	// enables setting the username for authentication
	Username string
}

// New creates an Vault client for reading secrets.
func (c *Config) New() (*vault.Client, error) {
	logrus.Trace("creating new Vault client from plugin configuration")

	// add the Vault specific config info to setup a client
	s := &vault.Setup{
		Addr:       c.Addr,
		AuthMethod: c.AuthMethod,
		Password:   c.Password,
		Token:      c.Token,
		Username:   c.Username,
	}

	// setup connection with Vault
	client, err := vault.New(s)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify Addr is provided
	if len(c.Addr) == 0 {
		return fmt.Errorf("no config address provided")
	}

	if !strings.Contains(c.Addr, "://") {
		return fmt.Errorf("config address must be <scheme>://<hostname> format")
	}

	// verify Addr is provided
	if len(c.AuthMethod) == 0 {
		return fmt.Errorf("no auth method provided")
	}

	// verify provided authentication is valid for authentication
	switch c.AuthMethod {
	case vault.TokenAuthMethod:
		if len(c.Token) == 0 {
			return fmt.Errorf("invalid authentication passed. Must set token for %s auth method", vault.TokenAuthMethod)
		}

	case vault.LDAPAuthMethod:
		if len(c.Password) == 0 && len(c.Username) == 0 {
			return fmt.Errorf("invalid authentication passed. Must set username and password for %s auth method", vault.LDAPAuthMethod)
		}
	}

	return nil
}
