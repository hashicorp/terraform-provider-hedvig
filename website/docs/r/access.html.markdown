---
layout: "hedvig"
page_title: "Hedvig: hedvig_access"
sidebar_current: "docs-hedvig-access"
description: |-
  Adds an ACL address to a controller resource of a vdisk resource.
---

# hedvig\_access

A Hedvig Access adds an address to an ACL of a controller of a vdisk. This allows for management of access resources.

## Example Usage

Example creating an Access resource.

```
resource "hedvig_access" "example-access" {
  vdisk = "${hedvig_vdisk.example-vdisk.name}"
  host = "${hedvig_lun.example-lun.controller}"
  address = "172.26.53.99"
  type = "host"
}
```

## Argument Reference

The following arguments are supported:

* `vdisk` - (Required) The name of the Vdisk that this Access is associated with.

* `host` - (Required) The fully qualified domain name of the controller this Access is associated with.

* `address` - (Required) The actual address that this Access is providing access to.

* `type` - (Required) The type of address provided in `address`. Can be `host`, `ip` or `iqn`.
