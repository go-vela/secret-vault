// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Read is a function to capture
// the secret for the provided path.
func (c *Client) Read(path string) (*api.Secret, error) {
	// send API call to capture the secret
	vault, err := c.Vault.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve secret %s: %w", path, err)
	}

	// return nil if secret does not exist
	if vault == nil {
		return nil, fmt.Errorf("unable to retrieve secret %s", path)
	}

	return vault, nil
}
