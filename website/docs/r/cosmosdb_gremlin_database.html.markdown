---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_gremlin_database"
description: |-
  Manages a Gremlin Database within a Cosmos DB Account.
---

# azurerm_cosmosdb_gremlin_database

Manages a Gremlin Database within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_gremlin_database" "example" {
  name                = "tfex-cosmos-gremlin-db"
  resource_group_name = data.azurerm_cosmosdb_account.example.resource_group_name
  account_name        = data.azurerm_cosmosdb_account.example.name
  throughput          = 400
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Gremlin Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Gremlin Database is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the CosmosDB Account to create the Gremlin Database within. Changing this forces a new resource to be created.

* `throughput` - (Optional) The throughput of the Gremlin database (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.

~> **Note:** throughput has a maximum value of `1000000` unless a higher limit is requested via Azure Support

* `autoscale_settings` - (Optional) An `autoscale_settings` block as defined below. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.

~> **Note:** Switching between autoscale and manual throughput is not supported via Terraform and must be completed via the Azure Portal and refreshed.

---

An `autoscale_settings` block supports the following:

* `max_throughput` - (Optional) The maximum throughput of the Gremlin database (RU/s). Must be between `1,000` and `1,000,000`. Must be set in increments of `1,000`. Conflicts with `throughput`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CosmosDB Gremlin Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB Gremlin Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Gremlin Database.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB Gremlin Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB Gremlin Database.

## Import

CosmosDB Gremlin Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_gremlin_database.db1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/gremlinDatabases/db1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DocumentDB` - 2024-08-15
