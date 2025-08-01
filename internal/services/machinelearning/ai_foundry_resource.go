// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyvaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AIFoundry struct{}

type AIFoundryModel struct {
	Name                        string                                     `tfschema:"name"`
	Location                    string                                     `tfschema:"location"`
	ResourceGroupName           string                                     `tfschema:"resource_group_name"`
	ApplicationInsightsId       string                                     `tfschema:"application_insights_id"`
	StorageAccountId            string                                     `tfschema:"storage_account_id"`
	KeyVaultId                  string                                     `tfschema:"key_vault_id"`
	ContainerRegistryId         string                                     `tfschema:"container_registry_id"`
	Encryption                  []Encryption                               `tfschema:"encryption"`
	ManagedNetwork              []ManagedNetwork                           `tfschema:"managed_network"`
	PublicNetworkAccess         string                                     `tfschema:"public_network_access"`
	Identity                    []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PrimaryUserAssignedIdentity string                                     `tfschema:"primary_user_assigned_identity"`
	HighBusinessImpactEnabled   bool                                       `tfschema:"high_business_impact_enabled"`
	Description                 string                                     `tfschema:"description"`
	FriendlyName                string                                     `tfschema:"friendly_name"`
	DiscoveryUrl                string                                     `tfschema:"discovery_url"`
	WorkspaceId                 string                                     `tfschema:"workspace_id"`
	Tags                        map[string]interface{}                     `tfschema:"tags"`
}

type ManagedNetwork struct {
	IsolationMode string `tfschema:"isolation_mode"`
}

type Encryption struct {
	IdentityClientID string `tfschema:"user_assigned_identity_id"`
	KeyVaultID       string `tfschema:"key_vault_id"`
	KeyID            string `tfschema:"key_id"`
}

func (r AIFoundry) ModelObject() interface{} {
	return &AIFoundryModel{}
}

func (r AIFoundry) ResourceType() string {
	return "azurerm_ai_foundry"
}

func (r AIFoundry) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

func (r AIFoundry) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.MachineLearning.Workspaces
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Kind == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if !strings.EqualFold(*resp.Model.Kind, "Hub") {
			return fmt.Errorf("importing %s: specified workspace is not of kind `Hub`, got `%s`", id, *resp.Model.Kind)
		}

		return nil
	}
}

var _ sdk.ResourceWithUpdate = AIFoundry{}

var _ sdk.ResourceWithCustomImporter = AIFoundry{}

func (r AIFoundry) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WorkspaceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"high_business_impact_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// NOTE: O+C creating a hub that has encryption enabled will set this property to true
			Computed: true,
			ForceNew: true,
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
					"key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvaultValidate.NestedItemId,
					},
					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						// Can be removed when https://github.com/Azure/azure-rest-api-specs/issues/30625 has been fixed
						DiffSuppressFunc: suppress.CaseDifference,
					},
				},
			},
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: components.ValidateComponentID,
		},

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},

		"managed_network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"isolation_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForIsolationMode(), false),
					},
				},
			},
		},

		"primary_user_assigned_identity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      workspaces.PublicNetworkAccessEnabled,
			ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForPublicNetworkAccess(), false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AIFoundry) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r AIFoundry) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AIFoundryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := workspaces.NewWorkspaceID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_ai_foundry", id.ID())
			}

			storageAccountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			keyVaultId, err := commonids.ParseKeyVaultID(model.KeyVaultId)
			if err != nil {
				return err
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			payload := workspaces.Workspace{
				Name:     pointer.To(id.WorkspaceName),
				Location: pointer.To(location.Normalize(model.Location)),
				Identity: expandedIdentity,
				Tags:     tags.Expand(model.Tags),
				Kind:     pointer.To("Hub"),
				Properties: &workspaces.WorkspaceProperties{
					KeyVault:            pointer.To(keyVaultId.ID()),
					PublicNetworkAccess: pointer.To(workspaces.PublicNetworkAccess(model.PublicNetworkAccess)),
					StorageAccount:      pointer.To(storageAccountId.ID()),
				},
			}

			if model.ApplicationInsightsId != "" {
				applicationInsightsId, err := components.ParseComponentID(model.ApplicationInsightsId)
				if err != nil {
					return err
				}
				payload.Properties.ApplicationInsights = pointer.To(applicationInsightsId.ID())
			}

			if model.ContainerRegistryId != "" {
				containerRegistryId, err := registries.ParseRegistryID(model.ContainerRegistryId)
				if err != nil {
					return err
				}
				payload.Properties.ContainerRegistry = pointer.To(containerRegistryId.ID())
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if model.FriendlyName != "" {
				payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
			}

			if model.HighBusinessImpactEnabled {
				payload.Properties.HbiWorkspace = pointer.To(model.HighBusinessImpactEnabled)
			}

			if model.PrimaryUserAssignedIdentity != "" {
				userAssignedId, err := commonids.ParseUserAssignedIdentityID(model.PrimaryUserAssignedIdentity)
				if err != nil {
					return err
				}
				payload.Properties.PrimaryUserAssignedIdentity = pointer.To(userAssignedId.ID())
			}

			if len(model.Encryption) > 0 {
				encryption := expandEncryption(model.Encryption)
				payload.Properties.Encryption = encryption
			}

			if len(model.ManagedNetwork) > 0 {
				payload.Properties.ManagedNetwork = expandManagedNetwork(model.ManagedNetwork)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AIFoundry) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AIFoundryModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("application_insights_id") {
				applicationInsightsId, err := components.ParseComponentID(state.ApplicationInsightsId)
				if err != nil {
					return err
				}
				payload.Properties.ApplicationInsights = pointer.To(applicationInsightsId.ID())
			}

			if metadata.ResourceData.HasChange("container_registry_id") {
				containerRegistryId, err := registries.ParseRegistryID(state.ContainerRegistryId)
				if err != nil {
					return err
				}
				payload.Properties.ContainerRegistry = pointer.To(containerRegistryId.ID())
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.To(workspaces.PublicNetworkAccess(state.PublicNetworkAccess))
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = pointer.To(state.FriendlyName)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("primary_user_assigned_identity") {
				userAssignedId, err := commonids.ParseUserAssignedIdentityID(state.PrimaryUserAssignedIdentity)
				if err != nil {
					return err
				}
				payload.Properties.PrimaryUserAssignedIdentity = pointer.To(userAssignedId.ID())
			}

			if metadata.ResourceData.HasChange("managed_network") {
				payload.Properties.ManagedNetwork = expandManagedNetwork(state.ManagedNetwork)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(state.Tags)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AIFoundry) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			hub := AIFoundryModel{
				Name:              id.WorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				hub.Location = location.NormalizeNilable(model.Location)

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				hub.Identity = flattenedIdentity
				hub.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if v := pointer.From(props.ApplicationInsights); v != "" {
						applicationInsightsId, err := components.ParseComponentIDInsensitively(v)
						if err != nil {
							return err
						}
						hub.ApplicationInsightsId = applicationInsightsId.ID()
					}

					if v := pointer.From(props.ContainerRegistry); v != "" {
						containerRegistryId, err := registries.ParseRegistryID(v)
						if err != nil {
							return err
						}
						hub.ContainerRegistryId = containerRegistryId.ID()
					}

					storageAccountId, err := commonids.ParseStorageAccountID(*props.StorageAccount)
					if err != nil {
						return err
					}
					hub.StorageAccountId = storageAccountId.ID()

					if props.KeyVault != nil {
						keyVaultId, err := commonids.ParseKeyVaultID(*props.KeyVault)
						if err != nil {
							return err
						}
						hub.KeyVaultId = keyVaultId.ID()
					}

					hub.Description = pointer.From(props.Description)
					hub.FriendlyName = pointer.From(props.FriendlyName)
					hub.HighBusinessImpactEnabled = pointer.From(props.HbiWorkspace)
					hub.PublicNetworkAccess = string(*props.PublicNetworkAccess)
					hub.DiscoveryUrl = pointer.From(props.DiscoveryURL)
					hub.WorkspaceId = pointer.From(props.WorkspaceId)
					hub.ManagedNetwork = flattenManagedNetwork(props.ManagedNetwork)

					if v := pointer.From(props.PrimaryUserAssignedIdentity); v != "" {
						userAssignedId, err := commonids.ParseUserAssignedIdentityID(v)
						if err != nil {
							return err
						}
						hub.PrimaryUserAssignedIdentity = userAssignedId.ID()
					}

					encryption, err := flattenEncryption(props.Encryption)
					if err != nil {
						return fmt.Errorf("flattening `encryption`: %+v", err)
					}
					hub.Encryption = encryption
				}
			}

			return metadata.Encode(&hub)
		},
	}
}

