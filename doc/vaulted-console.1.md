vaulted-console 1
=============

NAME
----

vaulted console - Opens the AWS console in the default web browser

SYNOPSIS
--------

`vaulted console`  
`vaulted console` *name* [*OPTIONS*]  
`vaulted console --assume` *arn* [*OPTIONS*]  
`vaulted console --duration` *duration* [*OPTIONS*]

DESCRIPTION
-----------

Opens the AWS console in the default web browser. Uses either the credentials in the current environment or the credentials in the specified vault. Console sessions either use the passed in duration, the provided vault's duration, or defaults to 1 hour.

Durations must be at least 15 minutes and less than 12 hours.

OPTIONS
-------

`--assume` *arn*
  Specifies the full ARN of the role to assume. See **ASSUMING A ROLE** below
  for details on how Vaulted assumes roles.

  Role assumption may be performed without specifying a vault to spawn from.
  When invoked this way, credentials are sourced from default locations (e.g.
  environment, configuration files, instance profile, etc.).

`--duration` *duration*
	Specifies the duration that the console session is valid. The duration must be within the range 15m-12h.

ASSUMING A ROLE
---------------

A role to assume can be specified either in a vault's configuration (via
`vaulted edit`) or specified via the `--assume` option.

Vaulted first opens the specified vault to retrieve the appropriate credentials. If a role is specified in the vault's configuration it will use that unless a role is explicitly passed in through the `--assume` option. 