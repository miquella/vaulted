// Code generated by go-bindata.
// sources:
// doc/man/vaulted-add.1
// doc/man/vaulted-cp.1
// doc/man/vaulted-dump.1
// doc/man/vaulted-edit.1
// doc/man/vaulted-env.1
// doc/man/vaulted-load.1
// doc/man/vaulted-ls.1
// doc/man/vaulted-rm.1
// doc/man/vaulted-shell.1
// doc/man/vaulted-upgrade.1
// doc/man/vaulted.1
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _vaultedAdd1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x8e\x31\x4e\x03\x31\x10\x45\xfb\x3d\xc5\x3f\x00\xb1\xc4\x11\x20\x44\x8a\x0b\x36\x56\x1c\x0a\x24\x37\xa3\xf5\x18\x2c\x25\xe3\xb0\x9e\x38\xe2\xf6\x68\x4d\x40\x54\xd4\xff\xbd\x37\x63\x0e\x5b\x34\xba\x1c\x95\x63\x58\x51\x8c\xb8\x1f\x8c\xdf\x62\x7c\x78\xde\x0c\xc6\xb9\xe1\xb6\x61\x99\xc2\x0a\x59\x94\x67\x9a\x34\x37\x3e\x7e\x62\x9a\x99\x94\x2b\xf4\x9d\x31\x15\x51\x16\x45\x49\x20\x08\x5f\xbf\xab\x3d\xe6\x5f\xc7\x9d\xf3\xd6\xf7\x60\x48\x8f\x21\xad\xff\x64\x43\xda\x23\x24\x2b\x74\xe2\x90\x5c\x17\x9e\x36\x7e\xbd\xb7\xee\x60\x77\x63\x77\xfc\x99\xae\x52\x41\xf2\x7b\xbf\x31\x4e\x25\x32\x52\x99\xc1\x31\x6b\x96\xb7\x7f\xbe\x30\xbd\xf2\x72\x2e\x82\x8f\x4b\xd6\x85\xbe\xeb\xf8\x42\xfc\x28\xb9\xa2\x52\xe3\x08\x2d\x7d\xbb\x99\x5f\x01\x00\x00\xff\xff\x34\xb1\x7b\x0e\x21\x01\x00\x00")

func vaultedAdd1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedAdd1,
		"vaulted-add.1",
	)
}

func vaultedAdd1() (*asset, error) {
	bytes, err := vaultedAdd1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-add.1", size: 289, mode: os.FileMode(420), modTime: time.Unix(1483131727, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedCp1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x90\xc1\x6a\xeb\x30\x10\x45\xf7\xfe\x8a\x59\xbd\x55\x62\x78\x9f\x90\x26\x86\x18\x5a\x47\x44\x69\x43\x41\x50\xc6\xf6\x08\x0b\x1c\xc9\x95\x14\x9b\xfc\x7d\x91\xa2\xc4\xb4\xd0\x2e\xda\x9d\xb1\x2e\xf7\x9c\xb9\xf9\x61\x0b\x23\x9e\x7b\x4f\xad\x58\x36\x03\xfc\xcf\x72\xbe\x85\x6a\xf5\x54\x64\x39\x63\x59\x7a\x82\x66\x00\xb1\x84\xc6\x0c\x8a\x1c\xf8\x8e\xa0\x31\xda\x93\xf6\x60\x24\xe0\xb5\x00\x50\xb7\xe0\x70\x24\x07\xca\x03\x3a\x40\xd0\x34\xa5\xb7\x49\xf9\x2e\xfd\x18\xd0\xb9\xc9\xd8\x36\x82\xf8\x6b\xb5\x63\xbc\xe4\x11\x26\xe4\x83\x90\xeb\x19\x29\xe4\x1e\x84\x2c\x4d\xdf\x0a\xc9\xc2\x97\xa6\x49\x48\x96\xe5\xb5\xfd\x9a\x35\xc3\xe5\xdb\x34\xdf\xc2\xa6\xe0\xeb\x7d\xc9\x0e\xe5\xae\x8a\xa4\x75\xb2\x57\x3a\x1e\x73\x0f\x27\x5b\xe5\xa0\xb1\x84\xa1\xd9\x58\xb0\x34\xf4\xd8\x50\x0b\xf5\xe5\x7e\xb6\xb4\xe6\x34\xd3\xc4\xbf\x3c\xd6\x96\x32\xd5\x05\xb7\x97\xd5\xf3\xe3\xa1\xd8\xbc\xb1\x15\xe7\xc7\xdd\x7e\x13\xfc\x48\x8f\xca\x1a\x7d\x0a\x15\x23\x5a\x85\x75\x4f\x81\xe6\xc8\x2f\xc2\x6a\x93\xea\x7b\xa8\x09\xce\x8e\xda\x30\xa1\xef\x28\xbb\xed\x05\xd2\xd8\x19\xb9\x00\xe3\x3b\xb2\x93\x72\x14\x99\xf7\xd4\xad\xc2\xd2\xfb\x99\x5c\x38\x61\x54\x18\x23\xde\x5f\x7e\xd0\xac\x8a\xe3\x5f\x54\xb3\x4f\x12\x49\xf5\x3a\xea\x6f\x55\x3f\x02\x00\x00\xff\xff\xce\x4b\xec\x94\x9b\x02\x00\x00")