func (r AIFoundry) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			opts := workspaces.DefaultDeleteOperationOptions()

			if metadata.Client.Features.MachineLearning.PurgeSoftDeletedWorkspaceOnDestroy {
				opts.ForceToPurge = pointer.To(true)
			}

			if err := client.DeleteThenPoll(ctx, *id, opts); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandEncryption(input []Encryption) *workspaces.EncryptionProperty {
	encryption := input[0]
	out := workspaces.EncryptionProperty{
		Identity: &workspaces.IdentityForCmk{},
		KeyVaultProperties: workspaces.EncryptionKeyVaultProperties{
			KeyVaultArmId: encryption.KeyVaultID,
			KeyIdentifier: encryption.KeyID,
		},
		Status: workspaces.EncryptionStatusEnabled,
	}

	if encryption.IdentityClientID != "" {
		out.Identity.UserAssignedIdentity = pointer.To(encryption.IdentityClientID)
	}

	return &out
}

func flattenEncryption(input *workspaces.EncryptionProperty) ([]Encryption, error) {
	out := make([]Encryption, 0)

	if input == nil || input.Status != workspaces.EncryptionStatusEnabled {
		return out, nil
	}

	encryption := Encryption{}
	if v := input.KeyVaultProperties.KeyVaultArmId; v != "" {
		keyVaultId, err := commonids.ParseKeyVaultID(v)
		if err != nil {
			return nil, err
		}
		encryption.KeyVaultID = keyVaultId.ID()
	}
	if v := input.KeyVaultProperties.KeyIdentifier; v != "" {
		keyId, err := keyvaultParse.ParseNestedItemID(v)
		if err != nil {
			return nil, err
		}
		encryption.KeyID = keyId.ID()
	}

	if input.Identity != nil && input.Identity.UserAssignedIdentity != nil {
		userAssignedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity.UserAssignedIdentity)
		if err != nil {
			return nil, err
		}
		encryption.IdentityClientID = userAssignedId.ID()
	}

	return append(out, encryption), nil
}

func expandManagedNetwork(input []ManagedNetwork) *workspaces.ManagedNetworkSettings {
	network := input[0]

	return &workspaces.ManagedNetworkSettings{
		IsolationMode: pointer.To(workspaces.IsolationMode(network.IsolationMode)),
	}
}

func flattenManagedNetwork(input *workspaces.ManagedNetworkSettings) []ManagedNetwork {
	out := make([]ManagedNetwork, 0)
	if input == nil {
		return out
	}

	return append(out, ManagedNetwork{
		IsolationMode: string(pointer.From(input.IsolationMode)),
	})
}
