vaulted-edit 1
==============

NAME
----

vaulted edit - interactively edits the content of an existing vault

SYNOPSIS
--------

`vaulted edit` *name*

DESCRIPTION
-----------

Spawns an interactve mode for editing the content of an existing vault.

Upon quitting, the new content is saved to the vault.

GLOBAL
------

* S - Show/Hide Secrets  
   Toggles whether secret values are displayed when viewing vault details
   in the edit menu. This setting is global and affects the entire editing
   session and is reset to hidden each time the vault is reopened for
   editing.

AWS KEY
-------

* k - Key  
   Manages Access Key ID and Secret Access Key credentials.
* m - MFA  
   Manages MFA (multi-factor authentication) ARN or serial number. If enabled,
   user will be prompted to enter MFA code when accessing vault. MFA will
   remain active for the vault duration.  
   If you are using STS credentials, which occurs by default (see the 't'
   option below for details), you must have MFA enabled to invoke any IAM calls.
* r - Role  
   ARN of role to be assumed when accessing vault.
   When assuming a role, the maximum duration allowed by AWS is 1 hour. If your
   duration is greater than 1 hour when setting a role, the duration will be
   adjusted to 1 hour.
* t - Substitute with temporary credentials  
   Toggles whether your AWS credentials are substituted with a set of temporary
   credentials. For more details on this process, see the documentation for
   `vaulted shell`.
* D - Delete  
   Removes all AWS details stored in the vault.

SSH KEYS
-------

* a - Add  
   Manages SSH keys loaded into the agent when a session starts.
* g - Generate Key  
   Generates a dynamic, in-memory SSH key that will be stored for the vault
   duration. This is particularly useful when coupled with SSH key signing.
* v - HashiCorp Vault Signing URL  
   The HashiCorp Vault URL to configure SSH key signing. This URL can be an
   absolute path, or if the VAULT_ADDR environment variable is set you may
   provide this as a relative path to join the two.
* u - HashiCorp Vault User Principals  
   Configures the user principals that keys will be signed for when SSH key
   signing is enabled through HashiCorp Vault.
* E - Expose external SSH agent  
   Toggles whether or not to expose an external SSH agent configured via
   the SSH_AUTH_SOCK environment variable to the spawned vault session.
* D - Delete  
   Removes a specified SSH key from the list.
