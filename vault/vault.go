// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/sirupsen/logrus"
)

const (
	// LDAPAuthMethod is used for creating a client capable of LDAP authentication.
	LDAPAuthMethod = "ldap"

	// TokenAuthMethod is used for creating a client capable of token authentication.
	TokenAuthMethod = "token"
)

type (
	// client represents an internal struct for managing calls to a Vault instance
	//
	// Vault client docs: https://pkg.go.dev/github.com/hashicorp/vault/api?tab=doc
	Client struct {
		Vault *api.Client
	}

	// Setup represents the configuration necessary for
	// creating a Vault client capable of integrating
	// with a Vault instance.
	Setup struct {
		// specifies the address of the vault instances
		Addr string
		// specifies the authentication method to use
		AuthMethod string
		// specifies the password for authentication with LDAP auth method
		Password string
		// specifies the token for the vault instances
		Token string
		// specifies the username for authentication with LDAP auth method
		Username string
	}
)

var (
	// ErrInvalidAuthMethod defines the error type when the
	// AuthMethod provided to the client is unsupported.
	ErrInvalidAuthMethod = errors.New("invalid auth method provided")

	// LDAPUserPath defines the path the user information gets
	// written to after success LDAP authentication.
	LDAPUserPath = "/auth/ldap/login/%s"
)

// New returns a Secret implementation that integrates with a Vault secrets engine.
func New(s *Setup) (*Client, error) {
	conf := &api.Config{Address: s.Addr}

	// add auth method specific vault client configuration
	switch s.AuthMethod {
	case LDAPAuthMethod:
		logrus.Tracef("creating vault client with %s auth method", LDAPAuthMethod)
		// create Vault client
		vault, err := api.NewClient(conf)
		if err != nil {
			return nil, err
		}

		// options for passing the password
		options := map[string]interface{}{
			"password": s.Password,
		}

		// the login path
		path := fmt.Sprintf(LDAPUserPath, s.Username)

		// call to get a user token
		user, err := vault.Logical().Write(path, options)
		if err != nil {
			return nil, fmt.Errorf("unable to get user token: %w", err)
		}

		// set Vault API token in client
		vault.SetToken(user.Auth.ClientToken)

		return &Client{Vault: vault}, nil
	case TokenAuthMethod:
		logrus.Tracef("creating vault client with %s auth method", TokenAuthMethod)

		// create Vault client
		vault, err := api.NewClient(conf)
		if err != nil {
			return nil, err
		}

		// set Vault API token in client
		vault.SetToken(s.Token)

		return &Client{Vault: vault}, nil
	default:
		return nil, fmt.Errorf("%w: %s (Valid auth methods: %s, %s)",
			ErrInvalidAuthMethod,
			s.AuthMethod,
			LDAPAuthMethod,
			TokenAuthMethod,
		)
	}
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
	c, err := api.NewClient(&api.Config{Address: addr})
	if err != nil {
		t.Fatal(err)
	}

	// set Vault API token in client
	c.SetToken(rootToken)

	return &Client{
		Vault: c,
	}, nil
}
