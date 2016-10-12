package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
	"golang.org/x/crypto/ssh"
)

var (
	green = color.New(color.FgGreen)
	cyan  = color.New(color.FgCyan)
	blue  = color.New(color.FgBlue)
)

type Edit struct {
	VaultName string
	rlMenu    *readline.Instance
	rlValue   *readline.Instance
}

func (e *Edit) Run(steward Steward) error {
	var password string
	var vault *vaulted.Vault
	var err error

	if vaulted.VaultExists(e.VaultName) {
		password, vault, err = steward.OpenVault(e.VaultName, nil)
		if err != nil {
			return err
		}
	} else {
		vault = &vaulted.Vault{}
	}

	e.edit(e.VaultName, vault)

	var newPassword *string
	if password != "" {
		newPassword = &password
	}
	err = steward.SealVault(e.VaultName, newPassword, vault)
	if err != nil {
		return err
	}

	return nil
}

func mainMenu() {
	color.Set(color.FgYellow)
	print("")
	print("a - AWS Key")
	print("s - SSH Keys")
	print("v - Variables")
	print("d - Environment Duration")
	print("? - Help")
	print("q - Quit")
	color.Unset()
}

func awsMenu() {
	color.Set(color.FgYellow)
	print("")
	print("k - Key")
	print("m - MFA")
	print("r - Role")
	print("s - Show Key")
	print("D - Delete")
	print("? - Help")
	print("b - Back")
	color.Unset()
}

func sshKeysHelp() {
	color.Set(color.FgYellow)
	print("")
	print("a - Add")
	print("D - Delete")
	print("? - Help")
	print("b - Back")
	color.Unset()
}

func variableMenu() {
	color.Set(color.FgYellow)
	print("")
	print("a - Add")
	print("D - Delete")
	print("? - Help")
	print("b - Back")
	color.Unset()
}

func (e *Edit) edit(name string, v *vaulted.Vault) {
	exit := false
	for exit == false {
		cyan.Printf("\nVault: ")
		fmt.Printf("%s", name)
		printVariables(v)
		printAWS(v, false)
		printSSHKeys(v)
		printDuration(v)

		input := e.readMenu("Edit vault: [a,s,v,d,?,q]: ")
		switch input {
		case "a":
			e.aws(v)
		case "s":
			e.sshKeysMenu(v)
		case "v":
			e.variables(v)
		case "d":
			dur := e.readValue("Duration (e.g. 15m or 36h): ")
			duration, err := time.ParseDuration(dur)
			if err != nil {
				color.Red("%s", err)
				break
			}
			if duration < 15*time.Minute || duration > 36*time.Hour {
				color.Red("Duration must be between 15m and 36h")
				break
			}
			v.Duration = duration
		case "q":
			exit = true
		case "?", "help":
			mainMenu()
		default:
			color.Red("Command not recognized")
		}
	}
}

