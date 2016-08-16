package main

import (
	"fmt"
	"sort"
)

type List struct{}

func (l *List) Run(steward Steward) error {
	vaults, err := steward.ListVaults()
	if err != nil {
		return err
	}

	sort.Strings(vaults)
	for _, vault := range vaults {
		fmt.Println(vault)
	}

	return nil
}