func vaultedCp1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedCp1,
		"vaulted-cp.1",
	)
}

func vaultedCp1() (*asset, error) {
	bytes, err := vaultedCp1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-cp.1", size: 667, mode: os.FileMode(420), modTime: time.Unix(1483131727, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedDump1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\xcd\x41\x8a\x83\x30\x14\xc6\xf1\x7d\x4e\xf1\x5d\xc0\xc0\x1c\x61\x46\x05\x33\x30\x1a\x8c\x9b\x81\x6c\x42\xcd\xa3\x42\x93\x88\xbe\xb4\xd7\x2f\xa6\x5d\x94\x76\xf9\xf8\x78\xbf\xbf\x9c\x3a\x5c\x5d\xbe\xb0\x9f\x6d\x35\xe7\xb0\xe2\x4b\x48\xd3\xa1\xff\xfe\x6b\x85\xd4\x5a\x3c\x47\x94\xcd\x56\xb8\x6d\x0b\xfb\x1d\x7c\xf6\x38\xa5\xc8\x3e\x32\x12\xc1\x3d\x10\x70\xc2\xce\x73\xca\x0c\xb7\xe3\xd7\x0c\x7d\xc1\xcc\x7f\x3f\x68\xa3\x4c\x01\x2d\xfd\x58\xaa\x5f\x59\x4b\x23\x2c\xa9\xe8\x82\xb7\xa4\xcb\x47\xd3\x9a\x7a\x54\x7a\x52\x87\xa0\xb5\x68\x72\x58\x3f\xa2\xc7\xf9\x9e\x5d\x62\xc9\x82\xd2\x16\x1c\x4b\x71\x0f\x00\x00\xff\xff\xbe\x1d\xa8\x5d\xe0\x00\x00\x00")

func vaultedDump1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedDump1,
		"vaulted-dump.1",
	)
}

func vaultedDump1() (*asset, error) {
	bytes, err := vaultedDump1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-dump.1", size: 224, mode: os.FileMode(420), modTime: time.Unix(1483131728, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedEdit1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x91\xc1\x8e\xd3\x30\x10\x86\xef\x7e\x8a\xff\x84\x38\xec\x5a\x82\x37\x58\x96\xa2\x46\x88\x6e\x14\x17\xa1\x95\x7c\x71\xe3\x71\x63\x29\xb1\x8b\x67\xd2\xd0\xb7\x47\xf1\xd2\xaa\xe2\xb2\x57\xff\xfe\x3e\xcf\xfc\xd6\xfb\x2d\xce\x6e\x1e\x85\xbc\x7d\x24\x1f\x05\x9f\x94\x36\x5b\xec\x9e\x7e\x6c\x94\x6e\x5b\xf5\x2f\x44\xcd\xec\x23\x62\x12\x2a\xae\x97\x78\xa6\xf1\x52\x4f\x19\x32\x10\xfa\x9c\x84\x92\x20\x07\xb8\x04\xfa\x13\x59\x62\x3a\xbe\xb9\xab\xd1\xbc\xee\x5e\x5a\xd3\x98\x6a\xb5\xe1\x8b\x0d\xcf\xf7\x6e\x1b\x3a\xd8\xd0\x24\x37\x91\x0d\x6d\x25\xbe\x6e\xcc\x73\xd7\xb4\xfb\xe6\x65\x57\x21\x73\x72\x4b\xe2\x55\x7f\x1d\xe2\x4c\x98\xb2\x27\x84\x5c\xaa\x64\x7d\xf1\xbd\x61\x74\x75\xfd\x3c\xe5\x84\xdf\x73\x94\x35\x78\xa8\x50\xa2\xe5\x06\x46\x06\xbb\x33\x79\x48\xae\xd9\x95\x34\x5b\x3c\xfd\x32\xf8\xbe\x79\x55\xba\x33\x4a\x37\x2d\xec\xc7\xc3\x8c\xcf\xaa\x96\x63\xe6\x03\x4b\x94\x59\x08\x4b\x94\x01\x42\xd3\x29\x17\x57\x2e\xe8\x0b\x79\x4a\x12\xdd\xc8\x4a\x1f\x8a\xda\xe7\xe3\x71\x24\xc6\x32\x90\x0c\x54\x70\xc9\x73\xa9\xea\xbb\x8b\x70\x85\xc0\x37\xa5\x7f\x73\x3a\x30\xd5\xcd\x6e\x72\x75\xc7\x68\x7c\xcb\x05\x53\x2e\x04\x4f\xe2\xe2\xc8\xc8\x09\x32\x44\xc6\xa9\xe4\x9e\x98\x1f\xc0\x44\x75\x29\x9f\xfb\x79\xa2\x24\x4e\x62\x4e\x6b\x87\xff\x7d\x0b\x0f\x34\x8e\x36\x74\xf6\x83\x56\xba\xdb\xa8\xbf\x01\x00\x00\xff\xff\xbe\x8e\xbd\xbc\x2c\x02\x00\x00")

