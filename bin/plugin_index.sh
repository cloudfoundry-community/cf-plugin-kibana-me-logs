#!/bin/bash

export AUTHOR=${AUTHOR:-"Dr Nic Williams"}
export EMAIL=${EMAIL:-"drnicwilliams@gmail.com"}
export GH_AUTHOR=${GH_AUTHOR:-drnic}
export HOMEPAGE=${HOMEPAGE:-"http://github.com/$GH_AUTHOR"}
export GH_ORG=${GH_ORG:-cloudfoundry-community}
export GH_REPO=${GH_REPO:-cf-plugin-kibana-me-logs}
export NAME=${NAME:-"kibana-me-logs"}
export DESCRIPTION=${DESCRIPTION:-"Launches the Kibana UI (from [kibana-me-logs](https://github.com/cloudfoundry-community/kibana-me-logs)\) for an application."}
export PKG_DIR=${PKG_DIR:=out}
export PROJECT_CREATED="2015-03-18"

VERSION=$(<VERSION)

if [[ "$(which md5sum)X" == "X" ]]; then
  echo "md5sum not installed"
  exit 1
fi

cat << EOS
- name: $NAME
  description: $DESCRIPTION
  version: $VERSION
  created: $PROJECT_CREATED
  updated: $(date +%F)
  company: $GH_ORG
  authors:
  - name: "$AUTHOR"
    homepage: $HOMEPAGE
    contact: $EMAIL
  homepage: http://github.com/$GH_ORG/$GH_REPO
  binaries:
  - platform: win64
    url: "https://github.com/$GH_ORG/$GH_REPO/releases/download/v$VERSION/${GH_REPO}_windows_amd64.exe"
    checksum: "$(md5sum out/${GH_REPO}_windows_amd64.exe | awk '{print $1}')"
  - platform: linux64
    url: "https://github.com/$GH_ORG/$GH_REPO/releases/download/v$VERSION/${GH_REPO}_linux_amd64"
    checksum: "$(md5sum out/${GH_REPO}_linux_amd64 | awk '{print $1}')"
  - platform: darwin
    url: "https://github.com/$GH_ORG/$GH_REPO/releases/download/v$VERSION/${GH_REPO}_darwin_amd64"
    checksum: "$(md5sum out/${GH_REPO}_darwin_amd64 | awk '{print $1}')"
EOS
