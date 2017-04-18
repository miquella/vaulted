vaulted 1
=========

NAME
----

vaulted - spawn environments from securely stored secrets

SYNOPSIS
--------

`vaulted` `-n` *name* [`-i`]  
`vaulted` `-n` *name* [`--`] *CMD*

`vaulted` *COMMAND* [*args...*]

DESCRIPTION
-----------

If no *COMMAND* is provided, `vaulted` either spawns *CMD* (if provided) or
spawns an interactive shell.

`--` may be used to differentiate the *CMD* from `vaulted`'s own arguments.

COMMANDS
--------

`add`
  Interactively creates the content of a new vault. See vaulted-add(1).

`cp` / `copy`
  Copies the content of a vault and saves it as a new vault with a new password. See vaulted-cp(1).

`dump`
  Writes the content of a vault to stdout as JSON. See vaulted-dump(1).

`edit`
  Interactively edits the content of an existing vault. See vaulted-edit(1).

`env`
  Outputs shell commands that load secrets for a vault into the shell. See vaulted-env(1).

`load`
  Uses JSON provided to stdin to create or replace the content of a vault. See vaulted-load(1).

`ls` / `list`
  Lists all vaults. See vaulted-ls(1).

`rm`
  Removes existing vaults. See vaulted-rm(1).

`shell`
  Starts an interactive shell with the secrets for the vault loaded into the shell. See vaulted-shell(1).

`upgrade`
  Upgrades legacy vaults to the current vault format. See vaulted-upgrade(1).

EXIT CODES
----------
|Exit code|Meaning|
|:-:|---|
| 0 | Success. |
| 64 | Invalid CLI usage (see message for more details). |
| 65 | There was an unrecoverable problem with the vault file. |
| 79 | Invalid password supplied. |
