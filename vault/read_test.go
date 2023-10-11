// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestVault_Read(t *testing.T) {
	// step types
	vault, _ := NewMock(t)
	path := "secret/foo"
	want := &api.Secret{
		Data: map[string]interface{}{
			"secret": "bar",
		},
	}

	// initialize vault with test data
	_, _ = vault.Vault.Logical().Write("secret/foo", map[string]interface{}{
		"secret": "bar",
	})

	// run
	got, err := vault.Read(path)
	if err != nil {
		t.Errorf("Count returned err: %v", err)
	}

	if !reflect.DeepEqual(got.Data, want.Data) {
		t.Errorf("Read is %+v, want %+v", got, want)
	}
}
