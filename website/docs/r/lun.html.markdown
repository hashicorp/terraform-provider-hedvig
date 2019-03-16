---
layout: "hedvig"
page_title: "Hedvig: hedvig_lun"
sidebar_current: "docs-hedvig-lun"
description: |-
  Adds a virtual disk as a LUN to a controller.
---

# hedvig\_lun

A Hedvig Lun adds a vdisk resource to a particular controller, enabling the addition of ACL access resources to the vdisk.

## Example Usage

Example creating a Lun resource.

```
resource "hedvig_lun" "example-lun" {
  vdisk = "${hedvig_vdisk.example-vdisk.name}"
  controller = "examplevip1.hedviginc.com"
}
```

## Argument Reference

The following arguments are supported:

 * `vdisk` - (Required) The name of the vdisk the LUN is on.

 * `controller` - (Required) The fully qualified domain name for the controller that the LUN is to attach to.
