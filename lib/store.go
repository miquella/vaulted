package vaulted

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/miquella/ssh-proxy-agent/lib/proxyagent"
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

	CreateSession(vault *Vault, name, password string) (*Session, error)
	GetSession(vault *Vault, name, password string) (*Session, error)
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

func (s *store) listRecursive(extraParts string) ([]string, error) {
	var vaults []string
	var err error
	if extraParts == "" {
		vaults, err = xdg.DATA.Glob(filepath.Join("vaulted", "*"))
	} else {
		vaults, err = xdg.DATA.Glob(filepath.Join("vaulted", extraParts, "*"))
	}
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
		if info.Mode().IsDir() {
			var innerVaults []string
			if extraParts == "" {
				innerVaults, err = s.listRecursive(info.Name())
			} else {
				innerVaults, err = s.listRecursive(filepath.Join(extraParts, info.Name()))
			}

			if err != nil {
				return nil, err
			}
			found = append(found, innerVaults...)
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}

		if !emitted[info.Name()] {
			emitted[info.Name()] = true
			if extraParts == "" {
				found = append(found, info.Name())
			} else {
				found = append(found, extraParts+"/"+info.Name())
			}
		}
	}

	return found, nil
}

func (s *store) ListVaults() ([]string, error) {
	return s.listRecursive("")
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

	removeSessionCache(name)

	return os.Remove(existing)
}

func (s *store) GetSession(v *Vault, name, password string) (*Session, error) {
	sessionCache, err := s.openSessionCache(name, password)
	if err != nil {
		sessionCache = &SessionCache{}
		removeSessionCache(name)
	} else {
		session, err := sessionCache.GetVaultSession(v)
		if err == nil && !session.Expired(15*time.Minute) {
			return session, nil
		}
	}

	return s.CreateSession(v, name, password)
}

func (s *store) CreateSession(v *Vault, name, password string) (*Session, error) {
	var session *Session
	var err error

	if !s.VaultExists(name) {
		return nil, os.ErrNotExist
	}

	// actually create the session
	if v.AWSKey.RequiresMFA() {
		var mfaToken string
		mfaToken, err = s.steward.GetMFAToken(name)
		if err == nil {
			session, err = v.NewSessionWithMFA(name, mfaToken)
		}
	} else {
		session, err = v.NewSession(name)
	}
	if err != nil {
		return nil, err
	}

	// create a fresh generated key if we are not using a cached session
	if v.SSHOptions != nil && v.SSHOptions.GenerateRSAKey {
		var keyPair *proxyagent.KeyPair
		keyPair, err = proxyagent.GenerateRSAKeyPair()
		if err != nil {
			return nil, err
		}
		session.GeneratedSSHKey = keyPair.PrivateKey
	}

	// we ignore errors because the session is viable even if saving the cache fails
	sessionCache, err := s.openSessionCache(name, password)
	if err != nil {
		sessionCache = &SessionCache{}
		removeSessionCache(name)
	}

	sessionCache.PutVaultSession(v, session)
	s.sealSessionCache(sessionCache, name, password)

	return session, nil
}

func (s *store) sealSessionCache(sessionCache *SessionCache, name, password string) error {
	// read the vault file (to get key details)
	vf, err := readVaultFile(name)
	if err != nil {
		return err
	}

	// normalize the session cache first
	sessionCache.SessionCacheVersion = SessionCacheVersion
	sessionCache.RemoveExpiredSessions()

	// marshal the session cache content
	content, err := json.Marshal(sessionCache)
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

func (s *store) openSessionCache(name, password string) (*SessionCache, error) {
	vf, err := readVaultFile(name)
	if err != nil {
		return nil, err
	}

	sf, err := readSessionFile(name)
	if err != nil {
		return nil, err
	}

	sessionCache := SessionCache{}

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

		err = json.Unmarshal(plaintext, &sessionCache)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Invalid encryption method: %s", sf.Method)
	}

	if sessionCache.SessionCacheVersion != SessionCacheVersion {
		return nil, fmt.Errorf("Invalid session version: %s", sessionCache.SessionCacheVersion)
	}

	return &sessionCache, nil
}
