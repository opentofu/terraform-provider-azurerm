# AzureRM Provider Contributor Guides

**First,** thank you for your interest in contributing to the Azure Provider! And if you're unsure or anything, please do reach out for help. You can open a draft pull request (PR) or an issue with what you know or join the [Slack Workspace for Contributors](https://terraform-azure.slack.com) ([Request Invite](https://join.slack.com/t/terraform-azure/shared_invite/zt-37y31pz9c-chP0FWRtmsagfpF8r4m37g)) and we'll do our best to guide you in the right direction.

> **Note:** this documentation is a work-in-progress - if you see something that's not quite right or missing, we'd really appreciate a PR!

This contribution guide assumes you have at least a basic understanding of both Go and Terraform itself (for example you know what a Data Source and a Resource are) - more information on those can be found [in the Terraform documentation](https://www.terraform.io/docs/language/index.html).

---

The AzureRM Provider is a Plugin which is invoked by Terraform (Core) and comprised of Data Sources and Resources.

Within the AzureRM Provider, these Data Sources and Resources are grouped into Service Packages - which are logical groupings of Data Sources/Resources based on the Azure Service they're related to.

Each of these Data Sources and Resources has both Acceptance Tests and Documentation associated with each Data Source/Resource - the Acceptance Tests are also located within this Service Package, however the Documentation exists within a dedicated folder.

More granular documentation covers how these fit together - and the most common types of contribution we see:

## Topics

Basics:

* [Overview of the Provider](topics/high-level-overview.md)
* [Building the Provider](topics/building-the-provider.md)
* [Debugging the Provider](topics/debugging-the-provider.md)
* [Running the Tests](topics/running-the-tests.md)
* [Opening a Pull Request](topics/guide-opening-a-pr.md)

Common Topics/Guides:

* [Adding a new Feature](topics/guide-new-feature.md)
* [Adding a new Service Package](topics/guide-new-service-package.md)
* [Adding a new Data Source](topics/guide-new-data-source.md)
* [Adding a new Resource](topics/guide-new-resource.md)
* [Adding fields to an existing Data Source](topics/guide-new-fields-to-data-source.md)
* [Adding fields to an existing Resource](topics/guide-new-fields-to-resource.md)
* [Adding State Migrations](topics/guide-state-migrations.md)
* [Adding Write-Only Attributes](topics/guide-new-write-only-attribute.md)
* [Breaking Changes and Deprecations](topics/guide-breaking-changes.md)
* [When to create a new Resource vs Inline Block](topics/guide-new-resource-vs-inline.md)

References:

* [Acceptance Testing](topics/reference-acceptance-testing.md)
* [Best Practices](topics/best-practices.md)
* [Glossary](topics/reference-glossary.md)
* [Naming](topics/reference-naming.md)
* [Provider Documentation Standards](topics/reference-documentation-standards.md)
* [Resource IDs](topics/guide-resource-ids.md)
* [Schema Design](topics/schema-design-considerations.md)
* [Working with Errors](topics/reference-errors.md)


Maintainer specific:

* [Merging Prs](topics/maintainer-merging.md)

FAQ:

* [Frequently Asked Questions](topics/frequently-asked-questions.md)
