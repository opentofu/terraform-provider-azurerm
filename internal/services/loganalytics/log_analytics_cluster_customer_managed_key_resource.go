// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsClusterCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsClusterCustomerManagedKeyCreate,
		Read:   resourceLogAnalyticsClusterCustomerManagedKeyRead,
		Update: resourceLogAnalyticsClusterCustomerManagedKeyUpdate,
		Delete: resourceLogAnalyticsClusterCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ClusterCustomerManagedKeyV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"log_analytics_cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: clusters.ValidateClusterID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
			},
		},
	}
}

func resourceLogAnalyticsClusterCustomerManagedKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Get("log_analytics_cluster_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	if props.KeyVaultProperties != nil {
		if keyProps := *props.KeyVaultProperties; keyProps.KeyName != nil && *keyProps.KeyName != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_cluster_customer_managed_key", id.ID())
		}
	}

	// Ensure `associatedWorkspaces` is not present in request, this is a read only property and cannot be sent to the API
	// Error: updating Customer Managed Key for Cluster
	//		performing CreateOrUpdate: unexpected status 400 (400 Bad Request) with error:
	//		InvalidParameter: 'properties.associatedWorkspaces' is a read only property and cannot be set.
	//		Please refer to https://docs.microsoft.com/en-us/azure/azure-monitor/log-query/logs-dedicated-clusters#link-a-workspace-to-the-cluster for more information on how to associate a workspace to the cluster.
	props.AssociatedWorkspaces = nil

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Get("key_vault_key_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Key Vault Key ID: %+v", err)
	}

	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(keyId.KeyVaultBaseUrl),
		KeyName:     pointer.To(keyId.Name),
		KeyVersion:  pointer.To(keyId.Version),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
	}

	updateWait, err := logAnalyticsClusterWaitForState(ctx, client, *id)
	if err != nil {
		return err
	}
	if _, err := updateWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish adding Customer Managed Key: %+v", *id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}

func resourceLogAnalyticsClusterCustomerManagedKeyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Get("key_vault_key_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Key Vault Key ID: %+v", err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	if props := model.Properties; props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	// This is a read only property, please see comment in the create function.
	model.Properties.AssociatedWorkspaces = nil

	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(keyId.KeyVaultBaseUrl),
		KeyName:     pointer.To(keyId.Name),
		KeyVersion:  pointer.To(keyId.Version),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
	}

	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}

func resourceLogAnalyticsClusterCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keyVaultKeyId := ""
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if kvProps := props.KeyVaultProperties; kvProps != nil {
				keyVaultUri := pointer.From(kvProps.KeyVaultUri)
				keyName := pointer.From(kvProps.KeyName)
				keyVersion := pointer.From(kvProps.KeyVersion)

				if keyVaultUri != "" && keyName != "" {
					keyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, keyVaultParse.NestedItemTypeKey, keyName, keyVersion)
					if err != nil {
						return err
					}
					keyVaultKeyId = keyId.ID()
				}
			}
		}
	}

	if keyVaultKeyId == "" {
		log.Printf("[DEBUG] %s has no Customer Managed Key - removing from state", *id)
		d.SetId("")
		return nil
	}

	d.Set("log_analytics_cluster_id", d.Id())
	d.Set("key_vault_key_id", keyVaultKeyId)

	return nil
}

func resourceLogAnalyticsClusterCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	if props.KeyVaultProperties == nil {
		return fmt.Errorf("deleting `azurerm_log_analytics_cluster_customer_managed_key` %s: `customer managed key does not exist!`", *id)
	}

	if props.KeyVaultProperties.KeyName == nil || *props.KeyVaultProperties.KeyName == "" {
		return fmt.Errorf("deleting `azurerm_log_analytics_cluster_customer_managed_key` %s: `customer managed key does not exist!`", *id)
	}

	// This is a read only property, please see comment in the create function.
	props.AssociatedWorkspaces = nil

	// The API only removes the CMK when it is sent empty string values, sending nil for each property or an empty object does not work.
	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(""),
		KeyName:     pointer.To(""),
		KeyVersion:  pointer.To(""),
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("deleting Customer Managed Key from %s: %+v", *id, err)
	}

	return nil
}
