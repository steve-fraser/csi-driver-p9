#!/bin/bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eo pipefail

function cleanup {
  echo 'pkill -f p9plugin'
  pkill -f p9plugin
  echo 'Deleting CSI sanity test binary'
  rm -rf csi-test
  echo 'Uninstalling P9 server on localhost'
  docker rm p9 -f
}
trap cleanup EXIT

function install_csi_sanity_bin {
  echo 'Installing CSI sanity test binary...'
  git clone https://github.com/kubernetes-csi/csi-test.git -b v2.2.0
  pushd csi-test/cmd/csi-sanity
  make
  popd
}

function provision_p9_server {
  echo 'Installing P9 server on localhost'
  apt-get update -y
  apt-get install -y nfs-common
  docker run -d --name p9 --privileged -p 2049:2049 -v $(pwd)/p9share:/p9share -e SHARED_DIRECTORY=/p9share itsthenetwork/nfs-server-alpine:latest
}

provision_p9_server
install_csi_sanity_bin

readonly endpoint='unix:///tmp/csi.sock'
nodeid='CSINode'
if [[ "$#" -gt 0 ]] && [[ -n "$1" ]]; then
  nodeid="$1"
fi

bin/p9plugin --endpoint "$endpoint" --nodeid "$nodeid" -v=5 &

echo 'Begin to run sanity test...'
readonly CSI_SANITY_BIN='csi-test/cmd/csi-sanity/csi-sanity'
"$CSI_SANITY_BIN" --ginkgo.v --csi.testvolumeparameters="$(pwd)/test/sanity/params.yaml" --csi.endpoint="$endpoint" --ginkgo.skip="should not fail when requesting to create a volume with already existing name and same capacity|should fail when requesting to create a volume with already existing name and different capacity|should work|should fail when the requested volume does not exist"
