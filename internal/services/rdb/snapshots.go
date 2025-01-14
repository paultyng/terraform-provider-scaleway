package rdb

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceRdbSnapshotCreate,
		ReadContext:   ResourceRdbSnapshotRead,
		UpdateContext: ResourceRdbSnapshotUpdate,
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
				Computed:    true,
				Description: "Expiration date of the snapshot in ISO 8601 format (RFC 3339).",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration date of the snapshot in ISO 8601 format (RFC 3339).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration date of the snapshot in ISO 8601 format (RFC 3339).",
			},
			"node_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of database instance you want to create",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of volume where data are stored",
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
	rdbAPI, region, err := newAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	_, instanceID, err := regional.ParseID(d.Get("instance_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &rdb.CreateSnapshotRequest{
		Region:     region,
		InstanceID: instanceID,
	}
	if _, ok := d.GetOk("name"); ok {
		createReq.Name = d.Get("name").(string)
	}

	if _, ok := d.GetOk("expires_at"); ok {
		createReq.ExpiresAt = types.ExpandTimePtr(d.Get("expires_at").(string))
	}
	res, err := rdbAPI.CreateSnapshot(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = waitForRDBSnapshot(ctx, rdbAPI, region, res.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(regional.NewIDString(region, res.ID))
	return ResourceRdbSnapshotRead(ctx, d, meta)
}

func ResourceRdbSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := NewAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := rdbAPI.GetSnapshot(&rdb.GetSnapshotRequest{
		SnapshotID: ID,
		Region:     region,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Set resource data fields
	_ = d.Set("instance_id", regional.NewIDString(region, res.InstanceID))
	_ = d.Set("name", res.Name)
	_ = d.Set("expires_at", res.ExpiresAt.Format(time.RFC3339))
	_ = d.Set("created_at", res.CreatedAt.Format(time.RFC3339))
	if res.UpdatedAt != nil {
		_ = d.Set("updated_at", res.UpdatedAt.Format(time.RFC3339))
	}
	_ = d.Set("node_type", res.NodeType)
	_ = d.Set("volume_type", res.VolumeType.Type)
	_ = d.Set("status", res.Status.String())
	if res.Size != nil {
		_ = d.Set("size", int(*res.Size))
	}
	_ = d.Set("region", region)

	return nil
}

func ResourceRdbSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, snapshotID, err := NewAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = rdbAPI.GetSnapshot(&rdb.GetSnapshotRequest{
		SnapshotID: snapshotID,
		Region:     region,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	snapshotUpdateRequest := &rdb.UpdateSnapshotRequest{
		SnapshotID: snapshotID,
	}
	needsUpdate := false

	if d.HasChange("name") {
		name := d.Get("name").(string)
		snapshotUpdateRequest.Name = &name
		needsUpdate = true
	}

	if d.HasChange("expires_at") {
		snapshotUpdateRequest.ExpiresAt = types.ExpandTimePtr(d.Get("expires_at").(string))
		needsUpdate = true
	}

	if needsUpdate {
		_, err = rdbAPI.UpdateSnapshot(snapshotUpdateRequest, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	_, err = waitForRDBSnapshot(ctx, rdbAPI, region, snapshotID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceRdbSnapshotRead(ctx, d, meta)
}

func ResourceRdbSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, snapshotID, err := NewAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = rdbAPI.GetSnapshot(&rdb.GetSnapshotRequest{
		SnapshotID: snapshotID,
		Region:     region,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	_, err = rdbAPI.DeleteSnapshot(&rdb.DeleteSnapshotRequest{
		Region:     region,
		SnapshotID: snapshotID,
	})
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
