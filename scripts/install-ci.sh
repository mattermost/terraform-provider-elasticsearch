#!/usr/bin/env bash

# Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
# See LICENSE.txt for license information.

set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

GOBIN=$(go env GOPATH)/bin

if [ -x $GOBIN/golangci-lint ]; then
    echo "Golang-CI linter is already installed"
    exit 0
fi

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $GOBIN v1.39.0
echo "Golang-CI linter installed succesfully"
