#!/usr/bin/env bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}"/../..)
source "${IAM_ROOT}/scripts/install/common.sh"

function iam::install::install_cfssl() {
    mkdir $HOME/bin/
    wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -O $HOME/bin/cfssl
    wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -O $HOME/bin/cfssljson
    wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64 -O $HOME/bin/cfssl-certinfo
    chomd +x $HOME/bin/{cfssl, cfssljson, cfssl-certinfo}
    iam::log::info "install cfssl tools successfully"
}

if [[ $* =~ "iam::install::" ]]; then
    eval $*
fi