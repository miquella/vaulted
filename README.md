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
