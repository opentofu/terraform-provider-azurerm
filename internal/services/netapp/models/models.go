// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

const (
	MaxQuotaTargetIDSizeInKiB int64 = 4294967295
)

type NetAppAccountEncryption struct {
	NetAppAccountID                   string `tfschema:"netapp_account_id"`
	UserAssignedIdentityID            string `tfschema:"user_assigned_identity_id"`
	SystemAssignedIdentityPrincipalID string `tfschema:"system_assigned_identity_principal_id"`
	EncryptionKey                     string `tfschema:"encryption_key"`
}

type NetAppAccountEncryptionDataSourceModel struct {
	NetAppAccountID                   string `tfschema:"netapp_account_id"`
	UserAssignedIdentityID            string `tfschema:"user_assigned_identity_id"`
	SystemAssignedIdentityPrincipalID string `tfschema:"system_assigned_identity_principal_id"`
	EncryptionKey                     string `tfschema:"encryption_key"`
}

type NetAppVolumeGroupSAPHanaVolume struct {
	Id                           string                         `tfschema:"id"`
	Name                         string                         `tfschema:"name"`
	VolumePath                   string                         `tfschema:"volume_path"`
	ServiceLevel                 string                         `tfschema:"service_level"`
	SubnetId                     string                         `tfschema:"subnet_id"`
	Protocols                    []string                       `tfschema:"protocols"`
	SecurityStyle                string                         `tfschema:"security_style"`
	StorageQuotaInGB             int64                          `tfschema:"storage_quota_in_gb"`
	ThroughputInMibps            float64                        `tfschema:"throughput_in_mibps"`
	Tags                         map[string]string              `tfschema:"tags"`
	SnapshotDirectoryVisible     bool                           `tfschema:"snapshot_directory_visible"`
	CapacityPoolId               string                         `tfschema:"capacity_pool_id"`
	ProximityPlacementGroupId    string                         `tfschema:"proximity_placement_group_id"`
	VolumeSpecName               string                         `tfschema:"volume_spec_name"`
	ExportPolicy                 []ExportPolicyRule             `tfschema:"export_policy_rule"`
	MountIpAddresses             []string                       `tfschema:"mount_ip_addresses"`
	DataProtectionReplication    []DataProtectionReplication    `tfschema:"data_protection_replication"`
	DataProtectionSnapshotPolicy []DataProtectionSnapshotPolicy `tfschema:"data_protection_snapshot_policy"`
}

type NetAppVolumeGroupSAPHanaModel struct {
	Name                  string                           `tfschema:"name"`
	ResourceGroupName     string                           `tfschema:"resource_group_name"`
	Location              string                           `tfschema:"location"`
	AccountName           string                           `tfschema:"account_name"`
	GroupDescription      string                           `tfschema:"group_description"`
	ApplicationIdentifier string                           `tfschema:"application_identifier"`
	Volumes               []NetAppVolumeGroupSAPHanaVolume `tfschema:"volume"`
}

type NetAppVolumeGroupSAPHanaDataSourceModel struct {
	Name                  string                           `tfschema:"name"`
	ResourceGroupName     string                           `tfschema:"resource_group_name"`
	Location              string                           `tfschema:"location"`
	AccountName           string                           `tfschema:"account_name"`
	GroupDescription      string                           `tfschema:"group_description"`
	ApplicationIdentifier string                           `tfschema:"application_identifier"`
	Volumes               []NetAppVolumeGroupSAPHanaVolume `tfschema:"volume"`
}

type NetAppVolumeGroupOracleVolume struct {
	Id                           string                         `tfschema:"id"`
	Name                         string                         `tfschema:"name"`
	VolumePath                   string                         `tfschema:"volume_path"`
	ServiceLevel                 string                         `tfschema:"service_level"`
	SubnetId                     string                         `tfschema:"subnet_id"`
	Protocols                    []string                       `tfschema:"protocols"`
	SecurityStyle                string                         `tfschema:"security_style"`
	StorageQuotaInGB             int64                          `tfschema:"storage_quota_in_gb"`
	ThroughputInMibps            float64                        `tfschema:"throughput_in_mibps"`
	Tags                         map[string]string              `tfschema:"tags"`
	SnapshotDirectoryVisible     bool                           `tfschema:"snapshot_directory_visible"`
	CapacityPoolId               string                         `tfschema:"capacity_pool_id"`
	ProximityPlacementGroupId    string                         `tfschema:"proximity_placement_group_id"`
	VolumeSpecName               string                         `tfschema:"volume_spec_name"`
	ExportPolicy                 []ExportPolicyRule             `tfschema:"export_policy_rule"`
	MountIpAddresses             []string                       `tfschema:"mount_ip_addresses"`
	DataProtectionReplication    []DataProtectionReplication    `tfschema:"data_protection_replication"`
	DataProtectionSnapshotPolicy []DataProtectionSnapshotPolicy `tfschema:"data_protection_snapshot_policy"`
	Zone                         string                         `tfschema:"zone"`
	EncryptionKeySource          string                         `tfschema:"encryption_key_source"`
	KeyVaultPrivateEndpointId    string                         `tfschema:"key_vault_private_endpoint_id"`
	NetworkFeatures              string                         `tfschema:"network_features"`
}

