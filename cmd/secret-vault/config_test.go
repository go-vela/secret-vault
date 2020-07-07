// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/go-vela/secret-vault/vault"
)

func TestVault_Config_New(t *testing.T) {
	// setup types
	tests := []struct {
		config *Config
		err    error
	}{
		{ // valid config with token auth method
			config: &Config{
				Addr:       "https://myvault.com/",
				AuthMethod: vault.TokenAuthMethod,
				Token:      "superSecretAPIKey",
			},
			err: nil,
		},
		// TODO: investigate how to put mock vault in a mode with a fake LDAP provider
		// { // valid config with ldap auth method
		// 	config: &Config{
		// 		Addr:       "https://myvault.com/",
		// 		AuthMethod: vault.LDAPAuthMethod,
		// 		Password:   "superSecretPassword",
		// 		Username:   "myusername",
		// 	},
		// 	err: nil,
		// },
	}

	// run test
	for _, test := range tests {
		got, err := test.config.New()
		if err != nil {
			t.Errorf("New returned err: %v", err)
		}

		if got == nil {
			t.Errorf("New is nil")
		}
	}
}

func TestVault_Config_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		config *Config
		err    error
	}{
		{ // valid config with token auth method
			config: &Config{
				Addr:       "https://myvault.com/",
				AuthMethod: vault.TokenAuthMethod,
				Token:      "superSecretAPIKey",
			},
			err: nil,
		},
		{ // valid config with ldap auth method
			config: &Config{
				Addr:       "https://myvault.com/",
				AuthMethod: vault.LDAPAuthMethod,
				Password:   "superSecretPassword",
				Username:   "myusername",
			},
			err: nil,
		},
	}

	// run test
	for _, test := range tests {
		err := test.config.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
