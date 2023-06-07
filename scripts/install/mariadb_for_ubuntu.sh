#!/bin/bash

# The root of the build/dist directory.
IAM_ROOT=$(dirname "$BASH_SOURCE[0]")/../..

[[ -z ${COMMON_SOURCED} ]] && source "${IAM_ROOT}"/scripts/install/common.sh

# Print necessary information after installation.
function iam::mariadb::info() {
    cat << EOF
MariaDB Login: mysql -h127.0.0.1 -u${MARIADB_ADMIN_USERNAME} -p'${MARIADB_ADMIN_PASSWORD}'
EOF
}

# Install MariaDB.
function iam::mariadb::install() {
    # 1. Config MariaDB 10.5 apt source.
    iam::common::sudo "apt-get install -y software-properties-common dirmngr apt-transport-https"
    echo ${LINUX_PASSWORD} | sudo -S apt-key adv --fetch-keys 'https://mariadb.org/mariadb_release_signing_key.asc'
    echo ${LINUX_PASSWORD} | sudo -S add-apt-repository 'deb [arch=amd64,arm64,ppc64el,s390x] https://mirrors.aliyun.com/mariadb/repo/10.5/ubuntu focal main'

    # 2. Install mariadb server and client.
    iam::common::sudo "apt update"
    iam::common::sudo "apt -y install mariadb-server"

    # 3. Start mariadb and set startup when boot.
    iam::common::sudo "systemctl enable mariadb"
    iam::common::sudo "systemctl start mariadb"

    # 4. Set initial password of root.
    iam::common::sudo "mysqladmin -u${MARIADB_ADMIN_USERNAME} password ${MARIADB_ADMIN_PASSWORD}"

    # 5. Check status.
    iam::mariadb::status || return 1
    iam::mariadb::info
    iam::log::info "install MariaDB successfully"
}

# Uninstall MariaDB.
function iam::mariadb::uninstall() {
    set +o errexit
    iam::common::sudo "systemctl stop mariadb"
    iam::common::sudo "systemctl disable mariadb"
    iam::common::sudo "apt-get -y remove mariadb-server"
    iam::common::sudo "rm -rf /var/lib/mysql"
    iam::log::info "uninstall MariaDB successfully"
    set -o errexit
}

# Check MariaDB status.
function iam::mariadb::status() {
    # Check mariadb run status
    systemctl status mariadb | grep -q "active" || {
        iam::log::error "mariadb failed to start, maybe not installed properly"
        return 1
    }

    mysql -u${MARIADB_ADMIN_USERNAME} -p${MARIADB_ADMIN_PASSWORD} -e quit &> /dev/null || {
        iam::log::error "mariadb failed to login, maybe not installed properly"
        return 1
    }

    iam::log::info "MariaDB status active"
}

if [[ "$*" =~ iam::mariadb:: ]]; then
    eval "$*"
fi