#!/bin/bash

function get_master {
    echo "$FILE: Pulling master..."
    # update working copy
    git checkout master > /dev/null 2>&1
    git pull > /dev/null 2>&1
    git pull --tag > /dev/null 2>&1
}

function update_mod_files {
    rm go.mod > /dev/null 2>&1
    rm go.sum > /dev/null 2>&1
    go mod init > /dev/null 2>&1 && go mod tidy > /dev/null 2>&1

    # check for local changes to mod files
    git add go.mod > /dev/null 2>&1 
    git add go.sum > /dev/null 2>&1 
    git commit -m "update mod files" > /dev/null 2>&1 && echo "$FILE: Updating mod files..." && git push > /dev/null 2>&1 || echo "$FILE: Deps up to date!"
}

function update_tag {
    # get current tag
    TAG=$(git-tagger --action=get)
    TAGCOMMIT=$(git rev-list -n 1 "$TAG")
    COMMIT=$(git rev-parse HEAD)

    if [ "$TAGCOMMIT" == "$COMMIT" ]; then
        echo "$FILE: Tag up to date!"
        return
    fi

    # tag latest commit
    echo "$FILE: Setting new tag..."
    git-tagger > /dev/null 2>&1
    TAG=$(git-tagger --action=get)
    echo "$FILE: Updated tag to $TAG!"
}

function main {
    echo "Cleaning mod cache..."
    go clean --modcache

    # save base dir
    DIR=`pwd`

    # iterate over all files piped to input
    index=0
    while read FILE
    do
        # increment index
        ((index=index+1))

        echo ""
        echo "$index) Scanning $FILE"

        # attempt to cd into piped path
        cd $DIR
        cd $FILE || continue

        get_master
        continue

        # ignore dirs without go.mod
        if [ -f "go.mod" ]; then
            echo "$FILE: Checking deps..."
        else
            echo "$FILE: go.mod not found. Skipping"
            continue
        fi

        update_mod_files

        if [[ $FILE =~ "-plugin" ]]; then
            echo "$FILE: Not setting tag for plugins"
            continue
        fi

        update_tag
    done
}

### SCRIPT START ###
main
