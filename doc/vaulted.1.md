vaulted 1
=========

NAME
----

vaulted - spawn sessions from securely stored secrets

SYNOPSIS
--------

`vaulted` `-n` *name* [`-i`]  
`vaulted` `-n` *name* [`--`] *CMD*

`vaulted` *COMMAND* [*args...*]

DESCRIPTION
-----------

If no *COMMAND* is provided, `vaulted` either spawns *CMD* (if provided) or
spawns an interactive shell.

`--` may be used to differentiate the *CMD* from `vaulted`'s own arguments.

COMMANDS
--------

`add` / `create` / `new`
  Interactively creates the content of a new vault. See vaulted-add(1).

`cp` / `copy`
  Copies the content of a vault and saves it as a new vault with a new password. See vaulted-cp(1).

`dump`
  Writes the content of a vault to stdout as JSON. See vaulted-dump(1).

`edit`
  Interactively edits the content of an existing vault. See vaulted-edit(1).

`env`
  Outputs shell commands that load secrets for a vault into the shell. See vaulted-env(1).

`load`
  Uses JSON provided to stdin to create or replace the content of a vault. See vaulted-load(1).

`ls` / `list`
  Lists all vaults. See vaulted-ls(1).

`rm` / `delete` / `remove`
  Removes existing vaults. See vaulted-rm(1).

`shell`
  Starts an interactive shell with the secrets for the vault loaded into the shell. See vaulted-shell(1).

`upgrade`
  Upgrades legacy vaults to the current vault format. See vaulted-upgrade(1).

EXIT CODES
----------
|Exit code|Meaning|
|:-:|---|
| 0 | Success. |
| 64 | Invalid CLI usage (see message for more details). |
| 65 | There was an unrecoverable problem with the vault file. |
| 69 | A required service is presently unavailable (e.g. askpass). |
| 79 | Invalid password supplied. |

GUI Password Prompts
--------------------

Although Vaulted tries to make sure you can redirect `stdin` and friends,
sometimes it is still preferable to use a GUI-based password prompt. For this
reason, Vaulted can be configured to use an askpass implementation. Vaulted's
askpass integration is triggered when the `VAULTED_ASKPASS` variable is set.

Pointing `VAULTED_ASKPASS` to an executable file that implements askpass will
cause Vaulted to use execute the file specified to prompt the user for
passwords. The first parameter provided to the executable is prompt text
intended to be shown to the user. The askpass implementation then writes the
password to `stdout` and returns a success code (0). If a failure code (non-0)
is returned, the password input is aborted.

Vaulted is intended to integrate seamlessly with existing askpass
implementations (e.g. `ssh-askpass`).

On macOS, a simple AppleScript askpass implementation can be used:

```AppleScript
#!/usr/bin/env osascript

on run argv
    if length of argv > 0 then
        set message to item 1 of argv
    else
        set message to "Password:"
    end if

    set frontmost_application to name of (info for (path to frontmost application))
    tell application frontmost_application
        display dialog message with title "Vaulted" with icon caution with hidden answer default answer ""

        text returned of result
    end tell
end run
```
