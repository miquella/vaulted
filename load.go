package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/miquella/vaulted/lib"
)

type Load struct {
	VaultName string
}

func (l Load) Run(steward Steward) error {
	jvault, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	vault := &vaulted.Vault{}
	err = json.Unmarshal(jvault, vault)
	if err != nil {
		return err
	}

	err = steward.SealVault(l.VaultName, nil, vault)
	if err != nil {
		return err
	}

	return nil
}