func vaultedEdit1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedEdit1,
		"vaulted-edit.1",
	)
}

func vaultedEdit1() (*asset, error) {
	bytes, err := vaultedEdit1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-edit.1", size: 556, mode: os.FileMode(420), modTime: time.Unix(1483131728, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedEnv1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x92\x41\x8b\xdb\x30\x10\x85\xef\xfe\x15\xf3\x03\xb2\x86\x5e\x7b\xdb\xa6\x01\x1b\xb6\x59\xb3\x0a\x94\x82\x2f\x8a\x35\x8a\x87\xc8\xa3\x20\x8d\xe2\xe6\xdf\x17\xc9\x5e\x6a\x96\x5c\xf5\xfc\xde\xd3\xfb\xac\xfa\xd4\xc0\x5d\x27\x27\x68\xfa\x17\xe4\x3b\x7c\xab\x6a\xd5\xc0\xf1\xf5\xd7\xa1\xaa\xbb\xae\x5a\x35\xc8\x52\xff\x02\x3e\xc9\x2d\x49\x84\x38\xa2\x73\x30\xf8\x69\xd2\x6c\x22\xc8\xa8\x05\x9c\xd7\x06\x22\x0e\x01\x25\x82\xf5\x01\xf4\x92\x0c\xc4\xe2\x41\x46\x5c\x5c\x25\x5f\xfd\x39\xbe\x77\xaa\x55\xa5\xa3\xb7\x3f\x7a\xbb\xdf\x34\xf5\xf6\x03\x7a\xdb\xb2\x9e\xb0\xb7\x5d\x31\xfc\x3c\xa8\xfd\x47\xdb\x9d\xda\xf7\x63\xf1\xec\x03\x6a\xc1\x08\x9a\xb3\x81\x82\xe7\x09\x59\x20\x45\xe2\x0b\xdc\x75\x20\x7d\x76\x45\x36\xa5\xf9\xf5\xb7\x82\x2b\x3e\x20\x8a\x0f\x68\x80\xb8\x9c\x96\xca\x1a\x4e\x23\x56\x01\x63\x72\x92\xcd\xdb\xb8\x4d\x50\x40\x48\x11\x0d\x88\x87\x0b\x32\x06\x2d\xf8\x94\xc2\x4c\xce\x55\x05\x45\x59\xbc\xe2\x28\x08\xf4\x62\xa8\xcb\x80\xd3\x27\x0f\xa0\x08\x3a\x89\x37\x28\x38\x64\x00\x36\xf8\xa9\x98\x17\x2e\xaa\x39\xbc\xbd\x65\x22\xcf\x2e\xb6\x03\xb2\x1b\xb4\x14\x21\xf1\x95\xfd\xcc\xe0\x03\x24\x8e\x37\x1c\xc8\x12\x9a\xdd\x1a\x16\xc7\x9c\x34\xf8\xe9\xa6\x85\xce\x0e\xff\x5f\x3e\x0f\xc4\x89\x44\xd0\xd4\xeb\x5f\x69\x8f\x5e\xf0\x7b\x6f\x3b\x50\xaa\xc9\xf8\x96\xaf\xe8\xc2\x05\xe2\x3c\x22\x7f\xb2\xf8\x02\x2e\xb3\xa0\x08\xb3\x7e\x64\xba\x14\xf3\x46\x93\xb0\x5a\xdf\x01\xb1\x3e\x93\x23\x79\x64\x9a\x12\xf4\x70\x2d\xc7\x8e\x2c\x0a\x4d\x08\x7e\xd9\xb4\x09\xdc\xc1\x3c\xd2\x30\xc2\x84\x9a\x63\x11\x95\x6a\x2a\x7d\xc9\x2c\x66\x9f\x9c\x01\xfc\x4b\x31\x3f\x35\x83\x96\x98\x04\xdd\xa3\xae\xfe\x05\x00\x00\xff\xff\x1b\x2e\x0b\x7d\xdd\x02\x00\x00")

