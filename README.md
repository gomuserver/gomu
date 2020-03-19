#gomu - Go Mod Utils
gomu is intended to make working with go.mod sane.

## Non Destructive ##
Non Destructive commands will not commit/push changes to the repository. However, they should not be interrupted.

# list: Sorted dependencies
List all dependencies in order of the dependency chain
`gomu list`

# pull: Update all libs
Pull all dependencies in order of the dependency chain
`gomu pull`

# replace-local: Add local replacement
Appends a local replacement clause to each lib's go.mod file, directing each dependency preceding lib in the chain to a local source
`gomu replace-local`

# help: Show usage
Show flags/args and usage help
`gomu help`

## Destrucive ##
Destructive commands can/will attempt to commit and push changes. If running with -name-only, it will NOT prompt you for a warning. Please be careful!

# sync: Update deps and tags for sorted dependencies in order
Update mod files up the chain, tagging new and pushing versions where applicable, cleaning go.sum
`gomu sync`

# deploy: Push local changes, then sync
Performs sync with the added functionality of committing/pushing local changes to a provided branch
`gomu deploy`

#Options
Flags and args

# -dir: Search within Target dir (accepts multiple -dir flags)
-dir will aggregate all files within provided directories. Omitting entirely will run from the current directory
`gomu -dir <your-organization> sync`

# -dep: Filter libs depending on dep (accepts multiple -dep flags)
-dep will ignore any libs wich do not the provided dep suffix in the go.sum file
`gomu -dep vroomy/plugins sync`

# -action: Interchangable with trailing command arg
-action is simply for convenience allowing training -flags instead of requiring the ending arg to be the command
`gomu -action replace-local`

# -name-only: Minimize output to goUrl of updated files only
-name-only is typically used for | command chains or simply less verbosity
`gomu -action list -name-only -dir hatchify -dir vroomy -dir <your-organization> -dep errors`

# -tag: Set Version tag
-tag will force set a specific tag upon sync or deploy. Omitting will attempt to increment tag if a vx.x.x tag is currently set for the lib

NOTE: go.mod seems to have trouble with versions greater than v1.0.0 - this is not tested or officially supported as of this version of gomu
`gomu -action sync -tag v0.1.0`

# -branch: Checkout branch 
-branch will checkout/create a provided branch before making any changes (and eventually manage pull requests)
`gomu -action pull -branch develop`
