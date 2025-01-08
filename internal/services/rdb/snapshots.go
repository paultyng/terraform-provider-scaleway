package rdb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceRdbSnapshotCreate,
		ReadContext:   ResourceRdbSnapshotRead,
		DeleteContext: ResourceRdbSnapshotDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultInstanceTimeout),
			Read:   schema.DefaultTimeout(defaultInstanceTimeout),
			Delete: schema.DefaultTimeout(defaultInstanceTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
				Description:      "UUID of the Database Instance on which the snapshot is applied.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the snapshot.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration date of the snapshot in ISO 8601 format (RFC 3339).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the snapshot.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the snapshot in bytes.",
			},
			"region": regional.Schema(),
		},
		CustomizeDiff: cdf.LocalityCheck("instance_id"),
	}
}

func ResourceRdbSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return ResourceRdbSnapshotRead(ctx, d, meta)
}

func ResourceRdbSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func ResourceRdbSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
