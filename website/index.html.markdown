---
layout: "hedvig"
page_title: "Provider: Hedvig"
sidebar_current: "docs-hedvig-index"
description: |-
  The Hedvig provider is used to interact with Hedvig resources. The provider needs to be configured with proper credentials and a working cluster before it can be used.
---

# Hedvig Provider

The Hedvig provider is used to interact with [Hedvig](http://hedvig.io). The provider needs to be configured with the proper credentials on a working Hedvig cluster before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```
// Configure the Hedvig provider
provider "hedvig" {
  username = "Example"
  password = "example"
  node = "example1.hedviginc.com"
}

// Create a new VDisk
resource "hedvig_vdisk" "example" {
  # ...
}
```

## Configuration Reference
