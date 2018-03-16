# Copyright 2015 Container Solutions
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

provider "hedvig" {
  username = "HedvigAdmin"
  password = "hedvig"
  node = "tfhashicorp1.external.hedviginc.com"
}

resource "hedvig_vdisk" "my-vdisk" {
  cluster = "tfhashicorp"
  name = "HedvigVdiskNN"
  size = 18
  type = "BLOCK"
}

resource "hedvig_lun" "my-lun3" {
  cluster = "tfhashicorp3"
  vdisk = "${hedvig_vdisk.my-vdisk.name}"
  controller = "tfhashicorpvip1.external.hedviginc.com"
}

resource "hedvig_vdisk" "my-vdisk2" {
  cluster = "tfhashicorp"
  name = "HedvigVdiskOO"
  size = 20
  type = "NFS"
}

#resource "hedvig_vdisk" "my-vdisk3" {
#  cluster = "tfhashicorp"
#  name = "HedvigVdiskC"
#  size = 22
#  type = "NFS"
#}

resource "hedvig_mount" "my-mount" {
  cluster = "tfhashicorp3"
  vdisk = "${hedvig_vdisk.my-vdisk2.name}"
  controller = "tfhashicorpvip1.external.hedviginc.com"
}

#resource "hedvig_access" "my-access" {
#  cluster = "tfhashicorp"
#  vdisk = "${hedvig_vdisk.my-vdisk.name}"
#  host = "${hedvig_lun.my-lun3.controller}"
#  address = "172.22.22.21"
#  type = "host"
#}

#resource "hedvig_access" "my-access2" {
#  cluster = "tfhashicorp"
#  vdisk = "${hedvig_vdisk.my-vdisk.name}"
#  host = "${hedvig_lun.my-lun3.controller}"
#  address = "172.22.22.25"
#  type = "host"
#}

#resource "hedvig_access" "my-access3" {
#  cluster = "tfhashicorp"
#  vdisk = "${hedvig_vdisk.my-vdisk.name}"
#  host = "${hedvig_lun.my-lun3.controller}"
#  address = "172.22.22.30"
#  type = "host"
#}

#resource "hedvig_access" "my-access4" {
#  cluster = "tfhashicorp"
#  vdisk = "${hedvig_vdisk.my-vdisk2.name}"
#  host = "${hedvig_mount.my-mount.controller}"
#  address = "172.22.22.31"
#  type = "host"
#}

