#!/bin/bash

# Install svu for version comparison
go install github.com/caarlos0/svu@v1.9.0

# Fetch the latest tags
git fetch --tags

CURRENT=$(svu current)
NEXT=$(svu next)

if [ $CURRENT != $NEXT ]
then
  echo "Tagging with" $NEXT
  git tag $NEXT
  git tag latest
  git push --tags

  # Do the release
  go install github.com/goreleaser/goreleaser@v1.4.1
  goreleaser --rm-dist
else
  echo "No new version detected. Skipping release."
fi
