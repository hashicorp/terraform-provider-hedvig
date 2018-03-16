---
layout: "hedvig"
page_title: "Hedvig: hedvig_mount"
sidebar_current: "docs-hedvig-mount"
description: |-
  Mounts vdisk resource with a particular controller.
---

# hedvig\_mount

A Hedvig Mount mounts a vdisk resource with a particular controller. It can then be used to connect ACL access resources to the vdisk as well.

## Example Usage

```
resource "hedvig_mount" "example-mount" {
  cluster = "example"
  vdisk = "${hedvig_vdisk.example-vdisk.name}"
  controller = "examplevip1.hedviginc.com"
}
```

## Argument Reference

The following arguments are supported:

* `cluster` - (Required) The name of the cluster hosting the Mount.

* `vdisk` - (Required) The name of the vdisk the Mount is on.

* `controller` - (Required) The fully qualified domain name for the controller that the Mount is to attach to.
