package vaulted

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/miquella/xdg"
	"golang.org/x/crypto/nacl/secretbox"
)

var (
	ErrInvalidPassword         = errors.New("Invalid password")
	ErrInvalidKeyConfig        = errors.New("Invalid key configuration")
	ErrInvalidEncryptionConfig = errors.New("Invalid encryption configuration")
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

func SealVault(password string, name string, vault *Vault) error {
	vf := &VaultFile{
		Method:  "secretbox",
		Details: make(Details),
	}

	// generate a new key (while trying to keeping the existing key derivation and encryption methods)
	existingVaultFile, err := readVaultFile(name)
	if err == nil {
		vf.Method = existingVaultFile.Method
		vf.Key = existingVaultFile.Key
	}

	vf.Key = newVaultKey(vf.Key)

	// marshal the vault content
	content, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	// encrypt the vault
	if vf.Method == "" {
		vf.Method = "secretbox"
	}

	switch vf.Method {
	case "secretbox":
		nonce := [24]byte{}
		_, err = rand.Read(nonce[:])
		if err != nil {
			return err
		}
		vf.Details.SetBytes("nonce", nonce[:])

		key := [32]byte{}
		derivedKey, err := vf.Key.key(password, len(key))
		if err != nil {
			return err
		}
		copy(key[:], derivedKey[:])

		vf.Ciphertext = secretbox.Seal(nil, content, &nonce, &key)

	default:
		return fmt.Errorf("Invalid encryption method: %s", vf.Method)
	}

	writeVaultFile(name, vf)

	return nil
}

func OpenVault(password string, name string) (*Vault, error) {
	vf, err := readVaultFile(name)
	if err != nil {
		return nil, err
	}

	v := Vault{}

	switch vf.Method {
	case "secretbox":
		if vf.Key == nil {
			return nil, ErrInvalidKeyConfig
		}

		nonce := vf.Details.Bytes("nonce")
		if len(nonce) == 0 {
			return nil, ErrInvalidEncryptionConfig
		}
		boxNonce := [24]byte{}
		copy(boxNonce[:], nonce)

		boxKey := [32]byte{}
		derivedKey, err := vf.Key.key(password, len(boxKey))
		if err != nil {
			return nil, err
		}
		copy(boxKey[:], derivedKey[:])

		plaintext, ok := secretbox.Open(nil, vf.Ciphertext, &boxNonce, &boxKey)
		if !ok {
			return nil, ErrInvalidPassword
		}

		err = json.Unmarshal(plaintext, &v)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Invalid encryption method: %s", vf.Method)
	}

	return &v, nil
}

func readVaultFile(name string) (*VaultFile, error) {
	existing := xdg.DATA.Find(filepath.Join("vaulted", name))
	if len(existing) == 0 {
		return nil, os.ErrNotExist
	}

	f, err := os.Open(existing[0])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	vf := VaultFile{}
	err = d.Decode(&vf)
	if err != nil {
		return nil, err
	}

	return &vf, nil
}

func writeVaultFile(name string, vaultFile *VaultFile) error {
	pathname := xdg.DATA_HOME.Join("vaulted")
	err := os.MkdirAll(pathname, 0700)
	if err != nil {
		return err
	}

	filename := xdg.DATA_HOME.Join(filepath.Join("vaulted", name))
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	err = e.Encode(vaultFile)
	if err != nil {
		return err
	}

	return nil
}
