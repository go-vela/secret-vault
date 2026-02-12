// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-envparse"
	"github.com/spf13/afero"

	"github.com/go-vela/secret-vault/vault"
	"github.com/go-vela/server/compiler/types/raw"
)

func TestVault_Read_Exec_Legacy(t *testing.T) {
	// step types
	vault, cluster, _ := vault.NewMock(t)
	defer cluster.Cleanup()

	source := "/secret/foo"
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
		"secret":             "bar",
		"dash-secret":        "baz",
		"crazy??//!.#secret": "bazzy",
	})

	err := r.Exec(vault)
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
	}

	t.Setenv("VELA_MASKED_BASE64_OUTPUTS", "/vela/outputs/masked.env")

	err = appFS.MkdirAll(filepath.Dir("/vela/outputs/masked.env"), 0777)
	if err != nil {
		t.Error(err)
	}

	r.OutputsPath = "/vela/outputs/masked.env"

	err = r.Exec(vault)
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
	}

	a := &afero.Afero{
		Fs: appFS,
	}

	rawOutputs, err := a.ReadFile("/vela/outputs/masked.env")
	if err != nil {
		t.Errorf("unable to read outputs file: %v", err)
	}

	envMap, err := envparse.Parse(bytes.NewReader(rawOutputs))
	if err != nil {
		t.Errorf("unable to parse outputs file: %v", err)
	}

	for k, v := range envMap {
		// decode the base64 value
		decodedValue, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			t.Errorf("unable to decode base64 value for key %s: %v", k, err)
		}

		envMap[k] = string(decodedValue)
	}

	if envMap["VELA_SECRETS_FOOBAR_SECRET"] != "bar" {
		t.Errorf("Exec is %v, want %v", envMap["VELA_SECRETS_FOOBAR_SECRET"], "bar")
	}

	if envMap["VELA_SECRETS_FOOBAR_DASH_SECRET"] != "baz" {
		t.Errorf("Exec is %v, want %v", envMap["VELA_SECRETS_FOOBAR_DASH_SECRET"], "baz")
	}

	if envMap["VELA_SECRETS_FOOBAR_CRAZY_______SECRET"] != "bazzy" {
		t.Errorf("Exec is %v, want %v", envMap["VELA_SECRETS_FOOBAR_CRAZY_______SECRET"], "bazzy")
	}
}

func TestVault_Read_Exec(t *testing.T) {
	// step types
	vault, cluster, _ := vault.NewMock(t)
	defer cluster.Cleanup()

	source := "/secret/foo"

	r := &Read{
		Items: []*Item{
			{
				Source: source,
				Keys: map[string]KeyItem{
					"secret": {
						Name:   "secret",
						Target: raw.StringSlice{"TEST_SECRET", "TEST_SECRET_COPY"},
					},
					"dash-secret": {
						Name:   "dash-secret",
						Target: raw.StringSlice{"TEST_DASH_SECRET"},
						Path:   raw.StringSlice{"custom"},
					},
					"crazy??//!.#secret": {
						Name:   "crazy??//!.#secret",
						Target: raw.StringSlice{"TEST_CRAZY_SECRET"},
					},
				},
			},
		},
	}

	// setup filesystem
	appFS = afero.NewMemMapFs()

	// initialize vault with test data
	//nolint: errcheck // error check not needed
	vault.Vault.Logical().Write(source, map[string]interface{}{
		"secret":             "bar",
		"dash-secret":        "baz",
		"crazy??//!.#secret": "bazzy",
	})

	err := r.Exec(vault)
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
	}

	t.Setenv("VELA_MASKED_BASE64_OUTPUTS", "/vela/outputs/masked.env")

	err = appFS.MkdirAll(filepath.Dir("/vela/outputs/masked.env"), 0777)
	if err != nil {
		t.Error(err)
	}

	r.OutputsPath = "/vela/outputs/masked.env"

	err = r.Exec(vault)
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
	}

	a := &afero.Afero{
		Fs: appFS,
	}

	rawOutputs, err := a.ReadFile("/vela/outputs/masked.env")
	if err != nil {
		t.Errorf("unable to read outputs file: %v", err)
	}

	envMap, err := envparse.Parse(bytes.NewReader(rawOutputs))
	if err != nil {
		t.Errorf("unable to parse outputs file: %v", err)
	}

	for k, v := range envMap {
		// decode the base64 value
		decodedValue, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			t.Errorf("unable to decode base64 value for key %s: %v", k, err)
		}

		envMap[k] = string(decodedValue)
	}

	if envMap["TEST_SECRET"] != "bar" && envMap["TEST_SECRET_COPY"] != "bar" {
		t.Errorf("Exec is %v, want %v", envMap["TEST_SECRET"], "bar")
	}

	if envMap["TEST_DASH_SECRET"] != "baz" {
		t.Errorf("Exec is %v, want %v", envMap["TEST_DASH_SECRET"], "baz")
	}

	if envMap["TEST_CRAZY_SECRET"] != "bazzy" {
		t.Errorf("Exec is %v, want %v", envMap["TEST_CRAZY_SECRET"], "bazzy")
	}
}

func TestVault_Read_Exec_Fail(t *testing.T) {
	// step types
	vault, cluster, _ := vault.NewMock(t)
	defer cluster.Cleanup()

	source := "secret"
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

func TestVault_Read_Unmarshal_Legacy(t *testing.T) {
	items, err := os.ReadFile("testdata/legacy.json")
	if err != nil {
		t.Errorf("unable to read test items file: %v", err)
	}

	// setup types
	r := &Read{
		RawItems: string(items),
	}

	want := []*Item{
		{
			Source: "secret/team/database",
			Path:   raw.StringSlice{"database"},
		},
		{
			Source: "secret/team/nua",
			Path:   raw.StringSlice{"docker", "artifactory"},
		},
	}

	err = r.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if diff := cmp.Diff(want, r.Items); diff != "" {
		t.Errorf("Unmarshal mismatch (-want +got):\n%s", diff)
	}
}

func TestVault_Read_Unmarshal(t *testing.T) {
	items, err := os.ReadFile("testdata/items.json")
	if err != nil {
		t.Errorf("unable to read test items file: %v", err)
	}

	// setup types
	r := &Read{
		RawItems: string(items),
	}

	want := []*Item{
		{
			Source: "secret/team/database",
			Keys: map[string]KeyItem{
				"connection": {
					Name:   "connection",
					Target: raw.StringSlice{"DB_CONNECTION"},
				},
			},
		},
		{
			Source: "secret/team/nua",
			Keys: map[string]KeyItem{
				"username": {
					Name:   "username",
					Target: raw.StringSlice{"DOCKER_USERNAME", "ARTIFACTORY_USERNAME"},
					Path:   raw.StringSlice{"docker/username", "artifactory/username"},
				},
				"password": {
					Name:   "password",
					Target: raw.StringSlice{"DOCKER_PASSWORD", "ARTIFACTORY_PASSWORD"},
					Path:   raw.StringSlice{"docker/password", "artifactory/password"},
				},
			},
		},
	}

	err = r.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if diff := cmp.Diff(want, r.Items); diff != "" {
		t.Errorf("Unmarshal mismatch (-want +got):\n%s", diff)
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
