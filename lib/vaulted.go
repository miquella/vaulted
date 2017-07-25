package vaulted

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/miquella/xdg"
	"golang.org/x/crypto/nacl/secretbox"
)

var (
	ErrInvalidPassword         = errors.New("Invalid password")
	ErrInvalidKeyConfig        = errors.New("Invalid key configuration")
	ErrInvalidEncryptionConfig = errors.New("Invalid encryption configuration")
)

func VaultExists(name string) bool {
	existing := xdg.DATA.Find(filepath.Join("vaulted", name))
	if len(existing) == 0 {
		return false
	}

	return true
}

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

func SealVault(name, password string, vault *Vault) error {
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

	return writeVaultFile(name, vf)
}

func OpenVault(name, password string) (*Vault, error) {
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

func RemoveVault(name string) error {
	existing := xdg.DATA_HOME.Find(filepath.Join("vaulted", name))
	if existing == "" {
		untouchable := xdg.DATA_DIRS.Find(filepath.Join("vaulted", name))
		if len(untouchable) == 0 {
			return os.ErrNotExist
		}

		return fmt.Errorf("Because %s is outside the vaulted managed directory (%s), it must be removed manually", untouchable[0], xdg.DATA_HOME.Join("vaulted"))
	}

	removeSession(name)

	return os.Remove(existing)
}

func GetSession(name, password string) (*Session, error) {
	v, err := OpenVault(name, password)
	if err != nil {
		return nil, err
	}

	session, err := getSession(v, name, password)
	if err != nil {
		return nil, err
	}

	if v.AWSKey != nil && v.AWSKey.Role != "" {
		session, err = session.Assume(v.AWSKey.Role)
		if err != nil {
			return nil, err
		}
	}

	return session, nil
}

func getSession(v *Vault, name, password string) (*Session, error) {
	session, err := openSession(name, password)
	if err != nil {
		removeSession(name)
	} else if session.Expiration.After(time.Now().Add(15 * time.Minute)) {
		return session, nil
	}

	session, err = v.CreateSession(name)
	if err != nil {
		return nil, err
	}

	// we have a valid session, so if saving fails, ignore the failure
	sealSession(name, password, session)
	return session, nil
}

func sealSession(name, password string, session *Session) error {
	// read the vault file (to get key details)
	vf, err := readVaultFile(name)
	if err != nil {
		return err
	}

	// marshal the session content
	content, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// encrypt the session
	sf := &SessionFile{
		Method:  "secretbox",
		Details: make(Details),
	}

	switch sf.Method {
	case "secretbox":
		nonce := [24]byte{}
		_, err = rand.Read(nonce[:])
		if err != nil {
			return err
		}
		sf.Details.SetBytes("nonce", nonce[:])

		key := [32]byte{}
		derivedKey, err := vf.Key.key(password, len(key))
		if err != nil {
			return err
		}
		copy(key[:], derivedKey[:])

		sf.Ciphertext = secretbox.Seal(nil, content, &nonce, &key)

	default:
		return err
	}

	return writeSessionFile(name, sf)
}

func openSession(name, password string) (*Session, error) {
	vf, err := readVaultFile(name)
	if err != nil {
		return nil, err
	}

	sf, err := readSessionFile(name)
	if err != nil {
		return nil, err
	}

	s := Session{}

	switch sf.Method {
	case "secretbox":
		if vf.Key == nil {
			return nil, ErrInvalidKeyConfig
		}

		nonce := sf.Details.Bytes("nonce")
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

		plaintext, ok := secretbox.Open(nil, sf.Ciphertext, &boxNonce, &boxKey)
		if !ok {
			return nil, ErrInvalidPassword
		}

		err = json.Unmarshal(plaintext, &s)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Invalid encryption method: %s", sf.Method)
	}

	return &s, nil
}
