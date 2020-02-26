package main

import (
	"fmt"
	"sort"

	"github.com/miquella/vaulted/v3/lib"
)

type List struct {
	Active string
}

func (l *List) Run(store vaulted.Store) error {
	vaults, err := store.ListVaults()
	if err != nil {
		return err
	}

	sort.Strings(vaults)
	for _, vault := range vaults {
		if vault == l.Active {
			vault = fmt.Sprintf("%s (active)", vault)
		}
		fmt.Println(vault)
	}

	return nil
}