func vaultedEnv1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedEnv1,
		"vaulted-env.1",
	)
}

func vaultedEnv1() (*asset, error) {
	bytes, err := vaultedEnv1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-env.1", size: 733, mode: os.FileMode(420), modTime: time.Unix(1483131728, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedLoad1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x5c\xce\x5f\x6a\x03\x21\x18\x04\xf0\x77\x4f\x31\x17\x88\xd0\x23\xb4\x69\x20\x16\xea\xca\x9a\x97\x82\x2f\xb2\x7e\x12\x61\xab\x41\xbf\xdd\x5e\xbf\x54\xfb\x8f\xbc\x0d\x0c\xc3\x6f\xe4\xe5\x8c\xdd\x6f\x2b\x53\x70\x87\xb5\xf8\x80\x07\x21\xed\x19\xfa\xf1\xf5\x24\xa4\x31\xe2\xbb\x44\xef\xdc\x01\x5b\xa3\x86\x17\x3b\x69\xdc\x6a\xd9\x53\xa0\x00\x2e\x68\x1c\x52\xfe\x0a\x4b\x25\xcf\x84\x52\x51\xe9\xb6\xfa\x85\xc0\x57\xc2\x52\x32\x53\x66\x94\x08\x3f\xb8\x8e\xd8\x37\x3d\x19\xab\x6c\x87\x5c\x7c\x72\xf1\xf8\x9f\x73\x71\x86\x8b\x2a\xfb\x77\x72\xd1\xf4\xc5\xf3\xc9\x1e\x67\x65\x2e\x6a\xd2\x7d\x34\x0f\xa4\xdd\x2b\x7f\x33\x7c\x24\xbe\x8e\xc3\x3f\xfd\xef\xf1\x3d\xf9\xf1\x5c\x8a\xcf\x00\x00\x00\xff\xff\x29\xac\xab\x44\x08\x01\x00\x00")

func vaultedLoad1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedLoad1,
		"vaulted-load.1",
	)
}

func vaultedLoad1() (*asset, error) {
	bytes, err := vaultedLoad1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-load.1", size: 264, mode: os.FileMode(420), modTime: time.Unix(1483131728, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedLs1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x5c\x8d\x41\x0a\x83\x30\x10\x45\xf7\x39\xc5\x3f\x80\x06\x7a\x84\xd6\x0a\x0a\xad\x06\xc7\x4d\x21\x1b\x8b\x11\x02\x83\x29\xce\xd8\xf3\x17\xd2\x42\xc1\xed\xbc\xf7\xe6\xdb\xb1\xc1\x7b\xda\x59\xc3\xec\x4b\x16\x9c\x8c\xa5\x06\xdd\xf9\x5e\x1b\xeb\x9c\xf9\x21\xb0\xc0\x97\xe0\x28\x2a\x98\x98\xbf\x89\x64\x97\x1e\x5d\xef\xa8\xa5\xec\xfb\xe5\xe2\x97\xea\x5f\xf9\x65\x30\xf6\xb9\x1d\xef\x51\x34\x13\x6a\x70\xad\xa9\x1a\x5a\x37\xb6\x7d\x97\x3f\xdc\x0e\x1b\x05\xd2\x1a\xf0\x0a\x1b\x38\xae\xa1\x80\x26\x88\xce\x69\x57\x6b\x3e\x01\x00\x00\xff\xff\x32\x37\x94\xc4\xbc\x00\x00\x00")

func vaultedLs1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedLs1,
		"vaulted-ls.1",
	)
}

