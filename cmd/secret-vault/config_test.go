// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestVault_Config_New(t *testing.T) {
	// setup types
	c := &Config{
		Addr:  "https://myvault.com/",
		Token: "superSecretAPIKey",
	}

	got, err := c.New()
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got == nil {
		t.Errorf("New is nil")
	}
}

func TestVault_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Addr:  "https://myvault.com/",
		Token: "superSecretAPIKey",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
