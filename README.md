# vaulted
Spawning and storage of secure environments

With so many secrets floating around in our modern lives, it's a wonder we're
able to keep track of any of them!

`vaulted` allows you to create vaults of related secrets and then spawn
environments with these secrets. Vaults can contain secure environment
variables, AWS credentials, or SSH keys (RSA & ECDSA).

`vaulted` also attempts to insulate spawned environments from other
environments on the system. Temporary AWS credentials are created for each
spawned environment as well as a new SSH agent. The SSH agent still allows
access keys in the parent environment's SSH agent, but any keys added inside
the environment are only available in the spawned environment.

# Usage

When decrypting a vault, `vaulted` will prompt for the vault's password.
Optionally, the password may be provided through the `VAULTED_PASSWORD`
environment variable.

## `vaulted add`

Interactively add a new environment named `NAME`.

```sh
vaulted add NAME
```

## `vaulted edit`

Interactively edit an environment named `NAME`.

```sh
vaulted edit NAME
```

## `vaulted cat`

Display the static environment variables defined in `NAME`.

```sh
vaulted cat NAME
```

## `vaulted dump`

Dumps a JSON formatted version of the `NAME` vault.

```sh
vaulted dump NAME
```

## `vaulted list` / `vaulted ls`

Lists all vaults.

```sh
vaulted list
```

## `vaulted load`

Creates or replaces the `NAME` vault with the JSON formatted version provided.

```sh
vaulted load NAME
```

## `vaulted rm`

Removes the `NAME` vault.

```sh
vaulted rm NAME
```

## `vaulted shell`

Spawns an interactive shell using the `NAME` vault's secrets.

The `SHELL` environment variable is used as the interactive shell to spawn.

```sh
vaulted shell NAME
```

## `vaulted upgrade`

Upgrades all vaults in the legacy vault format to the current format.

```sh
vaulted upgrade
```

# Spawned Environment

In addition to secrets, spawned environments also include the name of the
spawned environment in the `VAULTED_ENV` environment variable. This is
particularly useful if you would like to indicate the spawned environment in
your shell prompt.
