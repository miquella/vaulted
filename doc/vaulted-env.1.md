vaulted-env 1
=============

NAME
----

vaulted env - outputs shell commands that load secrets for a vault into the shell

SYNOPSIS
--------

`vaulted env` *name* [*OPTIONS*]  
`vaulted env --assume` *arn* [*OPTIONS*]  
`vaulted env --no-session` *name* [*OPTIONS*]

DESCRIPTION
-----------

Creates a session using variables and the AWS key stored in the vault. The
resulting environment variables are used to generate shell commands that will
load the secrets into a shell.

The shell is autodetected from the `SHELL` environment variable, if the shell
is unknown or unspecified, `sh` compatible commands are emitted.

*Note:* SSH keys are ignored when generating sessions this way. This is due to
the inability to track the lifetime of the session, which means the SSH agent
would exist indefinitely.

OPTIONS
-------

`--assume` *arn*
  Specifies the full ARN or short name of the role to assume. See
  **ASSUMING A ROLE** below for details on how Vaulted assumes roles.

  Role assumption may be performed without specifying a vault to spawn from.
  When invoked this way, credentials are sourced from default locations (e.g.
  environment, configuration files, instance profile, etc.).

`--no-session`
  Disables the generation of temporary credentials and role assumption. The
  permanent credentials stored in the vault are used instead.

  Role assumption can be performed after spawning a shell using the `--assume`
  command with the ARN of the role you wish to assume.

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

  As an example, to output just the AWS credentials in a JSON format, use the
  custom template: `--format '{{ json .AWSCreds }}'`

`--refresh`
  Start a new session with new temporary credentials and a refreshed expiration.

|Field Name|Type|Notes|
|---|---|---|
| Command | string | The full command line used to invoke vaulted including options |
| AWSCreds | struct | The AWS access key for the sesion |
| AWSCreds.ID | string | ID for the AWS access key |
| AWSCreds.Secret | string | Secret for the AWS access key |
| AWSCreds.Token | string | Session token for the AWS access key (used with temporary credentials) |
| Set | map[string]string | A map of environment variables to be set along with their values |
| Unset | []string | A slice of variables to be removed from the environment before setting new values |

AWS KEY
-------

[comment]: # (WHEN MODIFYING THESE DOCS, ALSO UPDATE: vaulted-shell.1.md)

Vaulted uses permanent credentials stored in the vault to generate temporary
credentials using AWS STS. The resulting credentials are configured to last for
the duration of the vault session (configured via `vaulted edit`).

Temporary credentials must be used in order to assume a role.

*Note:* because Vaulted defaults to substituting permanent credentials with
temporary credentials, the credentials you provide to Vaulted may not match
those in a spawned session.

Vaulted uses AWS standard environment variables for credentials:

 * `AWS_ACCESS_KEY_ID`  
   An identifier for the access key. Unique to each set of permanent or
   temporary credentials.
 * `AWS_SECRET_ACCESS_KEY`  
   The secret used to sign requests. Unique to each set of permanent or
   temporary credentials.
 * `AWS_SESSION_TOKEN` / `AWS_SECURITY_TOKEN`  
   Provided when using temporary credentials. `AWS_SECURITY_TOKEN` is provided
   to support legacy SDKs.

ASSUMING A ROLE
---------------

[comment]: # (WHEN MODIFYING THESE DOCS, ALSO UPDATE: vaulted-shell.1.md)

A role to assume can be specified either in a vault's configuration (via
`vaulted edit`) or specified via the `--assume` option.

Vaulted first generates and caches a session for the vault (assuming a vault
*name* was specified) and uses the resulting session to perform the assume role
call. This allows the session to be cached for a longer duration than the assume
roll call will allow (maximum of 1 hour). This also means that the assume role
call is performed each time a new environment is spawned from the session.

You can also specify a role name and Vaulted will attempt to create the full ARN
from it, like so:

```
vaulted shell --assume SuperRole
```

which would resolve to the full arn `arn:aws:iam::111222333444:role/SuperRole`.

When assuming a role, the maximum duration of the resulting credentials is 1
hour, regardless of the duration specified in the vault's configuration.

In addition to the variables specified above, Vaulted provides additional
environment variables with information about the role:

 * `VAULTED_ENV_ROLE_ARN`  
   The full ARN of the role assumed.
 * `VAULTED_ENV_ROLE_ACCOUNT_ID`  
   The account ID specified by the role ARN.
 * `VAULTED_ENV_ROLE_NAME`  
   The name of the role specified by the role ARN.
 * `VAULTED_ENV_ROLE_PATH`  
   The path of the role specified by the role ARN.

For example:

```
vaulted shell --assume arn:aws:iam::111222333444:role/SuperRole
```

would result in the following variables being available in {{.Set}}:

```
VAULTED_ENV_ROLE_ARN=arn:aws:iam::111222333444:role/SuperRole
VAULTED_ENV_ROLE_ACCOUNT_ID=111222333444
VAULTED_ENV_ROLE_NAME=SuperRole
VAULTED_ENV_ROLE_PATH=/
```

GUI Password Prompts
--------------------

GUI-based password prompts can be used by setting the `VAULTED_ASKPASS`
variable. See `vaulted(1)` for more details.
