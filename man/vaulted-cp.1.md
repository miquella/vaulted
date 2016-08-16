vaulted-cp 1
=============

NAME
----

vaulted cp - copies the content of a vault and saves it as a new vault with a new password

SYNOPSIS
--------

`vaulted cp` *&lt;old&gt;* *&lt;new&gt;*

`vaulted copy` *&lt;old&gt;* *&lt;new&gt;*

DESCRIPTION
-----------

Content in the *new* vault is created or replaced by content from *old*.

If the `VAULTED_PASSWORD` environment variable is set, it will be used as the
password for *old*, otherwise the password will be requested via the tty.

If the `VAULTED_NEW_PASSWORD` environment variable is set, it will be used as
the password for *new*, otherwise the password will be requested via the tty.
