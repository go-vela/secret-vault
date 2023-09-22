// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"testing"

	"github.com/go-vela/secret-vault/vault"
)

func TestVault_Plugin_Exec(t *testing.T) {
	// TODO write this test
}

func TestVault_Plugin_Validate(t *testing.T) {
	// setup types
	items, _ := json.Marshal([]Item{
		{
			Path:   "foobar",
			Source: "/path/to/secret",
		},
	})

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
					RawItems: string(items),
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
					RawItems: string(items),
				},
			},
			err: nil,
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
