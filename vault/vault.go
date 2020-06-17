// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

type Client struct {
	// https://pkg.go.dev/github.com/hashicorp/vault/api?tab=doc
	Vault *api.Client
}

// New returns a Secret implementation that integrates with a Vault secrets engine.
func New(addr, token string) (*Client, error) {
	conf := api.Config{Address: addr}

	// create Vault client
	c, err := api.NewClient(&conf)
	if err != nil {
		return nil, err
	}

	// set Vault API token in client
	c.SetToken(token)

	// return client
	return &Client{
		Vault: c,
	}, nil
}

// NewMock returns a test unsealed Vault
// to integrate with a Vault secret provider.
//
// This function is intended for running tests only.
//
// Docs: https://pkg.go.dev/github.com/hashicorp/vault/vault?tab=doc
func NewMock(t *testing.T) (*Client, error) {
	// Create an in-memory, unsealed core (the "backend", if you will).
	//
	// Pinned commit version in Go.sum:
	// https://github.com/hashicorp/vault/issues/9072
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	_, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	c, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	c.SetToken(rootToken)

	return &Client{
		Vault: c,
	}, nil
}
