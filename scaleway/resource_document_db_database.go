package scaleway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	document_db "github.com/scaleway/scaleway-sdk-go/api/document_db/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayDocumentDBDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayDocumentDBDatabaseCreate,
		ReadContext:   resourceScalewayDocumentDBDatabaseRead,
		DeleteContext: resourceScalewayDocumentDBDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultRdbInstanceTimeout),
			Delete:  schema.DefaultTimeout(defaultRdbInstanceTimeout),
			Default: schema.DefaultTimeout(defaultRdbInstanceTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationUUIDorUUIDWithLocality(),
				Description:  "Instance on which the database is created",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The database name",
			},
			"region":     regionSchema(),
			"project_id": projectIDSchema(),
		},
	}
}

func resourceScalewayDocumentDBDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, region, err := documentDBAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := expandID(d.Get("instance_id"))

	_, err = waitForDocumentDBInstance(ctx, api, region, instanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	database, err := api.CreateDatabase(&document_db.CreateDatabaseRequest{
		Region:     region,
		InstanceID: instanceID,
		Name:       expandOrGenerateString(d.Get("name").(string), "database"),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceScalewayDocumentDBDatabaseID(region, instanceID, database.Name))

	_, err = waitForDocumentDBInstance(ctx, api, region, instanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayDocumentDBDatabaseRead(ctx, d, meta)
}

func getDocumentDBDatabase(ctx context.Context, api *document_db.API, region scw.Region, instanceID string, dbName string) (*document_db.Database, error) {
	res, err := api.ListDatabases(&document_db.ListDatabasesRequest{
		Region:     region,
		InstanceID: instanceID,
		Name:       &dbName,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if len(res.Databases) == 0 {
		return nil, fmt.Errorf("database %q not found", dbName)
	}

	return res.Databases[0], nil
}

func resourceScalewayDocumentDBDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceLocalizedID, databaseName, err := resourceScalewayDocumentDBDatabaseName(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	api, region, instanceID, err := documentDBAPIWithRegionAndID(meta, instanceLocalizedID)
	if err != nil {
		return diag.FromErr(err)
	}

	instance, err := waitForDocumentDBInstance(ctx, api, region, instanceID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	database, err := getDocumentDBDatabase(ctx, api, region, instanceID, databaseName)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", database.Name)
	_ = d.Set("region", instance.Region)
	_ = d.Set("project_id", instance.ProjectID)

	return nil
}

func resourceScalewayDocumentDBDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceLocalizedID, databaseName, err := resourceScalewayDocumentDBDatabaseName(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	api, region, instanceID, err := documentDBAPIWithRegionAndID(meta, instanceLocalizedID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForDocumentDBInstance(ctx, api, region, instanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = api.DeleteDatabase(&document_db.DeleteDatabaseRequest{
		Region:     region,
		Name:       databaseName,
		InstanceID: instanceID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForDocumentDBInstance(ctx, api, region, instanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
