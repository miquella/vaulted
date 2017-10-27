package main

import (
	"github.com/keybase/go-keychain"
	"github.com/miquella/vaulted/lib"
)

const (
	keychainPath = "vaulted.keychain"
)

func (t *AskPassSteward) GetKeychainPassword(name string) (string, error) {
	return getKeychainPassword(name)
}

func (t *AskPassSteward) SetKeychainPassword(name, password string) error {
	return setKeychainPassword(name, password)
}

func (t *TTYSteward) GetKeychainPassword(name string) (string, error) {
	return getKeychainPassword(name)
}

func (t *TTYSteward) SetKeychainPassword(name, password string) error {
	return setKeychainPassword(name, password)
}

func newKeychainItemFor(name string) keychain.Item {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetAccount(name)
	item.SetService("vaulted")
	return item
}

func getKeychainPassword(name string) (string, error) {
	kc := keychain.NewWithPath(keychainPath)

	query := newKeychainItemFor(name)
	query.SetMatchSearchList(kc)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil || len(results) == 0 {
		return "", vaulted.ErrKeychainPasswordNotFound
	}

	return string(results[0].Data), nil
}

func setKeychainPassword(name, password string) error {
	kc := keychain.NewWithPath(keychainPath)
	err := kc.Status()
	if err == keychain.ErrorNoSuchKeychain {
		// create a new keychain if it doesn't exist
		kc, err = keychain.NewKeychainWithPrompt(keychainPath)
	}
	if err != nil {
		return err
	}

	item := newKeychainItemFor(name)
	item.UseKeychain(kc)
	item.SetLabel(name + " vault")
	item.SetDescription("vaulted password")
	item.SetData([]byte(password))

	query := newKeychainItemFor(name)
	query.SetMatchSearchList(kc)

	err = keychain.UpdateItem(query, item)
	if err == keychain.ErrorItemNotFound {
		err = keychain.AddItem(item)
	}
	return err
}