func vaultedLs1() (*asset, error) {
	bytes, err := vaultedLs1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-ls.1", size: 188, mode: os.FileMode(420), modTime: time.Unix(1483131728, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedRm1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x54\x8f\x4d\x6a\xc3\x30\x10\x85\xf7\x3e\xc5\x5b\x65\xd5\x08\x7a\x84\x36\x0d\xc4\x8b\x26\xc2\xf6\xa6\x30\x1b\xd9\x1a\xd5\x02\x5b\x4a\xf5\x13\x9a\xdb\x17\x2b\x82\xb6\xbb\x81\x79\xdf\xfb\x66\xc4\x70\xc2\x4d\xe5\x25\xb1\xa6\x7d\x58\xf1\xdc\x88\xfe\x84\xf3\xcb\xfb\xb1\x11\x52\x36\x75\x85\xb0\x82\xf6\x08\xbc\xfa\x1b\x47\xf0\xb7\x8d\xc9\xba\xcf\x07\x19\x0b\xd2\x7f\x9c\x2f\xb2\x6f\xfb\x82\x91\x79\x25\x73\xf8\x85\xc9\x74\x20\xd3\x3a\xb5\x32\x19\xb9\x8d\xb4\x13\x42\x90\x91\x85\x7d\x3b\xf6\x87\xae\x95\x43\x7b\x39\x17\xbc\xab\x9e\x34\x73\x55\x20\x5e\x79\xb2\xc6\xb2\xc6\x78\xff\x53\x45\x3b\x81\x61\xe6\xed\xa2\x84\xc9\x6b\x86\x8d\xe0\xaf\xac\x16\x24\x5f\x78\x97\xd7\x91\x03\xbc\x69\x6a\x53\x9a\xd5\x16\xcd\x8b\x86\xf3\x09\x23\xd7\xb7\xb4\x28\xee\xd6\x40\x3d\xa4\x98\x94\xfb\x9f\x78\x2a\x8d\x1c\x82\x0f\x9b\x47\xdb\x78\x5d\xd4\x9d\x35\xbc\x43\x4c\xda\xe7\x24\x9a\x9f\x00\x00\x00\xff\xff\x85\x0f\x9d\xfc\x51\x01\x00\x00")

func vaultedRm1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedRm1,
		"vaulted-rm.1",
	)
}

