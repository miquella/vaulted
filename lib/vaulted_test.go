package vaulted_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/miquella/vaulted/lib"
	"github.com/miquella/xdg"
)

const (
	VAULT_AAA    = `{"key":{"method":"pbkdf2-sha512","details":{"iterations":93241,"salt":"zWWHn8tQ5YdeWhheqwBWPtPvCf0c3QWwpqq7ircIKRw="}},"method":"secretbox","details":{"nonce":"P9Lhy5gabHJIk7mfQA5jlgHp+Kwa1S2b"},"ciphertext":"jUpP+K05sr+ab5qQR49Qdpnvz71QXncGhT17Qr/A0oiQJ8Bg1p4B"}`
	VAULT_BBB    = `{"key":{"method":"pbkdf2-sha512","details":{"iterations":93391,"salt":"gZRWGwWD8BC+ddVlrnXCgzEsmlvovBUtmwLMN/fqsiQ="}},"method":"secretbox","details":{"nonce":"yF6JxYfO23IjzDsjsLoJ8GnD5kqLQu/L"},"ciphertext":"lHUdCnXyaW1T0OGku00pmS6/bzeXl0WzJmfhZ7nDImfuIPQ6jesS"}`
	VAULT_CCC    = `{"key":{"method":"pbkdf2-sha512","details":{"iterations":93546,"salt":"O658ZVGHXHs1ucsRjQZoaYAYPrjQ9aOKsZdB85tRgwc="}},"method":"secretbox","details":{"nonce":"FrlANjPJRXFtpahvT4j8s63JfTRS+ePY"},"ciphertext":"7rWjYfAkDPu7gJS78dutppo7p+o4yQboYJAZ/1p2Yh3l7R8HpB94"}`
	VAULT_HIDDEN = `{"key":{"method":"pbkdf2-sha512","details":{"iterations":93648,"salt":"09kfyzbAeKYA7SLyoAOHDt3hVjwS4JmTm7pNe6kJ9o0="}},"method":"secretbox","details":{"nonce":"UFKyjfJFLWzLxy2dHu7W0aT3Jbm+I+Ce"},"ciphertext":"ds6Wp3lIdA/GpKsbv5LC0I85tYZuhswORj6a/Vs/l4P6h/EMBsAlvDEZ"}`
)

var (
	xdg_data_home xdg.Path
	xdg_data_dirs xdg.Paths
	xdg_data      xdg.Paths

	dir1 string
	dir2 string
)

func TestListVaults(t *testing.T) {
	err := setupVaults()
	if err != nil {
		t.Fatalf("failted to setup vaults: %v", err)
	}
	defer teardownVaults()

	vaults, err := vaulted.ListVaults()
	if err != nil {
		t.Fatalf("failed to list vaults: %v", err)
	}

	sort.Strings(vaults)
	expected := []string{"aaa", "bbb", "ccc"}
	if !reflect.DeepEqual(expected, vaults) {
		t.Fatalf("expected %#v, got %#v", expected, vaults)
	}
}

func TestOpenVault(t *testing.T) {
	err := setupVaults()
	if err != nil {
		t.Fatalf("failted to setup vaults: %v", err)
	}
	defer teardownVaults()

	vault, err := vaulted.OpenVault("bbb", "password")
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}

	if vault.Vars["TEST"] != "BBB" {
		t.Fatalf("expected: BBB, got %s", vault.Vars["TEST"])
	}
}

func TestSealVault(t *testing.T) {
	err := setupVaults()
	if err != nil {
		t.Fatalf("failed to setup vaults: %v", err)
	}
	defer teardownVaults()

	_, err = vaulted.OpenVault("doesn't exist", "password")
	if err != os.ErrNotExist {
		t.Fatalf("expected: %v, got %v", os.ErrNotExist, err)
	}

	v1 := vaulted.Vault{
		Vars: map[string]string{
			"TEST": "TESTING",
		},
	}
	err = vaulted.SealVault("testing", "another password", &v1)
	if err != nil {
		t.Fatalf("failed to seal vault: %v", err)
	}

	_, err = vaulted.OpenVault("testing", "invalid password")
	if err != vaulted.ErrInvalidPassword {
		t.Fatalf("expected: %v, got: %v", vaulted.ErrInvalidPassword, err)
	}

	v2, err := vaulted.OpenVault("testing", "another password")
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}
	if v2.Vars["TEST"] != "TESTING" {
		t.Fatalf("expected: TESTING, got: %s", v2.Vars["TEST"])
	}
}

func setupVaults() error {
	var err error
	if err == nil {
		dir1, err = ioutil.TempDir("", "vaulted")
	}
	if err == nil {
		err = os.Mkdir(filepath.Join(dir1, "vaulted"), 0700)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir1, "vaulted", "aaa"), []byte(VAULT_AAA), 0600)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir1, "vaulted", "bbb"), []byte(VAULT_BBB), 0600)
	}

	if err == nil {
		dir2, err = ioutil.TempDir("", "vaulted")
	}
	if err == nil {
		err = os.Mkdir(filepath.Join(dir2, "vaulted"), 0700)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir2, "vaulted", "bbb"), []byte(VAULT_HIDDEN), 0600)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir2, "vaulted", "ccc"), []byte(VAULT_CCC), 0600)
	}

	if err == nil {
		xdg_data_home = xdg.DATA_HOME
		xdg.DATA_HOME = xdg.Path(dir1)

		xdg_data_dirs = xdg.DATA_DIRS
		xdg.DATA_DIRS = xdg.Paths{xdg.Path(dir2)}

		xdg_data = xdg.DATA
		xdg.DATA = append(xdg.Paths{xdg.DATA_HOME}, xdg.DATA_DIRS...)
	}

	return err
}

func teardownVaults() {
	xdg.DATA_HOME = xdg_data_home
	xdg.DATA_DIRS = xdg_data_dirs
	xdg.DATA = xdg_data

	if dir1 != "" {
		os.RemoveAll(dir1)
	}
	if dir2 != "" {
		os.RemoveAll(dir2)
	}
}
