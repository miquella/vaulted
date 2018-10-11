vaulted-shell 1
===============

NAME
----

vaulted shell - starts an interactive shell with the secrets for the vault loaded into the shell

SYNOPSIS
--------

`vaulted shell` *name* [*OPTIONS*]  
`vaulted shell --assume` *arn* [*OPTIONS*]  
`vaulted shell --assume-duration` *duration* [*OPTIONS*]  
`vaulted shell --no-session` *name* [*OPTIONS*]

DESCRIPTION
-----------

Starts an interactive shell (uses the `SHELL` environment variable, if set;
otherwise defaults to `/bin/sh`).

OPTIONS
-------

`--assume` *arn*
  Specifies the full ARN or the role name of the role to assume. See
  **ASSUMING A ROLE** below for details on how Vaulted assumes roles.

  Role assumption may be performed without specifying a vault to spawn from.
  When invoked this way, credentials are sourced from default locations (e.g.
  environment, configuration files, instance profile, etc.).

`--assume-duration` *duration*
  Specifies the duration that a role should be assumed for (e.g. 15m or 2h)
  See **ASSUMING A ROLE** below for details on how Vaulted assumes roles.

  This can apply to the role specified in a vault's configuration (via
  `vaulted edit`) or a role specified via the `--assume` option. If multiple
  roles are assumed, this will be applied to the final assumed role.

`--no-session`
  Disables the generation of temporary credentials and role assumption. The
  permanent credentials stored in the vault are used instead.

  Role assumption can be performed after spawning a shell using the `--assume`
  command with the ARN of the role you wish to assume.

`--refresh`
  Start a new session with new temporary credentials and a refreshed expiration.

AWS KEY
-------

[comment]: # (WHEN MODIFYING THESE DOCS, ALSO UPDATE: vaulted-env.1.md and
vaulted-exec.1.md)

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
vaulted-exec.1.md)

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

would result in the following variables being set:

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
