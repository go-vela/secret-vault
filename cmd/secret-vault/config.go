// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/secret-vault/vault"
	"github.com/sirupsen/logrus"
)

// Config represents the plugin configuration for Vault config information.
type Config struct {
	// enables performing a request against a Vault instance
	Addr string
	// enables authenticating via a token against a Vault instance
	Token string
}

// New creates an Vault client for reading secrets.
func (c *Config) New() (*vault.Client, error) {
	logrus.Trace("creating new Vault client from plugin configuration")

	// setup connection with Vault
	client, err := vault.New(c.Addr, c.Token)
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

	// verify Token is provided
	if len(c.Token) == 0 {
		return fmt.Errorf("no config token provided")
	}

	return nil
}
