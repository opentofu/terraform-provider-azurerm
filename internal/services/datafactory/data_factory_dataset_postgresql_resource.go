// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryDatasetPostgreSQL() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetPostgreSQLCreateUpdate,
		Read:   resourceDataFactoryDatasetPostgreSQLRead,
		Update: resourceDataFactoryDatasetPostgreSQLCreateUpdate,
		Delete: resourceDataFactoryDatasetPostgreSQLDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"linked_service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"schema_column": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Byte",
								"Byte[]",
								"Boolean",
								"Date",
								"DateTime",
								"DateTimeOffset",
								"Decimal",
								"Double",
								"Guid",
								"Int16",
								"Int32",
								"Int64",
								"Single",
								"String",
								"TimeSpan",
							}, false),
						},
						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceDataFactoryDatasetPostgreSQLCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).DataFactory.DatasetClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_postgresql", id.ID())
		}
	}

	postgresqlDatasetProperties := datafactory.RelationalTableDatasetTypeProperties{
		TableName: d.Get("table_name").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)
	postgresqlTableset := datafactory.RelationalTableDataset{
		RelationalTableDatasetTypeProperties: &postgresqlDatasetProperties,
		LinkedServiceName:                    linkedService,
		Description:                          &description,
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		postgresqlTableset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		postgresqlTableset.Parameters = expandDataSetParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		postgresqlTableset.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		postgresqlTableset.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("schema_column"); ok {
		postgresqlTableset.Structure = expandDataFactoryDatasetStructure(v.([]interface{}))
	}

	datasetType := string(datafactory.TypeBasicDatasetTypePostgreSQLTable)
	dataset := datafactory.DatasetResource{
		Properties: &postgresqlTableset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryDatasetPostgreSQLRead(d, meta)
}

func resourceDataFactoryDatasetPostgreSQLRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	postgresqlTable, ok := resp.Properties.AsRelationalTableDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset PostgreSQL %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeRelationalTable, resp.Properties)
	}

	d.Set("additional_properties", postgresqlTable.AdditionalProperties)

	if postgresqlTable.Description != nil {
		d.Set("description", postgresqlTable.Description)
	}

	parameters := flattenDataSetParameters(postgresqlTable.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(postgresqlTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := postgresqlTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := postgresqlTable.RelationalTableDatasetTypeProperties; properties != nil {
		val, ok := properties.TableName.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `table_name` since it's not a string")
		} else {
			d.Set("table_name", val)
		}
	}

	if folder := postgresqlTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	structureColumns := flattenDataFactoryStructureColumns(postgresqlTable.Structure)
	if err := d.Set("schema_column", structureColumns); err != nil {
		return fmt.Errorf("setting `schema_column`: %+v", err)
	}

	return nil
}

func resourceDataFactoryDatasetPostgreSQLDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
