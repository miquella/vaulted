package vaulted

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type ProxyKeyring struct {
	keyring  agent.Agent
	upstream agent.Agent

	listener net.Listener
}

func NewProxyKeyring(upstreamAuthSock string) (*ProxyKeyring, error) {
	var err error
	var conn net.Conn
	var upstream agent.Agent

	if upstreamAuthSock != "" {
		conn, err = net.Dial("unix", upstreamAuthSock)
		if err != nil {
			return nil, err
		}

		upstream = agent.NewClient(conn)
	}

	return &ProxyKeyring{
		upstream: upstream,
		keyring:  agent.NewKeyring(),
	}, nil
}

func (pk *ProxyKeyring) Listen() (string, error) {
	if pk.listener != nil {
		return "", errors.New("Already listening")
	}

	dir, err := ioutil.TempDir("", "proxykeyring")
	if err != nil {
		return "", err
	}

	err = os.Chmod(dir, 0700)
	if err != nil {
		return "", err
	}

	listener := filepath.Join(dir, "listener")
	pk.listener, err = net.Listen("unix", listener)
	if err != nil {
		return "", err
	}

	err = os.Chmod(listener, 0600)
	if err != nil {
		return "", err
	}

	return listener, nil
}

func (pk *ProxyKeyring) Serve() error {
	if pk.listener == nil {
		return errors.New("Not listening")
	}

	for {
		c, err := pk.listener.Accept()
		if err != nil {
			return err
		}

		go agent.ServeAgent(pk, c)
	}
}

func (pk *ProxyKeyring) Close() error {
	if pk.listener != nil {
		return pk.listener.Close()
	}

	return nil
}

func (pk *ProxyKeyring) List() ([]*agent.Key, error) {
	keys, err := pk.keyring.List()
	if err != nil {
		return nil, err
	}

	if pk.upstream != nil {
		ukeys, err := pk.upstream.List()
		if err != nil {
			log.Printf("[ProxyKeyring] Upstream list error: %v", err)
		} else {
			keys = append(keys, ukeys...)
		}
	}

	return keys, nil
}

func (pk *ProxyKeyring) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	sig, err := pk.keyring.Sign(key, data)
	if err == nil {
		return sig, nil
	}

	if pk.upstream != nil {
		usig, uerr := pk.upstream.Sign(key, data)
		if uerr == nil {
			return usig, nil
		}
	}

	return nil, err
}

func (pk *ProxyKeyring) Add(key agent.AddedKey) error {
	return pk.keyring.Add(key)
}

func (pk *ProxyKeyring) Remove(key ssh.PublicKey) error {
	err := pk.keyring.Remove(key)

	if pk.upstream != nil {
		uerr := pk.upstream.Remove(key)
		if uerr == nil {
			err = nil
		}
	}

	return err
}

func (pk *ProxyKeyring) RemoveAll() error {
	err := pk.keyring.RemoveAll()

	if pk.upstream != nil {
		uerr := pk.upstream.RemoveAll()
		if err == nil {
			err = uerr
		}
	}

	return err
}

func (pk *ProxyKeyring) Lock(passphrase []byte) error {
	err := pk.keyring.Lock(passphrase)

	if pk.upstream != nil {
		uerr := pk.upstream.Lock(passphrase)
		if err == nil {
			err = uerr
		}
	}

	return err
}

func (pk *ProxyKeyring) Unlock(passphrase []byte) error {
	err := pk.keyring.Unlock(passphrase)

	if pk.upstream != nil {
		uerr := pk.upstream.Unlock(passphrase)
		if err == nil {
			err = uerr
		}
	}

	return err
}

func (pk *ProxyKeyring) Signers() ([]ssh.Signer, error) {
	signers, err := pk.keyring.Signers()
	if err != nil {
		return nil, err
	}

	if pk.upstream != nil {
		usigners, err := pk.upstream.Signers()
		if err != nil {
			log.Printf("[ProxyKeyring] Upstream signers error: %v", err)
		} else {
			signers = append(signers, usigners...)
		}
	}

	return signers, nil
}