func vaultedRm1() (*asset, error) {
	bytes, err := vaultedRm1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-rm.1", size: 337, mode: os.FileMode(420), modTime: time.Unix(1483131729, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedShell1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x93\x41\x6b\xe3\x30\x10\x85\xef\xfe\x15\x43\x4f\x2d\xa4\x2e\xbb\xc7\xee\x29\x4d\x0d\x35\xe9\x26\xc1\xf2\x6e\x09\x18\x82\x62\x8d\x62\x51\x5b\x32\xd2\x38\xc1\xff\x7e\x91\xec\x6c\xe2\x6d\xe9\x6d\x8f\x89\x66\xde\x7c\xef\xcd\x38\xce\x5f\xe0\xc8\xbb\x9a\x50\x14\xf7\xae\xc2\xba\x86\x6f\x51\xcc\x5e\x60\x35\xff\x99\x44\xf1\x66\x13\x8d\xaf\x30\x3c\x16\xf7\xe0\x88\x5b\x72\xc0\x35\x28\x4d\x68\x79\x49\xea\x88\xe3\xf3\x49\x51\x05\x54\x21\x38\x2c\x2d\x92\x03\x69\x6c\xf8\x1d\x54\xa0\x36\x5c\xa0\xf0\x7d\x66\xa8\xf2\x4d\x61\x1c\xdb\xae\xd6\x1b\x96\xb2\x30\xb2\x90\x4f\x85\x5c\x4c\x06\x17\x32\x83\x42\xa6\x9a\x37\x58\xc8\x4d\x68\x79\x4e\xd8\x22\x4b\x37\x79\xba\x5e\x85\x2e\xf6\x05\xd7\x6d\xe7\xd0\x85\x91\x83\x36\x7b\x49\x5e\x5f\xbd\x26\xea\xa3\xb2\x46\x37\xa8\x09\x8e\xdc\x2a\xbe\xaf\x71\x06\x4a\x82\x43\xfa\x11\x19\xaa\xd0\x9e\x94\x43\x10\x28\x3d\x8e\x03\x32\xa3\xc4\xc3\x5e\xe9\x07\x57\x15\x32\xbb\x8b\x03\xcf\xfc\x8d\xc1\x32\xd9\x06\x96\x95\x21\x7c\x84\xa7\xfe\xdc\x37\x83\xdf\x67\x37\xdd\xde\x91\xa2\x8e\xd0\x01\x07\xc2\xa6\x35\x96\xdb\xde\xcf\x03\x23\xa1\xb4\x28\x50\x93\xe2\xb5\x83\x53\x85\x1a\x5c\xcb\x4f\x5a\xe9\x83\xf7\x75\x05\x1b\x47\x79\x85\x61\xe4\x3b\xf6\xa0\x74\xdb\x11\x54\x68\x11\x1a\xde\x83\x36\x04\x0d\xa7\x72\x58\x85\x2f\xb8\x0e\xbe\x37\x9d\x9d\x48\x41\x5e\x29\x07\x12\x39\x75\x16\xa1\xe4\x1a\xf6\x08\x64\x0e\x87\x3a\xb4\x78\x91\x7f\x56\x82\x42\x91\x4f\xaf\x41\xdd\xc5\xc1\xf0\xd9\x5e\xc8\x39\x4c\x70\x64\x2c\x8a\x89\x21\x32\x70\x40\x8d\x96\x13\x02\x3f\x3b\xbe\x44\xc0\x72\x06\x64\xde\x51\xbb\x88\x2a\x4e\x60\xb1\xe1\x4a\xc3\xb8\xc8\xf3\x25\x89\xce\x72\x52\xc6\x27\x83\xa5\x92\xea\x2f\xe3\x70\x63\x03\x4e\x70\xa4\x9a\x96\x97\x34\xac\x5d\x9a\xba\x36\x27\x9f\xe3\x67\x1b\x77\x8f\x51\x9c\xb1\x28\x4e\x37\x50\xdc\xee\x3b\xf8\x3e\xfa\x9d\xbf\xb1\xdd\x7c\xb1\x48\x18\xdb\x2d\x93\xed\x2e\x7d\x2e\x64\x16\xc5\x7b\x3b\xea\x0f\xd2\x17\x03\xbc\x2c\xd1\xb9\x61\x23\x02\x46\xde\xe9\x9f\x63\x2a\xd7\xc4\xb3\x48\x60\x8b\x5a\x78\x38\xa3\xfd\xd6\xfd\xd5\x5d\x0e\xc5\x7b\x55\x0e\x50\x7b\x52\x11\x7f\x4e\xc9\x92\x45\x96\xe4\x57\xb0\x5f\x93\x0e\x9f\xe7\x35\xdb\x48\xfb\xf1\xe1\xff\x11\x33\x96\xae\x57\xbb\x7c\xbd\x4c\x56\xfe\x98\x6e\xb9\x10\xca\xf7\xf2\xba\xee\x67\x30\xf1\xf6\x2b\x4b\xf3\xed\xa5\x54\xb9\x70\x3d\x64\xc0\x75\x6d\x6b\x2c\x41\x8d\x07\x5e\xf6\xc0\x9e\x97\xee\x6e\xe2\x7b\xac\xbb\xb9\xf1\x1f\xf5\x25\x80\xcb\x59\x7e\xc0\x16\xca\x9d\xb9\xb3\x24\xfa\x13\x00\x00\xff\xff\x49\x80\xd2\x39\x21\x05\x00\x00")

func vaultedShell1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedShell1,
		"vaulted-shell.1",
	)
}

