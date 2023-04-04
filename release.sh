#!/bin/bash

# If GITHUB_TOKEN is not defined, exit
if [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_TOKEN is not defined"
  exit 1
fi

CURRENT_VERSION=`git tag --list | sort -V | tail -1`
CURRENT_MAJOR=`echo $CURRENT_VERSION | cut -f1 -d'.' | tr -d 'v'`
CURRENT_MINOR=`echo $CURRENT_VERSION | cut -f2 -d'.'`
CURRENT_PATCH=`echo $CURRENT_VERSION | cut -f3 -d'.'`

# Switch based on the first argument value
case $1 in
  "upgrade")
    NEW_MAJOR=$((CURRENT_MAJOR+1))
    NEW_MINOR=0
    NEW_PATCH=0
    ;;
  "update")
    NEW_MAJOR=$CURRENT_MAJOR
    NEW_MINOR=$((CURRENT_MINOR+1))
    NEW_PATCH=0
    ;;
  "patch")
    NEW_MAJOR=$CURRENT_MAJOR
    NEW_MINOR=$CURRENT_MINOR
    NEW_PATCH=$((CURRENT_PATCH+1))
    ;;
  *)
    echo "Usage: $0 {upgrade|update|patch}"
    exit 1
    ;;
esac

export RELEASE="v$NEW_MAJOR.$NEW_MINOR.$NEW_PATCH"
echo "Releasing $RELEASE..."
git tag -a $RELEASE -m "Release $RELEASE"
git push origin $RELEASE
goreleaser release --rm-dist
