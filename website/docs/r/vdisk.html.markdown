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
  size = 20
  type = "NFS"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to be used by the Vdisk for identification.

* `size` - (Required) The size of the disk in GB.

* `type` - (Required) The type of the disk; can be either `BLOCK` or `NFS`
