package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"io/ioutil"
)

const (
	Iterations        = 65536
	DefaultVaultPerms = 0600
)

type AccountVault struct {
	Accounts map[string]Account `json:"accounts"`
}

// Loads a vault of encrypted accounts
func LoadAccountVault(filename, password string) (*AccountVault, error) {
	vaultData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// parse vault format
	components := bytes.SplitN(vaultData, []byte(":"), 3)
	if len(components) < 3 {
		return nil, errors.New("invalid vault format")
	}

	salt, err := base64.StdEncoding.DecodeString(string(components[0]))
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(string(components[1]))
	if err != nil {
		return nil, err
	}

	// decrypt
	gcm := gcmCipher(salt, password)
	jsonData, err := gcm.Open(nil, nonce, components[2], salt)
	if err != nil {
		return nil, err
	}

	av := AccountVault{}
	err = json.Unmarshal(jsonData, &av)
	if err != nil {
		return nil, err
	}

	return &av, nil
}

// Saves a vault of encrypted accounts
func (av *AccountVault) SaveAccountVault(filename, password string) error {
	jsonData, err := json.Marshal(av)
	if err != nil {
		return err
	}

	// encrypt
	salt := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}
	gcm := gcmCipher(salt, password)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	sealedData := gcm.Seal(nil, nonce, jsonData, salt)

	header := base64.StdEncoding.EncodeToString(salt) + ":" + base64.StdEncoding.EncodeToString(nonce) + ":"
	vaultData := append([]byte(header), sealedData...)
	return ioutil.WriteFile(filename, vaultData, DefaultVaultPerms)
}

type Account struct {
	Name string            `json:"name"`
	Env  map[string]string `json:"env"`
}

// create GCM cipher using account name and password to generate a key (pbkdf2 w/sha256)
func gcmCipher(salt []byte, password string) cipher.AEAD {
	key := pbkdf2.Key([]byte(password), salt, Iterations, 16, sha256.New)
	aesCipher, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(aesCipher)
	return aead
}
