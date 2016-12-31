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

AWS KEY
-------

* k - Key  
   Manages Access Key ID and Secret Access Key credentials.
* m - MFA  
   Manages MFA (multi-factor authentication) ARN or serial number. If enabled, user will be prompted
   to enter MFA code when accessing vault. MFA will remain active for the
   vault duration.
* r - Role  
   ARN of role to be assumed when accessing vault.
* t - Substitute with temporary credentials  
   Toggles whether your AWS credentials are substituted with a set of temporary
   credentials. For more details on this process, see the documentation for
   `vaulted shell`.
* S - Show/Hide Key  
   Toggles whether Secret Access Key ID is displayed when viewing vault
   details in the edit menu. This setting only affects the current editing
   session and is reset to hidden each time the vault is reopened for editing.
* D - Delete  
   Removes all AWS details stored in the vault.