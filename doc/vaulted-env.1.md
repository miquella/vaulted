vaulted-env 1
=============

NAME
----

vaulted env - outputs shell commands that load secrets for a vault into the shell

SYNOPSIS
--------

`vaulted env` *name*

DESCRIPTION
-----------

Creates an environment using variables and the AWS key stored in the vault. The
resulting environment variables are used to generate shell commands that will
load the secrets into a shell.

The shell is autodetected from the `SHELL` environment variable, if the shell
is unknown or unspecified, `sh` compatible commands are emitted.

_Note:_ SSH keys are ignored when generating environments this way. This is due
to the inability to track the lifetime of the environment, which means the SSH
agent would exist indefinitely.
