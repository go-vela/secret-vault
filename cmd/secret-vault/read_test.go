// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/go-vela/secret-vault/vault"
	"github.com/spf13/afero"
)

func TestVault_Read_Exec(t *testing.T) {
	// step types
	vault, _ := vault.NewMock(t)
	source := "secret/foo"
	r := &Read{
		Items: []Item{
			{
				Path:   "foobar",
				Source: source,
			},
		},
	}

	// setup filesystem
	appFS = afero.NewMemMapFs()

	// initialize vault with test data
	vault.Vault.Logical().Write(source, map[string]interface{}{
		"secret": "bar",
	})

	err := r.Exec(vault)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestVault_Read_Validate_success(t *testing.T) {
	// setup types
	tests := []struct {
		read *Read
		err  error
	}{
		{
			// success
			read: &Read{
				Items: []Item{
					{
						Path:   "foobar",
						Source: "/path/to/secret",
					},
				},
			},
			err: nil,
		},
	}

	// run test
	for _, test := range tests {
		err := test.read.Validate()
		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestVault_Read_Validate_failure(t *testing.T) {
	// setup types
	tests := []struct {
		read *Read
		err  error
	}{
		{
			// error with no path
			read: &Read{
				Items: []Item{
					{
						Source: "/path/to/secret",
					},
				},
			},
			err: ErrNoPathProvided,
		},
		{
			// error with no source
			read: &Read{
				Items: []Item{
					{
						Path: "foobar",
					},
				},
			},
			err: ErrNoSourceProvided,
		},
	}

	// run test
	for _, test := range tests {
		err := test.read.Validate()
		if err == nil {
			t.Errorf("Validate should have returned err: %v", err)
		}
	}
}
