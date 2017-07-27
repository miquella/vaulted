package legacy

import (
	"github.com/miquella/vaulted/lib"
)

type LegacyStore interface {
	OpenLegacyVault() (map[string]Environment, string, error)
}

type legacyStore struct {
	steward vaulted.Steward
}

func New(steward vaulted.Steward) LegacyStore {
	return &legacyStore{
		steward: steward,
	}
}

func (s *legacyStore) OpenLegacyVault() (map[string]Environment, string, error) {
	maxTries := 1
	if getMax, ok := s.steward.(vaulted.StewardMaxTries); ok {
		maxTries = getMax.GetMaxOpenTries()
	}
	for i := 0; i < maxTries; i++ {
		password, err := s.steward.GetPassword(LegacyOperation, "Legacy vaults")
		if err != nil {
			return nil, "", err
		}

		if v, p, err := s.OpenLegacyVaultWithPassword(password); err != vaulted.ErrInvalidPassword {
			return v, p, err
		}
	}

	return nil, "", vaulted.ErrInvalidPassword
}

func (s *legacyStore) OpenLegacyVaultWithPassword(password string) (map[string]Environment, string, error) {
	legacyVault, err := ReadVault()
	if err != nil {
		return nil, "", err
	}

	environments, err := legacyVault.DecryptEnvironments(password)
	return environments, password, err
}
