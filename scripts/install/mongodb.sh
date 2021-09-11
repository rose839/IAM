#!/usr/bin/env bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

[[ -z $COMMON_SOURCED ]] && source ${IAM_ROOT}/scripts/install/common.sh

# Print info after install.
function iam::mongodb::info() {
cat << EOF
MongoDB Login: mongo mongodb://${MONGO_USERNAME}:'${MONGO_PASSWORD}'@${MONGO_HOST}:${MONGO_PORT}/iam_analytics?authSource=iam_analytics
EOF
}

# Install mongodb
function iam::mongodb::install() {
    # 1. Set MongoDB Yum source
    echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' > /etc/yum.repos.d/mongodb-org-4.4.repo
[mongodb-org-4.4]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/\$releasever/mongodb-org/4.4/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-4.4.asc
EOF"

    # 2. Install MongoDB and MongoDB client
    iam::common::sudo "yum install -y mongodb-org"

    # 3. Disable SELinux
	echo ${LINUX_PASSWORD} | sudo -S setenforce 0 || true
	echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^SELINUX=.*$/SELINUX=disabled/' /etc/selinux/config

	# 4. Turn on external network access
	echo ${LINUX_PASSWORD} | sudo -S sed -i '/bindIp/{s/127.0.0.1/0.0.0.0/}' /etc/mongod.conf
	echo ${LINUX_PASSWORD} | sudo -S sed -i '/^#security/a\security:\n  authorization: enabled' /etc/mongod.conf

    # 5. Start MongoDB
    iam::common::sudo "systemctl enable mongod"
    iam::common::sudo "systemctl start mongod"

    # 6. Create admin account
	mongo --quiet "mongodb://${MONGO_HOST}:${MONGO_PORT}" << EOF
use admin
db.createUser({user:"${MONGO_ADMIN_USERNAME}",pwd:"${MONGO_ADMIN_PASSWORD}",roles:["root"]})
db.auth("${MONGO_ADMIN_USERNAME}", "${MONGO_ADMIN_PASSWORD}")
EOF

	# 7. Create ${MONGO_USERNAME} user account
	mongo --quiet mongodb://${MONGO_ADMIN_USERNAME}:${MONGO_ADMIN_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/iam_analytics?authSource=admin << EOF
use iam_analytics
db.createUser({user:"${MONGO_USERNAME}",pwd:"${MONGO_PASSWORD}",roles:["dbOwner"]})
db.auth("${MONGO_USERNAME}", "${MONGO_PASSWORD}")
EOF

    iam::mongodb::status || return 1
    iam::mongodb::info
    iam::log::info "install MongoDB successfully"
}

# Uninstall
function iam::mongodb::uninstall() {
    set +o errexit
    iam::common::sudo "systemctl stop mongodb"
    iam::common::sudo "systemctl disable mongodb"
    iam::common::sudo "yum -y remove mongodb-org"
    iam::common::sudo "rm -rf /var/lib/mongo"
    iam::common::sudo "rm -f /etc/yum.repos.d/mongodb-10.5.repo"
    iam::common::sudo "rm -f /etc/mongod.conf"
    iam::common::sudo "rm -f /lib/systemd/system/mongod.service"
    iam::common::sudo "rm -f /tmp/mongodb-*.sock"
    set -o errexit
    iam::log::info "uninstall MongoDB successfully"
}

# Status Check
function iam::mongodb::status() {
    systemctl status mongod |grep -q 'active' || {
        iam::log::error "mongodb failed to start, maybe not installed properly"
        return 1
    }

    echo "show dbs" | mongo --quiet "mongodb://${MONGO_HOST}:${MONGO_PORT}" &>/dev/null || {
        iam::log::error "cannot connect to mongodb, mongo maybe not installed properly"
        return 1
    }

    echo "show dbs" | \
            mongo --quiet mongodb://${MONGO_ADMIN_USERNAME}:${MONGO_ADMIN_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/iam_analytics?authSource=admin &>/dev/null || {
        iam::log::error "can not login with ${MONGO_ADMIN_USERNAME}, mongo maybe not initialized properly"
        return 1
    }

    echo "show dbs" | \
            mongo --quiet mongodb://${MONGO_USERNAME}:${MONGO_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/iam_analytics?authSource=iam_analytics &>/dev/null|| {
        iam::log::error "can not login with ${MONGO_USERNAME}, mongo maybe not initialized properly"
        return 1
    }
}

if [[ "$*" =~ iam::mongodb:: ]];then
  eval $*
fi
