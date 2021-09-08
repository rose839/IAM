#!/bin/bash

# Common utilities, variables and checks for all scripts.

# If command's return value is not 0, then quit the execution of the scripts.
set -o errexit

# Turn down the utility which give an error when encounter an undefined variable.
set +o nounset

# Return non-zero value when any command in pipeline fails.
set -o pipefail

# Sourced flag(whether this script has been sourced).
COMMON_SOURCED=true

# The root dir of the project.
IAM_ROOT=$(dirname ${BASH_SOURCE[0]})/../..

# Redefine sudo, not need password.
function iam::common::sudo() {
    echo $LINUX_PASSWORD | sudo -S $1
}