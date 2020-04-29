# gomu - Go Mod Utils #
Designed to make working with mod files easier.


# Commands #
Commands are actions taken on the dependency chain.

This is usually parsed as the first non-flag argument passed to the utility.


## Help Commands ##
Help commands will print usage instructions and details.

Running help commands with arguments and flags provides more specific contextual details.

### gomu ###
  :: Designed to make working with mod files easier.
  To learn more, run `gomu help` or `gomu help <command>`
  (Flags can be added to either help command)

### gomu help ###
  :: Prints available commands and flags.
  Use `gomu help <command> <flags>` to get more specific info.

### gomu version ###
  :: Prints current version.
  Install using `gomu upgrade` to get version support.

### gomu upgrade ###
  :: Updates gomu itself!
  Optionally accepts a version number.
  Without argument, updates to latest tag.
  Otherwise updates to latest branch/tag provided by first arg or -b.
  Usage: `gomu upgrade` or `gomu upgrade -b master` or `gomu upgrade v0.5.1`


## Local Commands ##
Local commands can/will make stashes and edits to local files on your working copies.

However, they will not attempt to commit or push any changes by themselves.

### gomu list ###
  :: Prints each file in dependency chain.

### gomu pull ###
  :: Updates branch for file in dependency chain.
  Providing a -branch will checkout given branch.
  Creates branch if provided none exists.

### gomu replace ###
  :: Replaces each versioned file in the dependency chain.
  Uses the current checked out local copy.

### gomu reset ###
  :: Reverts go.mod and go.sum back to last committed version.
  Usage: `gomu reset mod-common parg`

### gomu test ###
  :: Runs `go test` on each library in the dependency chain.
  Prints names of failing libraries.
  Usage: `gomu test mod-common`


## Destructive Commands ##
Destructive commands can/will attempt to commit and push changes.

If running with -name-only, it will NOT prompt you for a warning.

Please be careful!

### gomu sync ###
  :: Updates modfiles.
  Conditionally performs extra tasks depending on flags.
  Usage: `gomu <flags> sync mod-common parg simply <flags>`

### gomu workflow ###
  :: Adds a github workflow to a repo.
  Requires -source <template path>.
  Usage: `gomu workflow mod-utils -c -b new-workflow -source workflows/templates/autotag.yml`

### gomu secret ###
  :: Adds a secret to a repo's github actions.
  Requires -source <file containing secret>.
  Usage: `gomu secret mod-utils -source ~/.ssh/server_key.crt`

# Flags #
Flags are options that can be set for some commands. 

Flags can be added before or after command name.

However, it is best not to split up flags and commands/args to avoid any accidental commands.

### [-i -in -include] ###
  :: Will aggregate files in 1 or more directories.
  Usage: `gomu list -i hatchify -i vroomy`

### [-b -branch] ###
  :: Will checkout or create said branch
  Updating or creating a pull request.
  Depending on command and other flags.
  Usage: `gomu pull -b feature/Jira-Ticket`

### [-name -name-only] ###
  :: Will reduce output to just the filenames changed.
  (ls-styled output for | chaining)
  Usage: `gomu list -name`

### [-direct -direct-import] ###
  :: Will avoid recursion in dependency sorting.
  Only includes deps in go.mod (not go.sum).
  Usage: `gomu list mod-utils -direct`

### [-c -commit] ###
  :: Will commit local changes if present.
  Includes all changed files in repository.
  Usage: `gomu sync -c`

### [-pr -pull-request] ###
  :: Will create a pull request if possible.
  Fails if on master, or if no changes.
  Usage: `gomu sync -pr`

### [-m -msg -message] ###
  :: Will set a custom commit message.
  Applies to -c and -pr flags.
  Usage: `gomu sync -c -m "Update all the things!"`

### [-t -tag] ###
  :: Will increment tag if new commits since last tag.
  Requires tag previously set.
  Usage: `gomu sync -t`

### [-set -set-version] ###
  :: Can be used with -tag to update semver.
  Will force tag version for all deps in chain.
  Usage: `gomu sync -t -set v0.5.0`

### [-source -source-path] ###
  :: Required for workflow and secret commands.
  Will provide a source template or secret file.
  Usage: `gomu workflow mod-utils -source path/to/template.yml`
