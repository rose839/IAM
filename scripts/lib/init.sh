#!/usr/bin/env bash

set -o errexit
set +o nounset
set -o pipefail

# Default use go modules
export GO111MODULE=on

IAM_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/../.. && pwd -P)

source "${IAM_ROOT}/scripts/lib/logging.sh"