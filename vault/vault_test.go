// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVault_New(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if s == nil {
		t.Error("New returned nil client")
	}
}

func TestVault_New_Error(t *testing.T) {
	// run test
	s, err := New("!@#$%^&*()", "")
	if err == nil {
		t.Errorf("New should have returned err")
	}

	if s != nil {
		t.Error("New should have returned nil client")
	}
}
