// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// config arguments loaded for the plugin
	Config *Config
	// read arguments loaded for the plugin
	Read *Read
}

// Exec runs the Vault plugin to read secrets into the Vela platform.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// setup connection with Vault
	vault, err := p.Config.New()
	if err != nil {
		return err
	}

	err = p.Read.Exec(vault)
	if err != nil {
		return err
	}

	logrus.Info("read secrets to volume")

	return nil
}

// Validate verifies the plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Config.Validate()
	if err != nil {
		return err
	}

	// validate read configuration
	err = p.Read.Validate()
	if err != nil {
		return err
	}

	return nil
}
