vaulted
=======

Spawn environments from securely stored secrets.

With so many secrets floating around in our modern lives, it's a wonder we're
able to keep track of any of them!

`vaulted` allows you to create vaults of related secrets and then spawn
sessions with these secrets. Vaults can contain secure environment
variables, AWS credentials, or SSH keys (RSA, DSA, & ECDSA).

`vaulted` also attempts to insulate spawned environments from other
environments on the system. Temporary AWS credentials are created for each
session as well as a new SSH agent. The SSH agent still allows access keys in
the parent environment's SSH agent, but any keys added inside the spawned
environment are only available in the that environment.

Installation
------------

### macOS

The easiest way to install `vaulted` on macOS is through
[Homebrew](http://brew.sh/).

```sh
brew install vaulted
```

### Manual

Installation on other platforms should be simple enough through `go get` as
long as you have a proper Go environment setup:

```sh
go get -u github.com/miquella/vaulted
```

Don't forget to add `$GOPATH/bin` to your `$PATH`!

Getting Started
---------------

`vaulted` is oriented around vaults of secrets that are used to spawn
environments. To get started, add a new vault:

```sh
vaulted add my-vault
```

This will start an interactive editing mode that will help you create your
first vault. AWS keys, SSH keys, and arbitrary environment variables can be
added to the vault. Once you have your vault arranged how you would like, use
`q` to exit the interactive mode and save the vault to disk.

While editing a vault, `Ctrl+C` may be used to discard changes to the
vault.

Now that your vault has been saved, the list of vaults will reflect your newly
saved vault:

```sh
vaulted ls
```

And you can use `vaulted` to spawn a command in an environment generated from
the secrets stored in the vault:

```sh
vaulted -n my-vault -- aws s3 ls
```

Sometimes it is useful to be able to issue multiple commands that require the
vault's secrets. In this case, you can spawn an interactive shell:

```sh
vaulted shell my-vault
```

_**Warning!** Leaving interactive shells with your credentials loaded can be
dangerous as you may inadvertently provide credentials to an application you
didn't intend!_

Using Vaulted from other software
---------------------------------
An `env` subcommand has been included with the intention of supplying machine
readable output for integration with shells and shell utilities. Every effort
has been made to supply meaningful return codes on failures along with a
description of what has gone wrong. See `vaulted(1)` for details on the return
values to expect and their meanings.


Going Further
-------------

While `vaulted` supports basic modification methods like copying, editing, and
removing, more advanced methods such as JSON-formatted dumping and loading are
also available. An environment can even be loaded into a running shell! See
`vaulted --help` for available commands.

Spawned Environment
-------------------

In addition to including secrets stored in the vault, spawned environments also
include environment variables that describe how the session and environment
were spawned. See `vaulted-env(1)` and `vaulted-shell(1)` for details.

GUI Password Prompts
--------------------

GUI-based password prompts can be used by setting the `VAULTED_ASKPASS`
variable. See [`vaulted(1)`](doc/vaulted.1.md) for more details.
