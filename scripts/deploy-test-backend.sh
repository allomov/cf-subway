#!/bin/bash

if [[ "${BOSH_TARGET}X" == "X" ]]; then
  echo "Required \$BOSH_TARGET, \$BOSH_USERNAME, \$BOSH_PASSWORD"
  exit 1
fi

cat > $HOME/.bosh_config << EOS
---
auth:
  ${BOSH_TARGET}:
    username: ${BOSH_USERNAME}
    password: ${BOSH_PASSWORD}
EOS

bosh target $BOSH_TARGET

if [[ "$(bosh releases | grep ' docker ')X" == "X" ]]; then
  bosh upload release https://bosh.io/d/github.com/cloudfoundry-community/docker-boshrelease
fi
if [[ "$(bosh releases | grep ' postgresql-docker ')X" == "X" ]]; then
  bosh upload release https://bosh.io/d/github.com/cloudfoundry-community/postgresql-docker-boshrelease
fi

cd /tmp
git clone https://github.com/cloudfoundry-community/postgresql-docker-boshrelease.git postgresql-docker
cd postgresql-docker

./templates/make_manifest warden broker
bosh -n deploy
