vaulted-exec 1
==============

NAME
----

vaulted exec - executes shell commands with a given vault or role

SYNOPSIS
--------

`vaulted exec` [*OPTIONS*] *name* *commands*  
`vaulted exec` [*OPTIONS*] *name* `--` *commands*  
`vaulted exec` [*OPTIONS*] `--assume` *arn* `--` *commands*

DESCRIPTION
-----------

Executes shell commands with a given vault or role. Exec combines the features
from `vaulted shell` with the command execution of `vaulted -n`.

OPTIONS
-------

`--assume` *arn*
  Specifies the full ARN or the role name of the role to assume. See
  **ASSUMING A ROLE** below for details on how Vaulted assumes roles.

  Role assumption may be performed without specifying a vault to spawn from.
  When invoked this way, credentials are sourced from default locations (e.g.
  environment, configuration files, instance profile, etc.).

`--no-session`
  Disables the generation of temporary credentials and role assumption. The
  permanent credentials stored in the vault are used instead.

  Role assumption can be performed after spawning a shell using the `--assume`
  command with the ARN of the role you wish to assume.

`--refresh`
  Start a new session with new temporary credentials and a refreshed expiration.

`--ssh-generate-key`
  Generate and load an RSA key into the spawned session's SSH agent.

`--ssh-proxy-agent`
  Exposes the external SSH agent to the spawned session.

`--ssh-signing-url`
  Configures the endpoint to use for SSH key signing. Currently only configured
  to work with HashiCorp Vault endpoints.

`--ssh-signing-users`
  Configures the users to sign SSH keys for when key signing is enabled.

AWS KEY
-------

[comment]: # (WHEN MODIFYING THESE DOCS, ALSO UPDATE: vaulted-env.1.md and
vaulted-shell.1.md)

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

[comment]: # (WHEN MODIFYING THESE DOCS, ALSO UPDATE: vaulted-env.1.md and
vaulted-shell.1.md)

A role to assume can be specified either in a vault's configuration (via
`vaulted edit`) or specified via the `--assume` option.

Vaulted first generates and caches a session for the vault (assuming a vault
*name* was specified) and uses the resulting session to perform the assume role
call. This allows the session to be cached for a longer duration than the assume
roll call will allow (maximum of 1 hour). This also means that the assume role
call is performed each time a new environment is spawned from the session.

You can also specify a role name and Vaulted will attempt
to create the full ARN from it, like so:

```
vaulted exec --assume SuperRole -- command
```

which would resolve to the full arn `arn:aws:iam::111222333444:role/SuperRole`.

When assuming a role, the maximum duration of the resulting credentials is 1
hour, regardless of the duration specified in the vault's configuration.

In addition to the variables specified above, Vaulted provides additional
environment variables with information about the role:

 * `VAULTED_ENV_ROLE_ACCOUNT_ID`  
   The account ID of the role assumed.
 * `VAULTED_ENV_ROLE_ARN`  
   The full ARN of the role assumed.
 * `VAULTED_ENV_ROLE_NAME`  
   The name of the role assumed.
 * `VAULTED_ENV_ROLE_PARTITION`  
   The AWS partition of the role assumed.
 * `VAULTED_ENV_ROLE_PATH`  
   The path of the role assumed.

For example:

```
vaulted exec --assume arn:aws:iam::111222333444:role/path/SuperRole -- command
```

would result in the following variables being set:

```
VAULTED_ENV_ROLE_ACCOUNT_ID=111222333444
VAULTED_ENV_ROLE_ARN=arn:aws:iam::111222333444:role/path/SuperRole
VAULTED_ENV_ROLE_NAME=SuperRole
VAULTED_ENV_ROLE_PARTITION=aws
VAULTED_ENV_ROLE_PATH=/path/
```

SSH KEY SIGNING
---------------

If you have access to a HashiCorp Vault instance that is configured for SSH key
signing, you can configure the Vaulted session to automatically sign all keys
added to the agent.

All keys added to the agent will be signed for a one hour duration and will be
automatically renewed prior to expiration to ensure that the keys are always
valid for the duration of the vault, so long as the HashiCorp Vault token
remains valid.

Key signing is particularly useful when combined with SSH key generation which
will generate a per-session, in-memory key to seed the agent that can be signed
on creation.

If the `VAULT_ADDR` environment variable is set you can provide a relative path
to the vault signing URL instead of the full address.

You must have a HashiCorp Vault token configured for this work, either via a
`.vault-token` file located in your home directory or a `VAULT_TOKEN`
environment variable. Using `vault login` with the appropriate options can be a
simple way to populate the `.vault-token` file.

For more details on the SSH proxy agent that Vaulted employs, see
`https://github.com/miquella/ssh-proxy-agent`

GUI Password Prompts
--------------------

GUI-based password prompts can be used by setting the `VAULTED_ASKPASS`
variable. See `vaulted(1)` for more details.
