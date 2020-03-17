#GoMu - Go Mod Utils
GoMu is intended to make working with go.mod sane.

#1) Sort
List all dependencies in order of the dependency chain, optionally warning of dependency cycles
`gomu sort`


#2) Sync
Update mod files up the chain, tagging new and pushing versions where applicable, cleaning go.sum, optionally warning of multiple dependencies
`gomu sync`

#3) Deploy
Performs sync with the added functionality of pushing local changes to a provided branch
`gomu deploy`

#Options
Flags and args
#1) Filter
-filter will ignore any libs wich do not include a provided string in the go.sum file
`gomu sort -filter vpm -filter github.com/vroomy/plugins`

#2) Target
-target will aggregate all filles withing provided directories. Omitting entirely will run from the current directory
`gomu sync -target github.com/vroomy -target github.com/<your-organization>`

#3) Tag
-tag will force set a specific tag upon sync or deploy. Omitting will attempt to increment tag if a vx.x.x tag is currently set for the lib
NOTE: go.mod seems to have trouble with versions greater than v1.0.0 - this is not tested or officially supported as of this version of gomu
`gomu sync -tag v0.1.0`
