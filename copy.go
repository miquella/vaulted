package main

type Copy struct {
	OldVaultName string
	NewVaultName string
}

func (c *Copy) Run(steward Steward) error {
	_, vault, err := steward.OpenVault(c.OldVaultName, nil)
	if err != nil {
		return err
	}

	err = steward.SealVault(c.NewVaultName, nil, vault)
	if err != nil {
		return err
	}

	return nil
}