func (e *Edit) aws(v *vaulted.Vault) {
	exit := false
	show := false

	for exit == false {
		var input string
		printAWS(v, show)
		if v.AWSKey == nil {
			input = e.readMenu("Edit AWS key [k,?,b]: ")
		} else {
			input = e.readMenu("Edit AWS key [k,m,r,s,D,?,b]: ")
		}

		switch input {
		case "k":
			awsAccesskey := e.readValue("Key ID: ")
			awsSecretkey := e.readValue("Secret: ")
			v.AWSKey = &vaulted.AWSKey{
				ID:     awsAccesskey,
				Secret: awsSecretkey,
				MFA:    "",
				Role:   "",
			}
		case "m":
			if v.AWSKey != nil {
				awsMfa := e.readValue("MFA ARN or serial number: ")
				v.AWSKey.MFA = awsMfa
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "r":
			if v.AWSKey != nil {
				awsRole := e.readValue("Role ARN: ")
				v.AWSKey.Role = awsRole
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "s":
			if v.AWSKey != nil {
				show = !show
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "D":
			if v.AWSKey != nil {
				removeKey := e.readValue("Delete your AWS key? (y/n): ")
				if removeKey == "y" {
					v.AWSKey = nil
				}
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "b":
			exit = true
		case "?", "help":
			awsMenu()
		default:
			color.Red("Command not recognized")
		}
	}
}

func (e *Edit) sshKeysMenu(v *vaulted.Vault) {
	exit := false

	for exit == false {
		printSSHKeys(v)
		input := e.readMenu("Edit ssh keys: [a,D,?,b]: ")
		switch input {
		case "a":
			e.addSSHKey(v)
		case "D":
			key := e.readValue("Key: ")
			_, ok := v.SSHKeys[key]
			if ok {
				delete(v.SSHKeys, key)
			} else {
				color.Red("Key '%s' not found", key)
			}
		case "b":
			exit = true
		case "?", "help":
			sshKeysHelp()
		default:
			color.Red("Command not recognized")
		}
	}
}

func (e *Edit) addSSHKey(v *vaulted.Vault) {
	homeDir := ""
	user, err := user.Current()
	if err == nil {
		homeDir = user.HomeDir
	} else {
		homeDir = os.Getenv("HOME")
	}

	defaultFilename := ""
	filename := ""
	if homeDir != "" {
		defaultFilename = filepath.Join(homeDir, ".ssh", "id_rsa")
		filename = e.readValue(fmt.Sprintf("Key file (default: %s): ", defaultFilename))
		if filename == "" {
			filename = defaultFilename
		}
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(filepath.Join(homeDir, ".ssh"), filename)
		}
	} else {
		filename = e.readValue("Key file: ")
	}

	decryptedBlock, err := loadAndDecryptKey(filename)
	if err != nil {
		color.Red("%v", err)
		return
	}

	comment := loadPublicKeyComment(filename + ".pub")
	var name string
	if comment != "" {
		name = e.readValue(fmt.Sprintf("Name (default: %s): ", comment))
		if name == "" {
			name = comment
		}
	} else {
		name = e.readValue("Name: ")
		if name == "" {
			name = filename
		}
	}

	if v.SSHKeys == nil {
		v.SSHKeys = make(map[string]string)
	}
	v.SSHKeys[name] = string(pem.EncodeToMemory(decryptedBlock))
}

func loadAndDecryptKey(filename string) (*pem.Block, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	if x509.IsEncryptedPEMBlock(block) {
		var passphrase string
		var decryptedBytes []byte
		for i := 0; i < 3; i++ {
			passphrase, err = ask.HiddenAsk("Passphrase: ")
			if err != nil {
				return nil, err
			}

			decryptedBytes, err = x509.DecryptPEMBlock(block, []byte(passphrase))
			if err == nil {
				break
			}
			if err != x509.IncorrectPasswordError {
				return nil, err
			}
		}

		if err != nil {
			return nil, err
		}

		return &pem.Block{
			Type:  block.Type,
			Bytes: decryptedBytes,
		}, nil
	}
	return block, nil
}

func loadPublicKeyComment(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}

	_, comment, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return ""
	}
	return comment
}

func (e *Edit) variables(v *vaulted.Vault) {
	exit := false

	for exit == false {
		printVariables(v)
		input := e.readMenu("Edit environment variables: [a,D,?,b]: ")
		switch input {
		case "a":
			variableKey := e.readValue("Name: ")
			variableValue := e.readValue("Value: ")
			if v.Vars == nil {
				v.Vars = make(map[string]string)
			}
			v.Vars[variableKey] = variableValue
		case "D":
			variable := e.readValue("Variable name: ")
			_, ok := v.Vars[variable]
			if ok {
				delete(v.Vars, variable)
			} else {
				color.Red("Variable '%s' not found", variable)
			}
		case "b":
			exit = true
		case "?", "help":
			variableMenu()
		default:
			color.Red("Command not recognized")
		}
	}
}

func print(message string) {
	fmt.Printf("%s\n", message)
}

func printVariables(v *vaulted.Vault) {
	color.Cyan("\nVariables:")
	if len(v.Vars) > 0 {
		var keys []string
		for key := range v.Vars {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			green.Printf("  %s: ", key)
			fmt.Printf("%s\n", v.Vars[key])
		}
	} else {
		print("  [Empty]")
	}
}

func printAWS(v *vaulted.Vault, show bool) {
	color.Cyan("\nAWS Key:")
	if v.AWSKey != nil {
		green.Printf("  Key ID: ")
		fmt.Printf("%s\n", v.AWSKey.ID)
		green.Printf("  Secret: ")
		if !show {
			fmt.Printf("%s\n", "<hidden>")
		} else {
			fmt.Printf("%s\n", v.AWSKey.Secret)
		}
		if v.AWSKey.MFA != "" {
			green.Printf("  MFA: ")
			fmt.Printf("%s\n", v.AWSKey.MFA)
		}
		if v.AWSKey.Role != "" {
			green.Printf("  Role: ")
			fmt.Printf("%s\n", v.AWSKey.Role)
		}
	} else {
		print("  [Empty]")
	}
}

func printSSHKeys(v *vaulted.Vault) {
	color.Cyan("\nSSH Keys:")
	if len(v.SSHKeys) > 0 {
		keys := []string{}
		for key := range v.SSHKeys {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			green.Printf("  %s\n", key)
		}
	} else {
		print("  [Empty]")
	}
}

func printDuration(v *vaulted.Vault) {
	cyan.Println("\nEnvironment:")
	green.Print("  Duration: ")
	var duration time.Duration
	if v.Duration == 0 {
		duration = vaulted.STSDurationDefault
	} else {
		duration = v.Duration
	}
	fmt.Printf("%s\n", duration.String())
}

func (e *Edit) readMenu(message string) string {
	if e.rlMenu == nil {
		var err error
		e.rlMenu, err = readline.New("")
		if err != nil {
			panic(err)
		}
	}

	print("")
	input := e.readInput(color.BlueString(message), e.rlMenu)
	print("")
	return input
}

func (e *Edit) readValue(message string) string {
	if e.rlValue == nil {
		var err error
		e.rlValue, err = readline.New("")
		if err != nil {
			panic(err)
		}
	}
	return e.readInput(color.GreenString(message), e.rlValue)
}

func (e *Edit) readInput(message string, rl *readline.Instance) string {

	rl.SetPrompt(message)
	line, err := rl.Readline()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(line)
}