func vaultedShell1() (*asset, error) {
	bytes, err := vaultedShell1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-shell.1", size: 1313, mode: os.FileMode(420), modTime: time.Unix(1483131729, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaultedUpgrade1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x90\x4d\x6e\xc3\x20\x10\x85\xf7\x9c\x62\x2e\x10\xa4\x1e\xa1\x4d\x23\xc5\x8b\x3a\x96\xf1\xa6\x12\x9b\x89\x67\x88\x23\xd9\x90\xf2\x93\xb6\xb7\xaf\xc0\xa1\x0b\x2f\xb2\x43\xbc\xf7\xbe\x4f\x20\x87\x23\xdc\x31\xcd\x91\x49\xef\xd2\xed\xe2\x91\x18\x5e\x84\x54\x47\x68\x5f\x3f\x0e\x42\x76\x9d\x78\xe4\x50\x63\xbd\xab\xc7\x00\x33\x5f\x70\xfc\x5d\x11\x01\xa2\x83\x38\x31\x8c\xc9\x7b\xb6\x71\xbd\x05\xe3\xfc\x82\xb1\x20\xd5\x67\x7b\xea\x54\xa3\x0a\x56\x9b\x37\x6d\xf6\x1b\xb8\x36\x7d\x69\xbe\x1f\xd4\xbe\x6f\xba\xa1\x39\xb5\xa5\xdc\x33\xd2\xd6\x86\x96\x60\x74\xf6\xce\x3e\xab\x27\x5e\x9e\xf9\x25\x0c\x13\x43\xc0\x85\xc5\x0d\x43\xf8\x76\x9e\xe0\x1a\x20\x05\xa6\xdc\x58\x77\x2b\x8c\xe9\x61\x90\x45\x9d\x77\xfc\x73\x8d\x30\x3a\xe2\xbc\xe1\xaf\x84\x73\x75\xd9\xb4\x9c\xd9\x83\x33\xff\x7f\x30\x61\xae\xa6\x99\xc0\xba\x08\x67\xae\x4f\x23\x29\xfe\x02\x00\x00\xff\xff\x93\xa5\x62\x52\x6e\x01\x00\x00")

func vaultedUpgrade1Bytes() ([]byte, error) {
	return bindataRead(
		_vaultedUpgrade1,
		"vaulted-upgrade.1",
	)
}

func vaultedUpgrade1() (*asset, error) {
	bytes, err := vaultedUpgrade1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted-upgrade.1", size: 366, mode: os.FileMode(420), modTime: time.Unix(1483131729, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _vaulted1 = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x94\xd1\x6e\xdb\x3c\x0c\x85\xef\xfd\x14\xbc\xfa\xff\x16\x58\x34\xf4\x11\xda\xb4\x40\x3d\x2c\x89\x11\x77\x18\x86\xb9\x17\xaa\x45\x25\x02\x6c\xc9\x10\x69\x67\x79\xfb\x41\x92\x93\x66\x8e\x8b\x01\xbb\xb3\xe4\xf0\x3b\x87\x87\x74\xc4\xcb\x33\x0c\xb2\x6f\x18\x15\xdc\x65\xa2\x7c\x86\xf5\xfd\xea\x29\x13\x45\x91\x9d\xae\xab\x05\x50\x27\x0f\x16\xd0\x0e\xc6\x3b\xdb\xa2\x65\x02\xed\x5d\x0b\x84\x75\xef\xb1\x39\x02\xb1\xf3\xa8\xc2\xd9\x23\x53\xe4\x94\x3f\xd6\x9b\xa2\xcc\xcb\xc8\xaa\xf4\x43\xa5\x97\x23\xb1\xd2\x5b\x48\x17\xd5\xc2\xa6\x43\x6e\x65\x8b\x95\x2e\xe0\xe7\xe9\x85\xa9\xf4\xf6\x35\x13\x6f\xfe\x1f\x6a\xab\x45\x28\x0e\xaf\x96\xab\xc7\x4a\x17\x1f\x59\xc8\x97\x9b\xd5\xea\x7e\xfd\x38\x16\xe7\xd2\xef\x48\x08\x51\xe9\xe2\x35\xb6\xf0\xf8\x54\x2e\xb7\x79\xf1\x92\x6f\xd6\x11\x91\x6b\xb0\x6e\x52\x67\x08\x3a\xef\x06\xa3\x50\x7d\x82\x2b\x0d\x34\xbc\x47\x9f\xf2\xa3\x77\x43\x70\x63\xf4\xb9\xec\x16\x9c\xcf\xc6\x5f\x48\x0b\xc6\x32\x7a\x59\xb3\x19\x10\x68\x8f\x4d\x23\x2e\xec\x8f\xbd\x41\x2b\x8f\xf0\x86\xd0\x13\x2a\x60\x07\xca\x68\x8d\x1e\x2d\x1b\xc9\x08\xbc\xc7\x0b\xa9\x38\xa8\xa9\xb1\xea\xbf\xff\x09\xdc\xc1\x82\xf4\xbb\x3e\x0e\x54\xc4\x8e\xc7\xc6\xca\x4c\xbc\x9c\x24\xa5\x0a\x05\x59\xfe\x6e\xab\x39\x42\xed\x51\x32\x52\x94\xaa\x9d\x65\xb4\x0c\x4e\x83\x04\x8b\x87\xb4\x4f\x02\x4a\x44\xc8\xc4\xc3\xf6\xb4\x5f\x0b\xa9\x14\xdc\xdc\xdd\x8a\x0b\x78\xdd\x85\x6e\x3e\x8f\xfe\x6a\xd7\x1d\x83\xd6\xd2\x75\x66\x0e\x1e\x41\x20\xad\x02\x92\x03\x12\x18\x06\x49\x97\xa2\x70\x30\xbc\x1f\x2f\x3a\x49\x74\x70\x5e\xcd\x18\xa9\xbb\xa9\x0f\xd5\xb7\xc1\x49\xf6\xdd\x9b\xd9\xb6\x12\x9d\x1d\x10\x2b\xd7\x47\xd9\x2f\xe5\x66\x3d\xc3\x0e\xa4\x29\x1d\x95\xe1\xeb\x0c\xc3\xed\xb5\x94\x05\xfc\x65\x88\x8d\xdd\x7d\x98\x63\x28\xbc\x92\xb0\x43\x50\xd8\xf4\xdc\xf5\x4c\x69\x71\xa0\x76\x6d\x2b\xad\x0a\x22\x92\xa1\x71\xf2\xfc\x85\x82\x76\xfe\xdc\x96\xb1\xec\xa2\x8f\xb4\x6e\x33\x82\x76\x98\xea\x05\x58\x10\xfc\x46\x98\xa2\x38\xaf\xf3\x98\x92\xb1\xe1\x21\xed\x09\x38\x0f\x1e\xbb\x46\xd6\xf8\x41\xb4\x33\xa2\xd1\xee\x54\x95\x2e\xd7\xa5\x31\x14\x63\xfd\x6a\x88\x09\x64\xd3\xa4\x5a\x9a\x83\xd1\x14\xe5\xdb\x50\xba\xc5\xd6\x85\x4d\xfa\x33\xf3\x39\x82\x6f\xa7\x84\x98\x56\x80\x94\x2c\x3d\xcf\x7f\xbb\x69\x21\x63\xb6\x17\xc1\x87\x73\x8a\x3e\x34\x89\xea\xef\x13\x48\xb0\x89\x81\xbe\xdb\x79\xa9\x30\x8e\x21\x3d\x12\x34\xb8\x93\xf5\x71\x6c\x03\x46\x6a\xdd\xfb\xf0\xe7\x30\x6a\x6a\xe7\x5b\x39\x97\xf8\xc8\x4b\x32\xbf\x03\x00\x00\xff\xff\x2d\x80\xe6\x9a\x19\x06\x00\x00")

