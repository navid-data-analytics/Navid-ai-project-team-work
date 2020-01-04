#!/bin/bash

set -e

USAGE="
Usage: $(basename $0) tagname
the tagname shall starts with 'manual' to indicate it is not for automated build
"

github_repo="git@github.com:callstats-io/ai-decision.git"
dockerhub_repo_name="ai-decision"

tag=$1

if [ -z "$tag" ]; then
    echo "tag not set"
    echo "${USAGE}"
    exit 1
fi

# the tag needs to starts with 'manual' prefix
if [[ ! $tag =~ ^manual ]]; then
    echo "tag needs to starts with 'manual'"
    exit 1
fi

# clone the github repo to a clean directory to make the build
tempdir="manual-builds-${tag}"

#if the tempdir already exists, prompt the user to check what is going on
if [ -d "$tempdir" ]; then
    echo "directory $tempdir already exists, please delete it and try again."
    exit 1
fi

git clone --branch $tag  $github_repo $tempdir

# make sure we delete the temporary stuff created by the script
trap "rm -rf $tempdir; docker rmi callstats/${dockerhub_repo_name}:${tag}" EXIT

echo "**********************"
echo "* build starts: $(date) "
echo "**********************"

#build and publish images
docker build -t callstats/${dockerhub_repo_name}:${tag} $tempdir
docker push callstats/${dockerhub_repo_name}:${tag}
docker rmi callstats/${dockerhub_repo_name}:${tag}

echo "**********************"
echo "* build done: $(date) "
echo "**********************"