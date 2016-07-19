package vaulted

import (
	"os"
	"path/filepath"

	"github.com/miquella/xdg"
)

func ListVaults() ([]string, error) {
	vaults, err := xdg.DATA.Glob(filepath.Join("vaulted", "*"))
	if err != nil {
		return nil, err
	}

	var found []string
	emitted := map[string]bool{}
	for _, vault := range vaults {
		info, err := os.Stat(vault)
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			continue
		}

		if !emitted[info.Name()] {
			emitted[info.Name()] = true
			found = append(found, info.Name())
		}
	}

	return found, nil
}
