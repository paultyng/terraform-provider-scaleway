package scaleway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayInstanceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayInstanceVolumeCreate,
		ReadContext:   resourceScalewayInstanceVolumeRead,
		UpdateContext: resourceScalewayInstanceVolumeUpdate,
		DeleteContext: resourceScalewayInstanceVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultInstanceVolumeDeleteTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the volume",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The volume type",
				ValidateFunc: validation.StringInSlice([]string{
					instance.VolumeVolumeTypeBSSD.String(),
					instance.VolumeVolumeTypeLSSD.String(),
				}, false),
			},
			"size_in_gb": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The size of the volume in gigabyte",
				ConflictsWith: []string{"from_snapshot_id", "from_volume_id"},
			},
			"from_volume_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "Create a copy of an existing volume",
				ValidateFunc:  validationUUIDorUUIDWithLocality(),
				ConflictsWith: []string{"from_snapshot_id", "size_in_gb"},
			},
			"from_snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "Create a volume based on a image",
				ValidateFunc:  validationUUIDorUUIDWithLocality(),
				ConflictsWith: []string{"from_volume_id", "size_in_gb"},
			},
			"server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server associated with this volume",
			},
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
			"zone":            zoneSchema(),
		},
	}
}

func resourceScalewayInstanceVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createVolumeRequest := &instance.CreateVolumeRequest{
		Zone:       zone,
		Name:       expandOrGenerateString(d.Get("name"), "vol"),
		VolumeType: instance.VolumeVolumeType(d.Get("type").(string)),
		Project:    expandStringPtr(d.Get("project_id")),
	}

	if definedZone, ok := d.GetOk("zone"); ok {
		zone = scw.Zone(definedZone.(string))
		createVolumeRequest.Zone = zone
	}

	if size, ok := d.GetOk("size_in_gb"); ok {
		volumeSizeInBytes := scw.Size(uint64(size.(int)) * gb)
		createVolumeRequest.Size = &volumeSizeInBytes
	}

	if volumeID, ok := d.GetOk("from_volume_id"); ok {
		createVolumeRequest.BaseVolume = expandStringPtr(expandID(volumeID))
	}

	if snapshotID, ok := d.GetOk("from_snapshot_id"); ok {
		createVolumeRequest.BaseSnapshot = expandStringPtr(expandID(snapshotID))
	}

	res, err := instanceAPI.CreateVolume(createVolumeRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("couldn't create volume: %s", err))
	}

	d.SetId(newZonedIDString(zone, res.Volume.ID))

	return resourceScalewayInstanceVolumeRead(ctx, d, meta)
}

func resourceScalewayInstanceVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, id, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := instanceAPI.GetVolume(&instance.GetVolumeRequest{
		VolumeID: id,
		Zone:     zone,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("couldn't read volume: %v", err))
	}

	_ = d.Set("name", res.Volume.Name)
	_ = d.Set("organization_id", res.Volume.Organization)
	_ = d.Set("project_id", res.Volume.Project)
	_ = d.Set("zone", string(zone))
	_ = d.Set("type", res.Volume.VolumeType.String())
	_ = d.Set("size_in_gb", int(res.Volume.Size/scw.GB))

	if res.Volume.Server != nil {
		_ = d.Set("server_id", res.Volume.Server.ID)
	} else {
		_ = d.Set("server_id", nil)
	}

	return nil
}

func resourceScalewayInstanceVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, id, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)

		_, err = instanceAPI.UpdateVolume(&instance.UpdateVolumeRequest{
			VolumeID: id,
			Zone:     zone,
			Name:     &newName,
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't update volume: %s", err))
		}
	}

	if d.HasChange("size_in_gb") {
		if d.Get("type") != instance.VolumeVolumeTypeBSSD.String() {
			return diag.FromErr(fmt.Errorf("only block volume can be resized"))
		}
		if oldSize, newSize := d.GetChange("size_in_gb"); oldSize.(int) > newSize.(int) {
			return diag.FromErr(fmt.Errorf("block volumes cannot be resized down"))
		}
		_, err := instanceAPI.WaitForVolume(&instance.WaitForVolumeRequest{
			VolumeID: id,
			Zone:     zone}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		volumeSizeInBytes := scw.Size(uint64(d.Get("size_in_gb").(int)) * gb)
		_, err = instanceAPI.UpdateVolume(&instance.UpdateVolumeRequest{
			VolumeID: id,
			Zone:     zone,
			Size:     &volumeSizeInBytes,
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't resize volume: %s", err))
		}
	}

	return resourceScalewayInstanceVolumeRead(ctx, d, meta)
}

func resourceScalewayInstanceVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, id, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		volumeResp, err := instanceAPI.GetVolume(&instance.GetVolumeRequest{
			Zone:     zone,
			VolumeID: id,
		})
		if err != nil {
			if is404Error(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if volumeResp.Volume.Server != nil {
			return resource.RetryableError(fmt.Errorf("volume is still attached to a server"))
		}

		deleteRequest := &instance.DeleteVolumeRequest{
			Zone:     zone,
			VolumeID: id,
		}

		err = instanceAPI.DeleteVolume(deleteRequest, scw.WithContext(ctx))
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
