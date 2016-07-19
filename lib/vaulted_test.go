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

func setupVaults() error {
	var err error
	if err == nil {
		dir1, err = ioutil.TempDir("", "vaulted")
	}
	if err == nil {
		err = os.Mkdir(filepath.Join(dir1, "vaulted"), 0700)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir1, "vaulted", "aaa"), []byte{}, 0600)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir1, "vaulted", "bbb"), []byte{}, 0600)
	}

	if err == nil {
		dir2, err = ioutil.TempDir("", "vaulted")
	}
	if err == nil {
		err = os.Mkdir(filepath.Join(dir2, "vaulted"), 0700)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir2, "vaulted", "bbb"), []byte{}, 0600)
	}
	if err == nil {
		err = ioutil.WriteFile(filepath.Join(dir2, "vaulted", "ccc"), []byte{}, 0600)
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
