#!/bin/bash

new_version=$1

if [[ "${new_version}X" == "X" ]]; then
  echo "USAGE ./bin/bump_version X.Y.Z"
  exit 1
fi

echo $new_version > VERSION
cat >version.go <<EOL
package main

// VERSION of plugin
var VERSION = "${new_version}"
EOL
