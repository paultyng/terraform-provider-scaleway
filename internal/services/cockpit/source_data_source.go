package cockpit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func DataSourceCockpitSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCockpitSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the data source.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the data source.",
			},
			"project_id": account.ProjectIDSchema(),
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the data source.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of the data source (e.g., 'metrics', 'logs', 'traces').",
				ValidateFunc: validation.StringInSlice([]string{
					"metrics", "logs", "traces",
				}, false),
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The origin of the data source (e.g., 'scaleway', 'external', 'custom').",
				ValidateFunc: validation.StringInSlice([]string{
					"scaleway", "external", "custom",
				}, false),
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the data source.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date of the data source.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date of the data source.",
			},
			"synchronized_with_grafana": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the data source is synchronized with Grafana.",
			},
			"retention_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The retention period of the data source in days.",
			},
		},
	}
}

func dataSourceCockpitSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("id"); ok {
		return fetchDataSourceByID(ctx, d, meta)
	}
	return fetchDataSourceByFilters(ctx, d, meta)
}

func fetchDataSourceByID(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	regionalID := d.Get("id").(string)
	api, region, id, err := NewAPIWithRegionAndID(meta, regionalID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	res, err := api.GetDataSource(&cockpit.RegionalAPIGetDataSourceRequest{
		Region:       region,
		DataSourceID: id,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	flattenDataSource(d, res)
	return nil
}

func fetchDataSourceByFilters(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, region, err := cockpitAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &cockpit.RegionalAPIListDataSourcesRequest{
		Region:    region,
		ProjectID: d.Get("project_id").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		req.Types = []cockpit.DataSourceType{cockpit.DataSourceType(v.(string))}
	}
	if v, ok := d.GetOk("origin"); ok {
		req.Origin = cockpit.DataSourceOrigin(v.(string))
	}

	res, err := api.ListDataSources(req, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return diag.FromErr(err)
	}

	if res.TotalCount == 0 {
		return diag.Errorf("no data source found matching the specified criteria")
	}

	if name, ok := d.GetOk("name"); ok {
		for _, ds := range res.DataSources {
			if ds.Name == name.(string) {
				flattenDataSource(d, ds)
				return nil
			}
		}
		return diag.Errorf("no data source found with name '%s'", name.(string))
	}

	flattenDataSource(d, res.DataSources[0])
	return nil
}

func flattenDataSource(d *schema.ResourceData, ds *cockpit.DataSource) {
	d.SetId(regional.NewIDString(ds.Region, ds.ID))
	_ = d.Set("project_id", ds.ProjectID)
	_ = d.Set("name", ds.Name)
	_ = d.Set("url", ds.URL)
	_ = d.Set("type", ds.Type.String())
	_ = d.Set("origin", ds.Origin.String())
	_ = d.Set("created_at", types.FlattenTime(ds.CreatedAt))
	_ = d.Set("updated_at", types.FlattenTime(ds.UpdatedAt))
	_ = d.Set("synchronized_with_grafana", ds.SynchronizedWithGrafana)
	_ = d.Set("retention_days", int(ds.RetentionDays))
	_ = d.Set("region", ds.Region.String())
}
