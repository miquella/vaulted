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
	ErrIncorrectPassword       = errors.New("Incorrect password")
	ErrInvalidKeyConfig        = errors.New("Invalid key configuration")
	ErrInvalidEncryptionConfig = errors.New("Invalid encryption configuration")
)

type Store interface {
	Steward() Steward

	ListVaults() ([]string, error)

	VaultExists(name string) bool
	OpenVault(name string) (*Vault, string, error)
	OpenVaultWithPassword(name, password string) (*Vault, string, error)
	SealVault(vault *Vault, name string) error
	SealVaultWithPassword(vault *Vault, name, password string) error
	RemoveVault(name string) error

	CreateSession(name string) (*Session, string, error)
	GetSession(name string) (*Session, string, error)
}

type store struct {
	steward Steward
}

func New(steward Steward) Store {
	return &store{
		steward: steward,
	}
}

func (s *store) Steward() Steward {
	return s.steward
}

func (s *store) ListVaults() ([]string, error) {
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

func (s *store) VaultExists(name string) bool {
	existing := xdg.DATA.Find(filepath.Join("vaulted", name))
	return len(existing) != 0
}

func (s *store) OpenVault(name string) (*Vault, string, error) {
	if !s.VaultExists(name) {
		return nil, "", os.ErrNotExist
	}

	maxTries := 1
	if getMax, ok := s.steward.(StewardMaxTries); ok {
		maxTries = getMax.GetMaxOpenTries()
	}
	for i := 0; i < maxTries; i++ {
		password, err := s.steward.GetPassword(OpenOperation, name)
		if err != nil {
			return nil, "", err
		}

		if v, p, err := s.OpenVaultWithPassword(name, password); err != ErrIncorrectPassword {
			return v, p, err
		}
	}

	return nil, "", ErrIncorrectPassword
}

func (s *store) OpenVaultWithPassword(name, password string) (*Vault, string, error) {
	if !s.VaultExists(name) {
		return nil, "", os.ErrNotExist
	}

	vf, err := readVaultFile(name)
	if err != nil {
		return nil, "", err
	}

	v := Vault{}

	switch vf.Method {
	case "secretbox":
		if vf.Key == nil {
			return nil, "", ErrInvalidKeyConfig
		}

		nonce := vf.Details.Bytes("nonce")
		if len(nonce) == 0 {
			return nil, "", ErrInvalidEncryptionConfig
		}
		boxNonce := [24]byte{}
		copy(boxNonce[:], nonce)

		boxKey := [32]byte{}
		derivedKey, err := vf.Key.key(password, len(boxKey))
		if err != nil {
			return nil, "", err
		}
		copy(boxKey[:], derivedKey[:])

		plaintext, ok := secretbox.Open(nil, vf.Ciphertext, &boxNonce, &boxKey)
		if !ok {
			return nil, "", ErrIncorrectPassword
		}

		err = json.Unmarshal(plaintext, &v)
		if err != nil {
			return nil, "", err
		}

	default:
		return nil, "", fmt.Errorf("Invalid encryption method: %s", vf.Method)
	}

	return &v, password, nil
}

func (s *store) SealVault(vault *Vault, name string) error {
	password, err := s.steward.GetPassword(SealOperation, name)
	if err != nil {
		return err
	}

	return s.SealVaultWithPassword(vault, name, password)
}

func (s *store) SealVaultWithPassword(vault *Vault, name, password string) error {
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

func (s *store) RemoveVault(name string) error {
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

func (s *store) CreateSession(name string) (*Session, string, error) {
	v, password, err := s.OpenVault(name)
	if err != nil {
		return nil, "", err
	}

	session, err := s.createSession(v, name, password)
	if err != nil {
		return nil, "", err
	}

	if v.AWSKey != nil && v.AWSKey.Role != "" {
		session, err = session.Assume(v.AWSKey.Role)
		if err != nil {
			return nil, "", err
		}
	}

	return session, "", nil
}

func (s *store) GetSession(name string) (*Session, string, error) {
	v, password, err := s.OpenVault(name)
	if err != nil {
		return nil, "", err
	}

	session, err := s.getSession(v, name, password)
	if err != nil {
		return nil, "", err
	}

	if v.AWSKey != nil && v.AWSKey.Role != "" {
		session, err = session.Assume(v.AWSKey.Role)
		if err != nil {
			return nil, "", err
		}
	}

	return session, "", nil
}

func (s *store) getSession(v *Vault, name, password string) (*Session, error) {
	session, err := s.openSession(name, password)
	if err != nil {
		removeSession(name)
	} else if session.Expiration.After(time.Now().Add(15 * time.Minute)) {
		return session, nil
	}

	return s.createSession(v, name, password)
}

func (s *store) createSession(v *Vault, name, password string) (*Session, error) {
	var session *Session
	var err error
	if v.AWSKey.RequiresMFA() {
		var mfaToken string
		mfaToken, err = s.steward.GetMFAToken(name)
		if err == nil {
			session, err = v.CreateSessionWithMFA(name, mfaToken)
		}
	} else {
		session, err = v.CreateSession(name)
	}
	if err != nil {
		return nil, err
	}

	// we have a valid session, so if saving fails, ignore the failure
	s.sealSession(session, name, password)
	return session, nil
}

func (s *store) sealSession(session *Session, name, password string) error {
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

func (s *store) openSession(name, password string) (*Session, error) {
	vf, err := readVaultFile(name)
	if err != nil {
		return nil, err
	}

	sf, err := readSessionFile(name)
	if err != nil {
		return nil, err
	}

	session := Session{}

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
			return nil, ErrIncorrectPassword
		}

		err = json.Unmarshal(plaintext, &session)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Invalid encryption method: %s", sf.Method)
	}

	return &session, nil
}
