# vaulted

Allows the secure storage of environments and execution of processes in those environments

# Usage

To prevent having to pass the `-n` argument for each invocation, the `VAULTED_ENV` environment variable can be set to the name of the environment.

## interactive add

```sh
vaulted -n env_name -a
```

## delete account

```sh
vaulted -n env_name -D
```

## list accounts

```sh
vaulted -L
```

## execute command

```sh
vaulted -n env_name command arg1 arg2
```

## spawn interactive shell

```sh
vaulted -n env_name -i
```

# Environment Self-Awareness

Vaulted will set an environment variable in spawned environments for commands or interactive shells to allow the spawned process to be aware that it is within a Vaulted-spawned environment. This is particularly useful if you would like to indicate such status in your shell prompt.

By default this environment variable is named `VAULTED_ACTIVE_ENV`, but another name may be specified with the `-E` option. Or set the value to an empty string (`-E ''`) to prevent Vaulted from adding the extra environment variable.
