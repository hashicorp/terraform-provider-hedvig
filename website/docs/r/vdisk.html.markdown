---
layout: "hedvig"
page_title: "Hedvig: hedvig_vdisk"
sidebar_current: "docs-hedvig-vdisk"
description: |-
  Manages a Vdisk resource on a Hedvig cluster.
---

# hedvig\_vdisk

Manages a Vdisk resource on a Hedvig cluster. For more information, visit [Hedvig's webpage](http://hedvig.io).

## Example Usage

Example creating a Vdisk resource.

```
resource "hedvig_vdisk" "example-vdisk" {
  name = "HedvigVdisk01"
  residence = "HDD"
  size = 20
  type = "NFS"
  blocksize = "4096"
  cacheenabled = "false"
  compressed = "true"
  deduplication = "false"
  description = "Short description of this disk."
  encryption = "false"
  replicationfactor = 3
  replicationpolicy = "Agnostic"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to be used by the Vdisk for identification.

* `residence` - (Optional) Disk residence; can be either `HDD` or `Flash`

* `size` - (Required) The size of the disk in GB

* `type` - (Required) The type of the disk; can be either `BLOCK` or `NFS`

* `scsi3pr` - (Required, defaults to false) Enables SCSI-3 Persistent Reservations for use wwith Clustered Shared Volumes (CSV)

* `blocksize` - (Optional, defaults to 4096)
 
* `cacheenabled` - (Optional, defaults to false) Enables client-side caching support for virtual disk blocks, to cache to local SSD or PCIe devices at the application compute tier

* `clusteredfilesystem` - (Optional, defaults to false) Formats a clustered file system on top of a virtual disk to be presented to multiple hosts

* `compressed` - (Optional, defaults to false)

* `deduplication` - (Optional, defaults to false)

* `description` - (Optional)

* `encryption` - (Optional, defaults to false) 

* `replicationfactor` - (Optional, defaults to 3) Can be any integer 1 - 6

* `replicationpolicy` - (Optional, defaults to Agnostic) Can be RackAware, DataCenterAware, or Agnostic (RackUnaware)
