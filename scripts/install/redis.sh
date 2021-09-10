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

    # 2.2 Delete # before `bind 127.0.0.1`
    echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' /etc/redis.conf

    # 2.3 Set password
    echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' /etc/redis.conf

    # 2.4 Turn down protected-mode
    echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' /etc/redis.conf

    # Disable firewall
    iam::common::sudo "systemctl stop firewalld.service"
    iam::common::sudo "systemctl disable firewalld.service"

    # 4. Start Redis
    iam::common::sudo "redis-server /etc/redis.conf"

    iam::redis::status || return 1
    iam::redis::info
    iam::log::info "install Redis successfully"
}

# Uninstall
function iam::redis::uninstall() {
    set +o errexit
    iam::common::sudo "killall redis-server"
    iam::common::sudo "yum -y remove redis"
    iam::common::sudo "rm -rf /var/lib/redis"
    set -o errexit
    iam::log::info "uninstall Redis successfully"
}

# Check redis status
function iam::redis::status() {
    if [[ -z "`pgrep redis-server`" ]]; then
      iam::log::error "Redis not running, maybe not installed properly"
      return 1
    fi

    redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
      iam::log::error "can not login with ${REDIS_USERNAME}, redis maybe not initialized properly"
      return 1
    }
}

if [[ "$*" =~ iam::redis:: ]];then
    eval $*
fi