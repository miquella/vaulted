package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/miquella/vaulted/vault"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
)

type CommandMode int

const (
	NoMode CommandMode = iota
	SpawnEnvironmentMode
	MutateEnvironmentMode
	ListEnvironmentsMode
)

var (
	// vault flags
	filename string

	// environment flags
	environment string
	commandMode CommandMode = NoMode

	listEnvironments    bool
	dumpJsonEnvironment bool

	interactiveAdd    bool
	deleteEnvironment bool
	interactiveShell  bool

	charDevice bool

	// vault
	vaultKey  []byte
	vaultData []byte
	v         vault.Vault
	envs      vault.Environments
)

func init() {
	environment = os.Getenv("VAULTED_ENV")

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&filename, "f", filepath.Join(u.HomeDir, ".vaulted"), "vault filename")
	flag.StringVar(&environment, "n", environment, "name of the environment")
	flag.BoolVar(&listEnvironments, "L", false, "list environments in vault")
	flag.BoolVar(&dumpJsonEnvironment, "j", false, "dump json version of environment")
	flag.BoolVar(&interactiveAdd, "a", false, "add to environment interactively")
	flag.BoolVar(&deleteEnvironment, "D", false, "delete environment from vault")
	flag.BoolVar(&interactiveShell, "i", false, "spawn a new shell populated with the environment")
	flag.Parse()

	// figure out which mode the user wants
	spawnMode := interactiveShell || flag.NArg() > 0
	mutateMode := interactiveAdd || deleteEnvironment
	listMode := listEnvironments || dumpJsonEnvironment

	if spawnMode {
		if listMode || mutateMode {
			fmt.Fprint(os.Stderr, "ERROR: cannot list/dump or update environments while spawning\n")
			os.Exit(2)
		}

		if os.Getenv("SHELL") == "" {
			fmt.Fprint(os.Stderr, "ERROR: SHELL environment variable not set, cannot spawn shell")
			os.Exit(2)
		}

		commandMode = SpawnEnvironmentMode

	} else if mutateMode {
		if listMode {
			fmt.Fprintf(os.Stderr, "ERROR: cannot update environments while listing/dumping\n")
			os.Exit(2)
		}
		commandMode = MutateEnvironmentMode

	} else if listMode {
		commandMode = ListEnvironmentsMode
	}

	stat, _ := os.Stdin.Stat()
	charDevice = (stat.Mode() & os.ModeCharDevice) != 0
}

func main() {
	if commandMode == NoMode {
		fmt.Fprint(os.Stderr, "invalid command mode\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	err := loadVault()
	if err != nil {
		println(err.Error())
		os.Exit(3)
	}

	switch commandMode {
	case ListEnvironmentsMode:
		listEnvironmentsMode()
	case SpawnEnvironmentMode:
		spawnEnvironmentMode()
	case MutateEnvironmentMode:
		mutateEnvironmentMode()
		err = saveVault()
		if err != nil {
			println(err.Error())
			os.Exit(3)
		}
	}
}

func getPassword() string {
	password, err := speakeasy.Ask("Vault password: ")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(4)
	}

	return password
}

func loadVault() error {
	var err error

	// skip loading if the file doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		envs = make(vault.Environments)
		return nil
	}

	// read the vault data
	vaultData, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// deserialize the vault data
	err = json.Unmarshal(vaultData, &v)
	if err != nil {
		return err
	}

	// decrypt the vault
	if vaultKey == nil {
		vaultKey, err = v.GenerateKey(getPassword())
		if err != nil {
			return err
		}
	}

	envs, err = v.DecryptEnvironments(vaultKey)
	return err
}

func saveVault() error {
	var err error

	// encrypt the vault
	if vaultKey == nil {
		vaultKey, err = v.GenerateKey(getPassword())
		if err != nil {
			return err
		}
	}

	err = v.EncryptEnvironments(vaultKey, envs)
	if err != nil {
		return err
	}

	// serialize the vault data
	tmpVaultData, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	// write the vault data
	if !bytes.Equal(tmpVaultData, vaultData) {
		return ioutil.WriteFile(filename, tmpVaultData, 0600)
	}

	return nil
}

func listEnvironmentsMode() {
	if listEnvironments {
		envList := make([]string, 0, len(envs))
		for name := range envs {
			envList = append(envList, name)
		}

		sort.Strings(envList)
		for _, name := range envList {
			fmt.Fprintf(os.Stdout, "%s\n", name)
		}
	} else {
		env, exists := envs[environment]
		if !exists {
			fmt.Fprintf(os.Stderr, "ERROR: environment '%s' does not exist\n", environment)
			os.Exit(2)
		}

		data, err := json.MarshalIndent(&env, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "%s\n", data)
	}
}

func spawnEnvironmentMode() {
	// validate that the environment exists
	env, exists := envs[environment]
	if !exists {
		fmt.Fprintf(os.Stderr, "ERROR: environment '%s' does not exist\n", environment)
		os.Exit(2)
	}

	// compile command line arguments
	var spawnArgs []string
	if interactiveShell {
		spawnArgs = append(spawnArgs, os.Getenv("SHELL"), "--login")
	}
	spawnArgs = append(spawnArgs, flag.Args()...)

	if len(spawnArgs) < 1 {
		fmt.Fprint(os.Stderr, "ERROR: invalid spawn arguments\n")
		os.Exit(2)
	}

	// build the environment
	vars := ParseEnviron(os.Environ())
	for key, val := range env.Vars {
		vars[key] = val
	}

	// locate the executable
	fullpath, err := exec.LookPath(spawnArgs[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot find executable: %s - %s\n", spawnArgs[0], err)
		os.Exit(2)
	}

	// start the process
	var attr os.ProcAttr
	attr.Env = CreateEnviron(vars)
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	proc, err := os.StartProcess(fullpath, spawnArgs, &attr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to execute command: %s\n", err)
		os.Exit(3)
	}

	// wait for the process to exit
	state, err := proc.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to execute command: %s\n", err)
		os.Exit(3)
	}
	if !state.Success() {
		os.Exit(3)
	}
}

func mutateEnvironmentMode() {
	if deleteEnvironment {
		delete(envs, environment)
	} else if interactiveAdd {
		env := envs[environment]
		env.Name = environment
		if env.Vars == nil {
			env.Vars = make(map[string]string)
		}

		if charDevice {
			fmt.Printf("input vars (VAR=VALUE); end with ctrl+d\n")
		}
		for {
			if charDevice {
				fmt.Print("var> ")
			}

			var varLine string
			_, err := fmt.Scanln(&varLine)
			if err == io.EOF || varLine == "" {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Failed to read env vars - %s\n", err)
				os.Exit(1)
			}

			key, value := ParseVar(varLine)
			env.Vars[key] = value
		}
		envs[environment] = env
	}
}
