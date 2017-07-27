vaulted-rm 1
=============

NAME
----

vaulted rm - removes existing vaults

SYNOPSIS
--------

`vaulted rm` *name* *...*

`vaulted delete` *name* *...*  
`vaulted remove` *name* *...*

DESCRIPTION
-----------

Removes the vaults specified by *name*. The exit code is equal to the number of
vaults that could not be removed.

If a vault cannot be removed, the error is displayed on stdout.