func vaulted1Bytes() ([]byte, error) {
	return bindataRead(
		_vaulted1,
		"vaulted.1",
	)
}

func vaulted1() (*asset, error) {
	bytes, err := vaulted1Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "vaulted.1", size: 1561, mode: os.FileMode(420), modTime: time.Unix(1483131729, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"vaulted-add.1": vaultedAdd1,
	"vaulted-cp.1": vaultedCp1,
	"vaulted-dump.1": vaultedDump1,
	"vaulted-edit.1": vaultedEdit1,
	"vaulted-env.1": vaultedEnv1,
	"vaulted-load.1": vaultedLoad1,
	"vaulted-ls.1": vaultedLs1,
	"vaulted-rm.1": vaultedRm1,
	"vaulted-shell.1": vaultedShell1,
	"vaulted-upgrade.1": vaultedUpgrade1,
	"vaulted.1": vaulted1,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"vaulted-add.1": &bintree{vaultedAdd1, map[string]*bintree{}},
	"vaulted-cp.1": &bintree{vaultedCp1, map[string]*bintree{}},
	"vaulted-dump.1": &bintree{vaultedDump1, map[string]*bintree{}},
	"vaulted-edit.1": &bintree{vaultedEdit1, map[string]*bintree{}},
	"vaulted-env.1": &bintree{vaultedEnv1, map[string]*bintree{}},
	"vaulted-load.1": &bintree{vaultedLoad1, map[string]*bintree{}},
	"vaulted-ls.1": &bintree{vaultedLs1, map[string]*bintree{}},
	"vaulted-rm.1": &bintree{vaultedRm1, map[string]*bintree{}},
	"vaulted-shell.1": &bintree{vaultedShell1, map[string]*bintree{}},
	"vaulted-upgrade.1": &bintree{vaultedUpgrade1, map[string]*bintree{}},
	"vaulted.1": &bintree{vaulted1, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

