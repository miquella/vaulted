package main

import (
	"fmt"
	"sort"
)

type List struct {
	Active []string
}

func (l *List) Run(steward Steward) error {
	vaults, err := steward.ListVaults()
	if err != nil {
		return err
	}

	active := map[string]bool{}
	for _, name := range l.Active {
		active[name] = true
	}

	sort.Strings(vaults)
	for _, vault := range vaults {
		if active[vault] {
			vault = fmt.Sprintf("%s (active)", vault)
		}
		fmt.Println(vault)
	}

	return nil
}
