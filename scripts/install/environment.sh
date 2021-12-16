#!/bin/bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# Set unified default password
readonly PASSWORD=${PASSWORD:-"iam59!z$"}

# Linux system user name and password
readonly LINUX_USERNAME=${LINUX_USERNAME:-"rose839"}
readonly LINUX_PASSWORD=${LINUX_PASSWORD:-${PASSWORD}}

# MariaDB configuration
readonly MARIADB_ADMIN_USERNAME=${MARIADB_ADMIN_USERNAME:-root} # MariaDB root user name
readonly MARIADB_ADMIN_PASSWORD=${MARIADB_ADMIN_PASSWORD:-${PASSWORD}} # MariaDB root user password
readonly MARIADB_HOST=${MARIADB_HOST:-127.0.0.1:3306} # MariaDB host address
readonly MARIADB_DATABASE=${MARIADB_DATABASE:-iam} # MariaDB iam database
readonly MARIADB_USERNAME=${MARIADB_USERNAME:-iam} # iam database username
readonly MARIADB_PASSWORD=${MARIADB_PASSWORD:-${PASSWORD}} # iam database password

# Redis configuration
readonly REDIS_HOST=${REDIS_HOST:-127.0.0.1} # Redis host
readonly REDIS_PORT=${REDIS_PORT:-6379} # Redis port
readonly REDIS_USERNAME=${REDIS_USERNAME:-''} # Redis user name
readonly REDIS_PASSWORD=${REDIS_PASSWORD:-${PASSWORD}} # Redis user password

# MongoDB configuration
readonly MONGO_ADMIN_USERNAME=${MONGO_ADMIN_USERNAME:-root} # MongoDB root user name
readonly MONGO_ADMIN_PASSWORD=${MONGO_ADMIN_PASSWORD:-${PASSWORD}} # MongoDB root user password
readonly MONGO_HOST=${MONGO_HOST:-127.0.0.1} # MongoDB address
readonly MONGO_PORT=${MONGO_PORT:-27017} # MongoDB port
readonly MONGO_USERNAME=${MONGO_USERNAME:-iam} # MongoDB user name
readonly MONGO_PASSWORD=${MONGO_PASSWORD:-${PASSWORD}} # MongoDB user password