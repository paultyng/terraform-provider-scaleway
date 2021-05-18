package scaleway

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayRdbDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayRdbDatabaseBackupCreate,
		ReadContext:   resourceScalewayRdbDatabaseBackupRead,
		UpdateContext: resourceScalewayRdbDatabaseBackupUpdate,
		DeleteContext: resourceScalewayRdbDatabaseBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultRdbInstanceTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationUUIDorUUIDWithLocality(),
				Description:  "Instance on which the user is created",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Database user name",
				Required:    true,
				ForceNew:    true,
			},
			// Common
			"region": regionSchema(),
		},
	}
}

func resourceScalewayRdbDatabaseBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, err := rdbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get("instance_id").(string)
	createReq := &rdb.CreateDatabaseBackupRequest{
		Region:       region,
		InstanceID:   expandID(instanceID),
		DatabaseName: "",
		Name:         d.Get("name").(string),
		ExpiresAt:    nil,
	}

	res, err := rdbAPI.CreateDatabaseBackup(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceScalewayRdbDatabaseBackupID(region, expandID(instanceID), res.Name))

	return resourceScalewayRdbDatabaseBackupRead(ctx, d, meta)
}

func resourceScalewayRdbDatabaseBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, err := rdbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID, userName, err := resourceScalewayRdbDatabaseBackupParseID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	res, err := rdbAPI.ListDatabaseBackups(&rdb.ListDatabaseBackupsRequest{
		Region:     region,
		Name:       &userName,
		OrderBy:    "",
		InstanceID: scw.StringPtr(instanceID),
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var dbBackup = res.DatabaseBackups[0]
	_ = d.Set("instance_id", newRegionalID(region, instanceID).String())
	_ = d.Set("name", dbBackup.Name)

	d.SetId(resourceScalewayRdbUserID(region, instanceID, dbBackup.Name))

	return nil
}

func resourceScalewayRdbDatabaseBackupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, err := rdbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	_, dbBackupName, err := resourceScalewayRdbDatabaseBackupParseID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	req := &rdb.UpdateDatabaseBackupRequest{
		Region:           region,
		DatabaseBackupID: d.Id(),
		Name:             scw.StringPtr(dbBackupName),
		ExpiresAt:        nil,
	}

	_, err = rdbAPI.UpdateDatabaseBackup(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayRdbDatabaseBackupRead(ctx, d, meta)
}

func resourceScalewayRdbDatabaseBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, err := rdbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = rdbAPI.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
		DatabaseBackupID: d.Id(),
		Region:           region,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = rdbAPI.DeleteDatabaseBackup(&rdb.DeleteDatabaseBackupRequest{
		DatabaseBackupID: d.Id(),
		Region:           region,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}

// Build the resource identifier
// The resource identifier format is "Region/InstanceId/UserName"
func resourceScalewayRdbDatabaseBackupID(region scw.Region, instanceID string, userName string) (resourceID string) {
	return fmt.Sprintf("%s/%s/%s", region, instanceID, userName)
}

// Extract instance ID and backup name from the resource identifier.
// The resource identifier format is "Region/InstanceId/UserName"
func resourceScalewayRdbDatabaseBackupParseID(resourceID string) (instanceID string, backupName string, err error) {
	idParts := strings.Split(resourceID, "/")
	if len(idParts) != 3 {
		return "", "", fmt.Errorf("can't parse user resource id: %s", resourceID)
	}
	return idParts[1], idParts[2], nil
}
