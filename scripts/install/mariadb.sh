#!/bin/bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

[[ -z $COMMON_SOURCED ]] && source ${IAM_ROOT}/scripts/install/common.sh

# Print info after install.
function iam::mariadb::info() {
cat << EOF
MariaDB Login: mysql -h127.0.0.1 -u${MARIADB_ADMIN_USERNAME} -p'${MARIADB_ADMIN_PASSWORD}'
EOF
}

# Install
function iam::mariadb::install() {
    
  # 1. Config MariaDB 10.2 Yum source.
  echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' > /etc/yum.repos.d/mariadb-10.2.repo
# MariaDB 10.2 CentOS repository list - created 2020-10-23 01:54 UTC
# http://downloads.mariadb.org/mariadb/repositories/
[mariadb]
name = MariaDB
baseurl = https://mirrors.ustc.edu.cn/mariadb/yum/10.2/centos74-amd64
module_hotfixes=1
gpgkey=https://mirrors.ustc.edu.cn/mariadb/yum/RPM-GPG-KEY-MariaDB
gpgcheck=1
EOF"

    # 2. Install mariadb server and client.
    iam::common::sudo "yum install -y MariaDB-server MariaDB-client"

    # 3. Start mariadb and set startup when boot.
    iam::common::sudo "systemctl enable mariadb"
    iam::common::sudo "systemctl start mariadb"

    # 4. Set initial password of root.
    iam::common::sudo "mysqladmin -u${MARIADB_ADMIN_USERNAME} password ${MARIADB_ADMIN_PASSWORD}"

    # 5. Check status.
    iam::mariadb::status || return 1
    iam::mariadb::info
    # iam::log::info "install MariaDB successfully"
}

# Uninstall
function iam::mariadb::uninstall() {
    iam::common::sudo "systemctl stop mariadb"
    iam::common::sudo "systemctl disable mariadb"
    iam::common::sudo "yum remove -y MariaDB-server MariaDB-client"
    iam::common::sudo "rm -rf /var/lib/mysql"
    iam::common::sudo "rm -f /etc/yum.repos.d/mariadb-10.2.repo"
    iam::log::info "uninstall MariaDB successfully"
}

# Status check
function iam::mariadb::status() {
    # Check mariadb run status
    systemctl status mariadb | grep -q "active" || {
        # iam::log::error "mariadb failed to start, maybe not installed properly"
        return 1
    }

    mysql -u${MARIADB_ADMIN_USERNAME} -p${MARIADB_ADMIN_PASSWORD} -e quit &>/dev/null || {
        # iam::log::error "can not login with root, mariadb maybe not initialized properly"
        return 1
    }
}

if [[ "$*" =~ iam::mariadb:: ]]; then
    eval $*
fi