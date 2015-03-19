package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"hash"
)

const KeySize = 32
const BlockSize = 16

type Vault struct {
	KeyDetails KeyDetails `json:"keyDetails"`

	MACDigest  string `json:"macDigest"`
	Cipher     string `json:"cipher"`
	CipherMode string `json:"cipherMode"`

	MAC          string `json:"mac"`
	IV           string `json:"iv"`
	Environments string `json:"environments"`
}

func (v *Vault) GenerateKey(password string) ([]byte, error) {
	err := v.setDefaults()
	if err != nil {
		return nil, err
	}

	salt, err := base64.StdEncoding.DecodeString(v.KeyDetails.Salt)
	if err != nil {
		return nil, err
	}
	return pbkdf2.Key([]byte(password), salt, v.KeyDetails.Iterations, KeySize, getHashFunc(v.KeyDetails.Digest)), nil
}

func (v *Vault) DecryptEnvironments(key []byte) (Environments, error) {
	// verify mac before decrypting
	mac, err := base64.StdEncoding.DecodeString(v.MAC)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(v.calculateMAC(key), mac) {
		return nil, errors.New("invalid password or corrupted environments")
	}

	// decrypt data
	decrypter, err := v.getCipherDecrypter(key)
	if err != nil {
		return nil, err
	}

	environmentsData, err := base64.StdEncoding.DecodeString(v.Environments)
	if err != nil {
		return nil, err
	}
	decrypter.XORKeyStream(environmentsData, environmentsData)

	// decode json
	envs := Environments{}
	err = json.Unmarshal(environmentsData, &envs)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (v *Vault) EncryptEnvironments(key []byte, envs Environments) error {
	// encode json
	environmentsData, err := json.Marshal(&envs)
	if err != nil {
		return err
	}

	// regenerate iv
	iv := make([]byte, BlockSize)
	_, err = rand.Reader.Read(iv)
	if err != nil {
		return err
	}
	v.IV = base64.StdEncoding.EncodeToString(iv)

	// encrypt data
	encrypter, err := v.getCipherEncrypter(key)
	if err != nil {
		return err
	}

	encrypter.XORKeyStream(environmentsData, environmentsData)
	v.Environments = base64.StdEncoding.EncodeToString(environmentsData)

	// generate mac
	v.MAC = base64.StdEncoding.EncodeToString(v.calculateMAC(key))

	return nil
}

func (v *Vault) setDefaults() error {
	if v.KeyDetails.Digest == "" {
		v.KeyDetails.Digest = "sha-512"
	}

	if v.KeyDetails.Iterations == 0 {
		v.KeyDetails.Iterations = 65536
	}

	if v.KeyDetails.Salt == "" {
		salt := make([]byte, KeySize)
		_, err := rand.Reader.Read(salt)
		if err != nil {
			return err
		}
		v.KeyDetails.Salt = base64.StdEncoding.EncodeToString(salt)
	}

	if v.MACDigest == "" {
		v.MACDigest = "sha-256"
	}

	if v.Cipher == "" {
		v.Cipher = "aes"
	}

	if v.CipherMode == "" {
		v.CipherMode = "ctr"
	}

	return nil
}

func (v *Vault) calculateMAC(key []byte) []byte {
	mac := hmac.New(getHashFunc(v.MACDigest), key)
	mac.Write([]byte(v.Environments))
	return mac.Sum(nil)
}

func (v *Vault) getCipher(key []byte) (cipher.Block, error) {
	switch v.Cipher {
	case "aes":
		return aes.NewCipher(key)
	default:
		return nil, errors.New("invalid cipher")
	}
}

func (v *Vault) getCipherEncrypter(key []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(v.IV)
	if err != nil {
		return nil, err
	}

	switch v.CipherMode {
	case "ctr":
		return cipher.NewCTR(block, iv), nil
	case "ofb":
		return cipher.NewOFB(block, iv), nil
	case "cfb":
		return cipher.NewCFBEncrypter(block, iv), nil
	default:
		return nil, errors.New("invalid cipher mode")
	}
}

func (v *Vault) getCipherDecrypter(key []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(v.IV)
	if err != nil {
		return nil, err
	}

	switch v.CipherMode {
	case "ctr":
		return cipher.NewCTR(block, iv), nil
	case "ofb":
		return cipher.NewOFB(block, iv), nil
	case "cfb":
		return cipher.NewCFBDecrypter(block, iv), nil
	default:
		return nil, errors.New("invalid cipher mode")
	}
}

type KeyDetails struct {
	Digest     string `json:"digest"`
	Salt       string `json:"salt"`
	Iterations int    `json:"iterations"`
}

type Environments map[string]Environment

type Environment struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

func getHashFunc(digest string) func() hash.Hash {
	switch digest {
	case "sha-1":
		return sha1.New
	case "sha-256":
		return sha256.New
	case "sha-512":
		return sha512.New
	default:
		return nil
	}
}
