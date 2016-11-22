package main

import (
	"fmt"
	"sort"
)

type List struct {
	Active string
}

func (l *List) Run(steward Steward) error {
	vaults, err := steward.ListVaults()
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
