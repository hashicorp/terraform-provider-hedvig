provider "hedvig" {
  username = "HedvigAdmin"
  password = "hedvig"
  node = "tfhashicorp1.external.hedviginc.com"
}

resource "hedvig_vdisk" "my-vdisk" {
  name = "HedvigVdiskNN"
  residence = "HDD"
  size = 18
  type = "BLOCK"
}

resource "hedvig_lun" "my-lun3" {
  vdisk = "${hedvig_vdisk.my-vdisk.name}"
  controller = "tfhashicorpvip1.external.hedviginc.com"
}

resource "hedvig_vdisk" "my-vdisk2" {
  name = "HedvigVdiskOO"
  size = 20
  type = "NFS"
}

#resource "hedvig_vdisk" "my-vdisk3" {
#  name = "HedvigVdiskC"
#  size = 22
#  type = "NFS"
#}

resource "hedvig_mount" "my-mount" {
  vdisk = "${hedvig_vdisk.my-vdisk2.name}"
  controller = "tfhashicorpvip1.external.hedviginc.com"
}

#resource "hedvig_access" "my-access" {
#  vdisk = "${hedvig_vdisk.my-vdisk.name}"
#  host = "${hedvig_lun.my-lun3.controller}"
#  address = "172.22.22.21"
#  type = "host"
#}

#resource "hedvig_access" "my-access2" {
#  vdisk = "${hedvig_vdisk.my-vdisk.name}"
#  host = "${hedvig_lun.my-lun3.controller}"
#  address = "172.22.22.25"
#  type = "host"
#}

#resource "hedvig_access" "my-access3" {
#  vdisk = "${hedvig_vdisk.my-vdisk.name}"
#  host = "${hedvig_lun.my-lun3.controller}"
#  address = "172.22.22.30"
#  type = "host"
#}

#resource "hedvig_access" "my-access4" {
#  vdisk = "${hedvig_vdisk.my-vdisk2.name}"
#  host = "${hedvig_mount.my-mount.controller}"
#  address = "172.22.22.31"
#  type = "host"
#}