type NetAppVolumeGroupOracleModel struct {
	Name                  string                          `tfschema:"name"`
	ResourceGroupName     string                          `tfschema:"resource_group_name"`
	Location              string                          `tfschema:"location"`
	AccountName           string                          `tfschema:"account_name"`
	GroupDescription      string                          `tfschema:"group_description"`
	ApplicationIdentifier string                          `tfschema:"application_identifier"`
	Volumes               []NetAppVolumeGroupOracleVolume `tfschema:"volume"`
}

type NetAppVolumeGroupOracleDataSourceModel struct {
	Name                  string                          `tfschema:"name"`
	ResourceGroupName     string                          `tfschema:"resource_group_name"`
	Location              string                          `tfschema:"location"`
	AccountName           string                          `tfschema:"account_name"`
	GroupDescription      string                          `tfschema:"group_description"`
	ApplicationIdentifier string                          `tfschema:"application_identifier"`
	Volumes               []NetAppVolumeGroupOracleVolume `tfschema:"volume"`
}

type ExportPolicyRule struct {
	RuleIndex         int64  `tfschema:"rule_index"`
	AllowedClients    string `tfschema:"allowed_clients"`
	Nfsv3Enabled      bool   `tfschema:"nfsv3_enabled"`
	Nfsv41Enabled     bool   `tfschema:"nfsv41_enabled"`
	UnixReadOnly      bool   `tfschema:"unix_read_only"`
	UnixReadWrite     bool   `tfschema:"unix_read_write"`
	RootAccessEnabled bool   `tfschema:"root_access_enabled"`
}

type DataProtectionReplication struct {
	EndpointType           string `tfschema:"endpoint_type"`
	RemoteVolumeLocation   string `tfschema:"remote_volume_location"`
	RemoteVolumeResourceId string `tfschema:"remote_volume_resource_id"`
	ReplicationFrequency   string `tfschema:"replication_frequency"`
}

type DataProtectionSnapshotPolicy struct {
	DataProtectionSnapshotPolicy string `tfschema:"snapshot_policy_id"`
}

type ReplicationSchedule string

const (
	ReplicationSchedule10Minutes ReplicationSchedule = "10minutes"
	ReplicationScheduleDaily     ReplicationSchedule = "daily"
	ReplicationScheduleHourly    ReplicationSchedule = "hourly"
)

func PossibleValuesForReplicationSchedule() []string {
	return []string{
		string(ReplicationSchedule10Minutes),
		string(ReplicationScheduleDaily),
		string(ReplicationScheduleHourly),
	}
}

type NetAppVolumeQuotaRuleModel struct {
	Name           string `tfschema:"name"`
	Location       string `tfschema:"location"`
	VolumeID       string `tfschema:"volume_id"`
	QuotaTarget    string `tfschema:"quota_target"`
	QuotaSizeInKiB int64  `tfschema:"quota_size_in_kib"`
	QuotaType      string `tfschema:"quota_type"`
}

type NetAppVolumeQuotaRuleDataSourceModel struct {
	Name           string `tfschema:"name"`
	Location       string `tfschema:"location"`
	VolumeID       string `tfschema:"volume_id"`
	QuotaTarget    string `tfschema:"quota_target"`
	QuotaSizeInKiB int64  `tfschema:"quota_size_in_kib"`
	QuotaType      string `tfschema:"quota_type"`
}

type NetAppBackupVaultModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	AccountName       string            `tfschema:"account_name"`
	Tags              map[string]string `tfschema:"tags"`
}

type NetAppBackupPolicyModel struct {
	Name                 string            `tfschema:"name"`
	ResourceGroupName    string            `tfschema:"resource_group_name"`
	Location             string            `tfschema:"location"`
	AccountName          string            `tfschema:"account_name"`
	Tags                 map[string]string `tfschema:"tags"`
	DailyBackupsToKeep   int64             `tfschema:"daily_backups_to_keep"`
	WeeklyBackupsToKeep  int64             `tfschema:"weekly_backups_to_keep"`
	MonthlyBackupsToKeep int64             `tfschema:"monthly_backups_to_keep"`
	Enabled              bool              `tfschema:"enabled"`
}
