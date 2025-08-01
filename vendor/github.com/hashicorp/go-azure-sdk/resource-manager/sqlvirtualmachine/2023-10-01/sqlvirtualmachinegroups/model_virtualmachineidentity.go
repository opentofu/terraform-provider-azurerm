package sqlvirtualmachinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineIdentity struct {
	ResourceId *string         `json:"resourceId,omitempty"`
	Type       *VMIdentityType `json:"type,omitempty"`
}
