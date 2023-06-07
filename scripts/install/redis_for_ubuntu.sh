#!/bin/bash

# The root of the build/dist directory.
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# Print necessary information after installation.
function iam::redis::info() {
    cat << EOF
Redis Login: redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a '${REDIS_PASSWORD}'
EOF
}

# Install Redis.
function iam::redis::install() {
    # 1. Install Redis.
    iam::common::sudo "apt-get install -y redis"

    # 2. Config.
    # 2.1 Set daemon.
    echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' /etc/redis/redis.conf

    # 2.2 Add '#' before 'bind 127.0.0.1' to comment it out. By default, only local connections are allowed. After commenting it out, the external network can connect to Redis.
    echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' /etc/redis/redis.conf

    # 2.3 Set password.
    echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' /etc/redis/redis.conf

    # 2.4 Turn down protected-mode.
    echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' /etc/redis/redis.conf

    # 3. Shutdown firewall.
    iam::common::sudo "sudo ufw disable"
    iam::common::sudo "sudo ufw status"

    # 4. Start Redis.
    iam::common::sudo "redis-server /etc/redis/redis.conf"

    iam::redis::status || return 1
    iam::redis::info
    iam::log::info "install Redis successfully"
}

# Uninstall Redis.
function iam::redis::uninstall() {
    set +o errexit
    iam::common::sudo "/etc/init.d/redis-server stop"
    iam::common::sudo "apt-get -y remove redis-server"
    iam::common::sudo "rm -rf /var/lib/redis"
    set -o errexit
    iam::log::info "uninstall Redis successfully"
}

# Status check.
function iam::redis::status() {
    if [[ -z "`pgrep redis-server`" ]];then
        iam::log::error_exit "Redis not running, maybe not installed properly"
        return 1
    fi


    redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
        iam::log::error "can not login with ${REDIS_USERNAME}, redis maybe not initialized properly"
        return 1
    }

    iam::log::info "redis-server status active"
}

if [[ "$*" =~ iam::redis:: ]];then
  eval $*
fi