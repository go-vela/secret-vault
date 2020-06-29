// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestVault_Read_Exec(t *testing.T) {
	// TODO write this test
}

func TestVault_Read_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		read *Read
		err  error
	}{
		{
			// success
			read: &Read{
				Path: "/path/to/secret",
				Keys: []string{"foobar"},
			},
			err: nil,
		},
		{
			// error with no path
			read: &Read{
				Keys: []string{"foobar"},
			},
			err: ErrNoPathProvided,
		},
		{
			// error with no path
			read: &Read{
				Path: "/path/to/secret",
			},
			err: ErrNoKeysProvided,
		},
	}

	// run test
	for _, test := range tests {
		err := test.read.Validate()
		if err != test.err {
			t.Errorf("Validate returned err: %v", err)
		}
	}

}
