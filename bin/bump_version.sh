#!/bin/bash

new_version=$1

if [[ "${new_version}X" == "X" ]]; then
  echo "USAGE ./bin/bump_version X.Y.Z"
  exit 1
fi

if [[ "$(which go-bindata)X" == "X" ]]; then
  echo "Installing go-bindata..."
  go get -u github.com/jteeuwen/go-bindata/...
fi

echo $new_version > VERSION
go-bindata -o version.go VERSION
