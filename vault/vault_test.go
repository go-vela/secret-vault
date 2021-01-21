// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVault_New(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	vault, _ := NewMock(t)

	// setup types
	tests := []struct {
		setup *Setup
		err   error
	}{
		{ // Success with token auth method
			setup: &Setup{
				Addr:       vault.Vault.Address(),
				AuthMethod: TokenAuthMethod,
				Token:      "supersecrettoken",
			},
			err: nil,
		},
		// TODO: investigate how to put mock vault in a mode with a fake LDAP provider
		// { // Success with token auth method
		// 	setup: &Setup{
		// 		Addr:       vault.Vault.Address(),
		// 		AuthMethod: LDAPAuthMethod,
		// 		Password:   "superSecretPassword",
		// 		Username:   "myusername",
		// 	},
		// 	err: nil,
		// },
	}

	// run test
	for _, test := range tests {
		vault, err := New(test.setup)
		if err != test.err {
			t.Errorf("New returned err: %v", err)
		}

		if vault == nil {
			t.Error("New returned nil client")
		}
	}
}

func TestVault_New_Error(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// setup types
	tests := []struct {
		setup *Setup
		err   error
	}{
		{ // failure with bad address and fake auth method
			setup: &Setup{
				Addr:       "!@#$%^&*()",
				AuthMethod: "fake",
				Token:      "",
			},
			err: fmt.Errorf("invalid auth method provided: fake (Valid auth methods: ldap, token)"),
		},
		{ // failure with no address
			setup: &Setup{
				AuthMethod: "fake",
				Token:      "",
			},
			err: fmt.Errorf("invalid auth method provided: fake (Valid auth methods: ldap, token)"),
		},
		{ // failure with no auth method
			setup: &Setup{
				Addr:  "!@#$%^&*()",
				Token: "",
			},
			err: fmt.Errorf("invalid auth method provided: fake (Valid auth methods: ldap, token)"),
		},
	}

	// run test
	for _, test := range tests {
		_, err := New(test.setup)
		if err == test.err {
			t.Errorf("New returned err: %v", err)
		}
	}
}
