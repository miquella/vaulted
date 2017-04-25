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
	xdgBackup xdg.XDG
)

func TestListVaults(t *testing.T) {
	setupVaults(t)
	defer teardownVaults(t)

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
	setupVaults(t)
	defer teardownVaults(t)

	vault, err := vaulted.OpenVault("bbb", "password")
	if err != nil {
		t.Fatalf("failed to open vault: %v", err)
	}

	if vault.Vars["TEST"] != "BBB" {
		t.Fatalf("expected: BBB, got %s", vault.Vars["TEST"])
	}
}

func TestSealVault(t *testing.T) {
	setupVaults(t)
	defer teardownVaults(t)

	_, err := vaulted.OpenVault("doesn't exist", "password")
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

func TestRemoveVault(t *testing.T) {
	setupVaults(t)
	defer teardownVaults(t)

	err := vaulted.RemoveVault("aaa")
	if err != nil {
		t.Fatalf("failed to remove vault: %v", err)
	}

	if _, err := os.Stat(filepath.Join(string(xdg.CACHE_HOME), "vaulted", "aaa")); !os.IsNotExist(err) {
		t.Error("cache for 'aaa' should have been removed and wasn't")
	}
}

func setupVaults(t *testing.T) {
	setupXDG(t)

	// XDG_DATA_HOME
	err := os.Mkdir(filepath.Join(string(xdg.DATA_HOME), "vaulted"), 0700)
	if err != nil {
		t.Fatalf("failed to create vaulted DATA_HOME dir: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(string(xdg.DATA_HOME), "vaulted", "aaa"), []byte(VAULT_AAA), 0600)
	if err != nil {
		t.Fatalf("failed to write 'aaa' home vault file: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(string(xdg.DATA_HOME), "vaulted", "bbb"), []byte(VAULT_BBB), 0600)
	if err != nil {
		t.Fatalf("failed to write 'bbb' home vault file: %v", err)
	}

	// XDG_DATA_DIRS
	err = os.Mkdir(filepath.Join(string(xdg.DATA_DIRS[0]), "vaulted"), 0700)
	if err != nil {
		t.Fatalf("failed to create vaulted DATA_DIRS dir: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(string(xdg.DATA_DIRS[0]), "vaulted", "bbb"), []byte(VAULT_HIDDEN), 0600)
	if err != nil {
		t.Fatalf("failed to write 'bbb' data vault file: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(string(xdg.DATA_DIRS[0]), "vaulted", "ccc"), []byte(VAULT_CCC), 0600)
	if err != nil {
		t.Fatalf("failed to write 'ccc' data vault file: %v", err)
	}

	// XDG_CACHE_HOME
	err = os.Mkdir(filepath.Join(string(xdg.CACHE_HOME), "vaulted"), 0700)
	if err != nil {
		t.Fatalf("failed to create vaulted CACHE_HOME dir: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(string(xdg.CACHE_HOME), "vaulted", "aaa"), []byte{}, 0600)
	if err != nil {
		t.Fatalf("failed to write 'aaa' session file: %v", err)
	}
}

func teardownVaults(t *testing.T) {
	teardownXDG(t)
}

func setupXDG(t *testing.T) {
	xdgBackup.DATA_HOME = xdg.DATA_HOME
	xdgBackup.DATA_DIRS = xdg.DATA_DIRS
	xdgBackup.DATA = xdg.DATA
	xdgBackup.CACHE_HOME = xdg.CACHE_HOME

	// DATA
	data_home, err := ioutil.TempDir("", "vaulted")
	if err != nil {
		t.Fatalf("failed to create XDG_DATA_HOME temp dir: %v", err)
	}
	xdg.DATA_HOME = xdg.Path(data_home)

	data_dirs, err := ioutil.TempDir("", "vaulted")
	if err != nil {
		t.Fatalf("failed to create XDG_DATA_DIRS temp dir: %v", err)
	}
	xdg.DATA_DIRS = xdg.Paths{xdg.Path(data_dirs)}

	xdg.DATA = append(xdg.Paths{xdg.DATA_HOME}, xdg.DATA_DIRS...)

	// CACHE
	cache_home, err := ioutil.TempDir("", "vaulted")
	if err != nil {
		t.Fatalf("failted to create XDG_CACHE_HOME temp dir: %v", err)
	}
	xdg.CACHE_HOME = xdg.Path(cache_home)
}

func teardownXDG(t *testing.T) {
	err := os.RemoveAll(string(xdg.DATA_HOME))
	if err != nil {
		t.Fatalf("failed to remove XDG_DATA_HOME temp dir: %v", err)
	}
	for _, dir := range xdg.DATA_DIRS {
		err := os.RemoveAll(string(dir))
		if err != nil {
			t.Fatalf("failed to remove XDG_DATA_HOME temp dir: %v", err)
		}
	}

	err = os.RemoveAll(string(xdg.CACHE_HOME))
	if err != nil {
		t.Fatalf("failed to remove XDG_CACHE_HOME temp dir: %v", err)
	}

	xdg.DATA_HOME = xdgBackup.DATA_HOME
	xdg.DATA_DIRS = xdgBackup.DATA_DIRS
	xdg.DATA = xdgBackup.DATA
	xdg.CACHE_HOME = xdgBackup.CACHE_HOME
}
