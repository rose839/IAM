#!/usr/bin/env bash

# Control the verbosity of script output and logging
IAM_VERBOSE=${IAM_VERBOSE:-5}

function iam::log::info() {
    local V=${V:-0}
    if ((${IAM_VERBOSE} < $V)); then
        return
    fi

    for message; do
        echo $message
    done
}

# Log error to stderr.
function iam::log::error() {
    local timestamp=$(date +"[%m%d %H:%M:%S]")
    echo "!!! ${timestamp} ${1-}" >&2
    shift
    for message; do
        echo "    ${message}" >&2
    done
}