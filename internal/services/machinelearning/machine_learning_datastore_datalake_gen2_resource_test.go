// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningDataStoreDataLakeGen2 struct{}

func TestAccMachineLearningDataStoreDataLakeGen2_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen2", "test")
	r := MachineLearningDataStoreDataLakeGen2{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2Basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningDataStoreDataLakeGen2_spn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen2", "test")
	r := MachineLearningDataStoreDataLakeGen2{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2Spn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccMachineLearningDataStoreDataLakeGen2_crossSubStorageAccount(t *testing.T) {
	if os.Getenv("ARM_SUBSCRIPTION_ID_ALT") == "" {
		t.Skip("ARM_SUBSCRIPTION_ID_ALT not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen2", "test")
	r := MachineLearningDataStoreDataLakeGen2{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2CrossSubStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningDataStoreDataLakeGen2_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen2", "test")
	r := MachineLearningDataStoreDataLakeGen2{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2Basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
		{
			Config: r.dataLakeGen2Spn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccMachineLearningDataStoreDataLakeGen2_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen2", "test")
	r := MachineLearningDataStoreDataLakeGen2{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2Basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r MachineLearningDataStoreDataLakeGen2) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	dataStoreClient := client.MachineLearning.Datastore
	id, err := datastore.ParseDataStoreID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := dataStoreClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Data Store %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (r MachineLearningDataStoreDataLakeGen2) dataLakeGen2Basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore_datalake_gen2" "test" {
  name                 = "accdatastore%[2]d"
  workspace_id         = azurerm_machine_learning_workspace.test.id
  storage_container_id = azurerm_storage_container.test.resource_manager_id
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStoreDataLakeGen2) dataLakeGen2Spn(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
resource "azuread_application" "test" {
  display_name = "acctestspa-%[2]d"
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application.test.client_id
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = azuread_service_principal.test.object_id
}

resource "azurerm_machine_learning_datastore_datalake_gen2" "test" {
  name                 = "accdatastore%[2]d"
  workspace_id         = azurerm_machine_learning_workspace.test.id
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  tenant_id            = azuread_service_principal.test.application_tenant_id
  client_id            = azuread_service_principal.test.client_id
  client_secret        = azuread_service_principal_password.test.value
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStoreDataLakeGen2) dataLakeGen2CrossSubStorageAccount(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore_datalake_gen2" "test" {
  name                 = "acctestdatastore%[2]d"
  workspace_id         = azurerm_machine_learning_workspace.test.id
  storage_container_id = azurerm_storage_container.test.resource_manager_id
}

provider "azurerm-alt" {
  subscription_id = "%[3]s"

  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "testalt" {
  provider = azurerm-alt
  name     = "acctestRG-alt-%[2]d"
  location = "%[4]s"
}

resource "azurerm_storage_account" "testalt" {
  provider                 = azurerm-alt
  name                     = "acctestsaalt%[5]d"
  location                 = azurerm_resource_group.testalt.location
  resource_group_name      = azurerm_resource_group.testalt.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "testalt" {
  provider              = azurerm-alt
  name                  = "acctestcontaineralt%[5]d"
  storage_account_name  = azurerm_storage_account.testalt.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore_datalake_gen2" "crosssub" {
  name                 = "acctestdcrosssub%[5]d"
  workspace_id         = azurerm_machine_learning_workspace.test.id
  storage_container_id = azurerm_storage_container.testalt.resource_manager_id
}
	`, template, data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID_ALT"), data.Locations.Primary, data.RandomIntOfLength(10))
}

func (r MachineLearningDataStoreDataLakeGen2) requiresImport(data acceptance.TestData) string {
	template := r.dataLakeGen2Basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_datastore_datalake_gen2" "import" {
  name                 = azurerm_machine_learning_datastore_datalake_gen2.test.name
  workspace_id         = azurerm_machine_learning_datastore_datalake_gen2.test.workspace_id
  storage_container_id = azurerm_machine_learning_datastore_datalake_gen2.test.storage_container_id
}
`, template)
}

func (r MachineLearningDataStoreDataLakeGen2) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomIntOfLength(15))
}
