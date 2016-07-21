package legacy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrInvalidPassword         = errors.New("Invalid password")
	ErrInvalidKeyConfig        = errors.New("Invalid key configuration")
	ErrInvalidEncryptionConfig = errors.New("Invalid encryption configuration")
)

type Vault struct {
	KeyDetails KeyDetails `json:"keyDetails"`

	MACDigest  string `json:"macDigest"`
	Cipher     string `json:"cipher"`
	CipherMode string `json:"cipherMode"`

	MAC          []byte `json:mac"`
	IV           []byte `json:"iv"`
	Environments []byte `json:"environments"`
}

type KeyDetails struct {
	Digest     string `json:"digest"`
	Iterations int    `json:"iterations"`
	Salt       []byte `json:"salt"`
}

type Environment struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

func ReadVault() (*Vault, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(u.HomeDir, ".vaulted")
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	v := &Vault{}
	d := json.NewDecoder(f)
	err = d.Decode(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vault) DecryptEnvironments(password string) (map[string]Environment, error) {
	if v.MACDigest != "sha-256" && v.Cipher != "aes" && v.CipherMode != "ctr" {
		return nil, ErrInvalidEncryptionConfig
	}

	// derive the key
	key, err := v.KeyDetails.key(password)
	if err != nil {
		return nil, err
	}

	// validate the mac
	if !hmac.Equal(v.mac(key), v.MAC) {
		return nil, ErrInvalidPassword
	}

	// decrypt the environments
	plaintext := make([]byte, len(v.Environments))

	block, err := aes.NewCipher(key)
	decrypter := cipher.NewCTR(block, v.IV)
	decrypter.XORKeyStream(plaintext, v.Environments)

	// unmarshal the environments
	environments := map[string]Environment{}
	err = json.Unmarshal(plaintext, &environments)
	if err != nil {
		return nil, err
	}

	return environments, nil
}

func (v *Vault) mac(key []byte) []byte {
	encodedEnvironments := make([]byte, base64.StdEncoding.EncodedLen(len(v.Environments)))
	base64.StdEncoding.Encode(encodedEnvironments, v.Environments)

	mac := hmac.New(sha256.New, key)
	mac.Write(encodedEnvironments)
	return mac.Sum(nil)
}

func (kd *KeyDetails) key(password string) ([]byte, error) {
	if kd.Digest != "sha-512" {
		return nil, ErrInvalidKeyConfig
	}

	return pbkdf2.Key([]byte(password), kd.Salt, kd.Iterations, 32, sha512.New), nil
}
