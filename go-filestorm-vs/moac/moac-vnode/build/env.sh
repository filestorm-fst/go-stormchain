#!/bin/sh
#modified to fit GOOGLE CLOUD
set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace under build if it doesn't exist yet.
workspace="$PWD/build/_workspace"
# Record the root of the repository
root="$PWD"
moacdir="$workspace/src/github.com/filestorm/go-filestorm/moac"
if [ ! -L "$moacdir/moac-vnode" ]; then
    mkdir -p "$moacdir"
    cd "$moacdir"
    echo "Make" $moacdir
    ln -s ../../../../../. moac-vnode
    #Add a library path
    ln -s ../../../../../../moac-lib moac-lib
    cd "$root"
fi
# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH
echo Set GOPATH $GOPATH
# Run the command inside the workspace.
cd "$moacdir/moac-vnode"
PWD="$moacdir/moac-vnode"
# Launch the arguments with the configured environment.
exec "$@"
