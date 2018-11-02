Terraform Provider for Hedvig
=============================

- Website: http://www.hedvig.io
- Documentation: https://www.terraform.io/docs/providers/hedvig/index.html
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg width="600pxx">

Maintainers
-----------

This provider plugin is maintained by:

* The [Hedvig Infrastructure Team](https://www.hedvig.io/blog/hedvig-terraform-simplified-apis-for-muliple-providers)
# The Terraform team at [HashiCorp](https://www.hashicorp.com)

Requirements
------------

-      [Terraform](https://www.terraform.iodownloads.html) 0.10+
-      [Go](https://golang.org/doc/install) 1.11.0 or higher

Building the Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-hedvig`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPAth/src/github.com/terraform/providers
$ git clone git@github.com:terraform-providers/terraform-provider-hedvig
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-hedvig
$ make build
```

Using the provider
-----------------

See the [Hedvig Provider documentation](https://www.terraform.io/docs/providers/hedvig/index.html) to get started using the Hedvig provider

Upgrading the provider
----------------------

To upgrade to the latest stable version of the Hedvig provider run `terraform init -upgrade`. See the [Terraform website](https://www.terraform.io/docs/configuration/providers.html#provider-versions) for more information.

Developing the Provider
----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-hedvig
...
```

For guidance on common development practices such as testing changes or vendoring libraries, see the [contribution guidelines](https://github.com/terraform-providers/terraform-provider-google/blob/master/.github/CONTRIBUTING.md). If you have other development questions we don't cover, please file an issue!
