// SPDX-License-Identifier: Apache-2.0

package main

import (
	"reflect"
	"testing"

	"github.com/go-vela/secret-vault/vault"
	"github.com/spf13/afero"
)

func TestVault_Read_Exec(t *testing.T) {
	// step types
	vault, _ := vault.NewMock(t)
	source := "/secret/foo"
	path := []string{"foobar", "foobar2"}
	keys := map[string]string{
		"secret": "my_secret",
	}

	r := &Read{
		Items: []*Item{
			{
				Path:   path,
				Source: source,
				Keys:   keys,
			},
		},
	}

	// setup filesystem
	appFS = afero.NewMemMapFs()

	// initialize vault with test data
	//nolint: errcheck // error check not needed
	vault.Vault.Logical().Write(source, map[string]interface{}{
		"secret": "bar",
	})

	err := r.Exec(vault)
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
	}
}

func TestVault_Read_Exec_Fail(t *testing.T) {
	// step types
	vault, _ := vault.NewMock(t)
	source := ""
	path := []string{"foobar", "foobar2"}
	r := &Read{
		Items: []*Item{
			{
				Path:   path,
				Source: source,
			},
		},
	}

	// setup filesystem
	appFS = afero.NewMemMapFs()

	// initialize vault with test data
	//nolint: errcheck // error check not needed
	vault.Vault.Logical().Write(source, map[string]interface{}{
		"secret": "bar",
	})

	err := r.Exec(vault)
	if err == nil {
		t.Errorf("Exec should have returned err: %v", err)
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
				Items: []*Item{
					{
						Path:   []string{"foobar", "foobar2"},
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
			// error with no items
			read: &Read{
				Items: []*Item{},
			},
			err: ErrNoItemsProvided,
		},
		{
			// error with no path
			read: &Read{
				Items: []*Item{
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
				Items: []*Item{
					{
						Path: []string{"foobar", "foobar2"},
					},
				},
			},
			err: ErrNoSourceProvided,
		},
		{
			// error with nil path
			read: &Read{
				Items: []*Item{
					{
						Source: "/path/to/secret",
						Path:   []string{""},
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

func TestVault_Read_Unmarshal(t *testing.T) {
	// setup types
	r := &Read{
		RawItems: `
		[
			{"path":["foo", "foo2"],"source":"secret/vela/hello_world","keys":{"foo":"bar"}}
		]
		`}

	want := []*Item{
		{
			Path:   []string{"foo", "foo2"},
			Source: "secret/vela/hello_world",
			Keys: map[string]string{
				"foo": "bar",
			},
		},
	}

	err := r.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if !reflect.DeepEqual(r.Items, want) {
		t.Errorf("Unmarshal is %v, want %v", r.Items, want)
	}
}

func TestVault_Read_Unmarshal_Single_Path(t *testing.T) {
	// setup types
	r := &Read{
		RawItems: `
		[
			{"path":"foo","source":"secret/vela/hello_world"}
		]
		`}

	want := []*Item{
		{
			Path:   []string{"foo"},
			Source: "secret/vela/hello_world",
		},
	}

	err := r.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if !reflect.DeepEqual(r.Items, want) {
		t.Errorf("Unmarshal is %v, want %v", r.Items, want)
	}
}

func TestVault_Read_Unmarshal_Fail(t *testing.T) {
	// setup types
	r := &Read{
		RawItems: `
		[
			{"path":"foo,"source":"secret/vela/hello_world"}
		]
		`}

	err := r.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err: %v", err)
	}
}

func TestVault_Read_Unmarshal_Fail_BadKeyMap(t *testing.T) {
	// setup types
	r := &Read{
		RawItems: `
		[
			{"path":["foo", "foo2"],"source":"secret/vela/hello_world","keys":["foo=bar"]}
		]
		`}

	err := r.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err: %v", err)
	}
}
