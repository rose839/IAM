#!/usr/bin/env bash

# script feature: Generate IAM component yaml config files based on scripts/install/environment.sh.
# example: genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml

env_file="$1"
template_file="$2"

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${IAM_ROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ]; then
    iam::log::error "Usage: genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml"
    exit 1
fi

source "${env_file}"

set +u

for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file}); do
    if [ -z "$(eval echo \$${env})" ]; then
        iam::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ -z $missing ]; then
    iam::log::error 'You may run `source scripts/environment.sh` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"