package edgeservices

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	edgeservices "github.com/scaleway/scaleway-sdk-go/api/edge_services/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func ResourcePipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourcePipelineCreate,
		ReadContext:   ResourcePipelineRead,
		UpdateContext: ResourcePipelineUpdate,
		DeleteContext: ResourcePipelineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultEdgeServicesTimeout),
			Read:    schema.DefaultTimeout(defaultEdgeServicesTimeout),
			Update:  schema.DefaultTimeout(defaultEdgeServicesTimeout),
			Delete:  schema.DefaultTimeout(defaultEdgeServicesTimeout),
			Default: schema.DefaultTimeout(defaultEdgeServicesTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The pipeline name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The pipeline description",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The pipeline description",
			},
			"dns_stage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The pipeline description",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The pipeline description",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The pipeline description",
			},
			"project_id": account.ProjectIDSchema(),
		},
	}
}

func ResourcePipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := NewEdgeServicesAPI(m)

	pipeline, err := api.CreatePipeline(&edgeservices.CreatePipelineRequest{
		Description: d.Get("description").(string),
		ProjectID:   d.Get("project_id").(string),
		Name:        d.Get("name").(string),
		DNSStageID:  types.ExpandStringPtr(locality.ExpandID(d.Get("dns_stage_id").(string))),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pipeline.ID)

	return ResourcePipelineRead(ctx, d, m)
}

func ResourcePipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := NewEdgeServicesAPI(m)

	pipeline, err := api.GetPipeline(&edgeservices.GetPipelineRequest{
		PipelineID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", pipeline.Name)
	_ = d.Set("description", pipeline.Description)
	_ = d.Set("dns_stage_id", types.FlattenStringPtr(pipeline.DNSStageID))
	_ = d.Set("created_at", types.FlattenTime(pipeline.CreatedAt))
	_ = d.Set("updated_at", types.FlattenTime(pipeline.UpdatedAt))
	_ = d.Set("status", pipeline.Status.String())
	_ = d.Set("project_id", pipeline.ProjectID)

	return nil
}

func ResourcePipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := NewEdgeServicesAPI(m)

	hasChanged := false

	updateRequest := &edgeservices.UpdatePipelineRequest{
		PipelineID: d.Id(),
	}

	if d.HasChange("name") {
		updateRequest.Name = types.ExpandUpdatedStringPtr(d.Get("name"))
		hasChanged = true
	}

	if d.HasChange("description") {
		updateRequest.Description = types.ExpandUpdatedStringPtr(d.Get("description"))
		hasChanged = true
	}

	if d.HasChange("dns_stage_id") {
		updateRequest.DNSStageID = types.ExpandUpdatedStringPtr(d.Get("dns_stage_id"))
		hasChanged = true
	}

	if hasChanged {
		_, err := api.UpdatePipeline(updateRequest, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourcePipelineRead(ctx, d, m)
}

func ResourcePipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := NewEdgeServicesAPI(m)

	err := api.DeletePipeline(&edgeservices.DeletePipelineRequest{
		PipelineID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}
