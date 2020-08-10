provider "hedvig" {
  username = "admin"
  password = "password"
  node = "node.domainname.com"
}

resource "hedvig_vdisk" "my-vdisk-lumosBlock55" {
  name = "HedvigVdiskLumosBlock55"
#  clusteredfilesystem = "true"
  deduplication = "true"
  scsi3pr = "false"
  compressed = "true"
  encryption = "false"
#  description = "Stuff about this Vdisk."
  residence = "HDD"
  type = "Block"
  size = 7
  blocksize = "4096"
  replicationpolicy = "Agnostic"
  cacheenabled = "true"
}

#resource "hedvig_access" "my-access-fudgeFlash23" {
#  vdisk = "${hedvig_vdisk.my-vdisk-fudgeFlash23.name}"
#  host = "${hedvig_mount.my-mount-fudgeFlash23.controller}"
#  address = "172.22.22.8"
#  type = "host"
#}

#resource "hedvig_mount" "my-mount-fudgeNFS24" {
#  vdisk = "${hedvig_vdisk.my-vdisk-fudgeNFS24.name}"
#  controller = "hedvigucs3.r3.snc1.hedviginc.com.hedviginc.com"
#}

#resource "hedvig_lun" "my-lun-lumosHDD21" {
#  vdisk = "${hedvig_vdisk.my-vdisk-lumosHDD21.name}"
#  controller = "lumosvip3.hedviginc.com"
#}

