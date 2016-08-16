vaulted-rm 1
=============

NAME
----

vaulted rm - copies the content of a vault and saves it as a new vault with a new password

SYNOPSIS
--------

`vaulted rm` *&lt;name&gt;* *...*

DESCRIPTION
-----------

Removes the vaults specified by *name*. The exit code is equal to the number of
vaults that could not be removed.

If a vault cannot be removed, the error is displayed on stdout.
