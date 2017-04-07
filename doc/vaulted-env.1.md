vaulted-env 1
=============

NAME
----

vaulted env - outputs shell commands that load secrets for a vault into the shell

SYNOPSIS
--------

`vaulted env` *name* [*OPTIONS*]

DESCRIPTION
-----------

Creates an environment using variables and the AWS key stored in the vault. The
resulting environment variables are used to generate shell commands that will
load the secrets into a shell.

The shell is autodetected from the `SHELL` environment variable, if the shell
is unknown or unspecified, `sh` compatible commands are emitted.

*Note:* SSH keys are ignored when generating environments this way. This is due
to the inability to track the lifetime of the environment, which means the SSH
agent would exist indefinitely.

OPTIONS
-------

`--format` &lt;shell,fish,sh,json,*custom*&gt;
  Specify what format to use, defaults to `shell` which will autodetect which
  shell format to emit.

  `fish`, and `sh` will emit output suitable for evaluation into the specified
  shell. `json` will emit the map of values to be set in a format useful in
  other programming environments, this will not include the array of variables
  to be removed.

  If the format string supplied does not match one of the pre-defined formats
  it will be interpreted as a text template using golang's built-in
  `text/template` package. The user template will be supplied a struct with the
  fields outlined below. Also available are a `replace` function for gsub like
  behavior and a `json` function which accepts a value and emits a pretty
  printed JSON string.

  To receive just the names of the variables to remove formatted as JSON the
  custom template would be: `{{ json .Unset }}`

|Field Name|Type|Notes|
|---|---|---|
| Set | map[string]string | A map of environment variables to be set along with their values |
| Unset | []string | A slice of variables to be removed from the environment before setting new values |
| Command | string | The full command line used to invoke vaulted including options |

EXIT CODES
----------
|Exit code|Meaning|
|:-:|---|
| 0 | Complete success. |
| 64 | Invalid CLI usage (see message for more details).|
| 65 | There was an unrecoverable problem with the vault file. |
| 79 | The supplied password is incorrect.|
