Rustack Terraform Provider
==================

<!-- TODO: Update this link -->
<!-- - Documentation: https://registry.terraform.io/providers/rustack/rustack/latest/docs -->

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.0.10
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

Building The Provider
---------------------
<!-- TODO -->

Using the provider
----------------------

See the [Rustack Provider documentation](https://registry.terraform.io/providers/rustack/rustack/latest/docs) to get started using the Rustack provider.

Developing the Provider
---------------------------

To compile the provider, run `make build`.

For information about writing acceptance tests, see the main Terraform [contributing guide](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#writing-acceptance-tests).

Releasing the Provider
----------------------

A [Gorelaser](https://goreleaser.com/) configuration is provided that produces
build artifacts matching the [layout required](https://www.terraform.io/docs/registry/providers/publishing.html#manually-preparing-a-release)
to publish the provider in the Terraform Registry.

Releases will appear as drafts. Once marked as published on the GitHub Releases page,
they will become available via the Terraform Registry.
