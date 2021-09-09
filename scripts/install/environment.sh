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