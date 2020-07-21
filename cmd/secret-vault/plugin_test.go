// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/go-vela/secret-vault/vault"
)

func TestVault_Plugin_Exec(t *testing.T) {
	// TODO write this test
}

func TestVault_Plugin_Validate(t *testing.T) {
	// setup types
	// setup types
	tests := []struct {
		plugin *Plugin
		err    error
	}{
		{ // success with token config and read action
			plugin: &Plugin{
				Config: &Config{
					Addr:       "https://myvault.com/",
					AuthMethod: vault.TokenAuthMethod,
					Token:      "superSecretAPIKey",
				},
				Read: &Read{
					Path: "/path/to/secret",
					Keys: []string{"foobar"},
				},
			},
			err: nil,
		},
		{ // success with ldap config and read action
			plugin: &Plugin{
				Config: &Config{
					Addr:       "https://myvault.com/",
					AuthMethod: vault.LDAPAuthMethod,
					Password:   "superSecretPassword",
					Username:   "myusername",
				},
				Read: &Read{
					Path: "/path/to/secret",
					Keys: []string{"foobar"},
				},
			},
			err: nil,
		},
		Read: &Read{
			Path: "/path/to/secret",
			Keys: []string{"foobar"},
		},
	}

	// run test
	for _, test := range tests {
		err := test.plugin.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
