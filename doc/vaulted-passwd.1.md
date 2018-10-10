vaulted-passwd 1
================

NAME
----

vaulted passwd - changes the password of a vault

SYNOPSIS
--------

`vaulted passwd` *name*

`vaulted password` *name*

DESCRIPTION
-----------

Content in the *name* vault is untouched, only the password is changed.

If the `VAULTED_PASSWORD` environment variable is set, it will be used as the
current password for *name*, otherwise the user will be prompted for the
password.

If the `VAULTED_NEW_PASSWORD` environment variable is set, it will be used as
the new password for *name*, otherwise the user will be prompted for the
password.
