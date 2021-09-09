#!/usr/bin/env bash

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# Print info after install.
function iam::redis::info() {
cat << EOF
Redis Login: redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a '${REDIS_PASSWORD}'
EOF
}

# Install
function iam::redis::install() {
    # 1. Install redis
    iam::common::sudo "yum install -y redis"

    # 2. Config
    # 2.1 Set daemon
    iam::common::sudo "sed -i '/^daemonize/{s/no/yes/}' /etc/redis.conf"
}

if [[ "$*" =~ iam::redis:: ]];then
  eval $*
fi