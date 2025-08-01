---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_host_pool_registration_info"
description: |-
  Manages a Virtual Desktop Host Pool Registration Info.
---

# azurerm_virtual_desktop_host_pool_registration_info

Manages the Registration Info for a Virtual Desktop Host Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-hostpool"
  location = "westeurope"
}

resource "azurerm_virtual_desktop_host_pool" "example" {
  name                 = "example-HP"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  type                 = "Pooled"
  validate_environment = true
  load_balancer_type   = "BreadthFirst"

}

resource "azurerm_virtual_desktop_host_pool_registration_info" "example" {
  hostpool_id     = azurerm_virtual_desktop_host_pool.example.id
  expiration_date = "2022-01-01T23:40:52Z"
}
```

## Arguments Reference

The following arguments are supported:

* `expiration_date` - (Required) A valid `RFC3339Time` for the expiration of the token..

* `hostpool_id` - (Required) The ID of the Virtual Desktop Host Pool to link the Registration Info to. Changing this forces a new Registration Info resource to be created. Only a single virtual_desktop_host_pool_registration_info resource should be associated with a given hostpool. Assigning multiple resources will produce inconsistent results.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Host Pool Registration Info resource.

* `token` - The registration token generated by the Virtual Desktop Host Pool for registration of session hosts.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the AVD Registration Info.
* `read` - (Defaults to 5 minutes) Used when retrieving the AVD Registration Info.
* `update` - (Defaults to 30 minutes) Used when updating the AVD Registration Info.
* `delete` - (Defaults to 30 minutes) Used when deleting the AVD Registration Info.

## Import

AVD Registration Infos can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_host_pool_registration_info.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1/registrationInfo/default
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DesktopVirtualization` - 2024-04-03
