// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/docker"
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

		// vault will return a nil Auth struct with no error if path is correct but password fails
		if user == nil || user.Auth == nil {
			return nil, fmt.Errorf("unable to set user token: authentication failed")
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
func NewMock(t *testing.T) (*Client, *docker.DockerCluster, error) {
	opts := &docker.DockerClusterOptions{
		ImageRepo: "hashicorp/vault", // or "hashicorp/vault-enterprise"
		ImageTag:  "latest",
	}
	cluster := docker.NewTestDockerCluster(t, opts)

	client := cluster.Nodes()[0].APIClient()

	err := client.Sys().Mount("secret", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "1",
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return &Client{Vault: client}, cluster, nil
}
