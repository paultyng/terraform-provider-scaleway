package edgeservices

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	edgeservices "github.com/scaleway/scaleway-sdk-go/api/edge_services/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func ResourceTLSStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceTLSStageCreate,
		ReadContext:   ResourceTLSStageRead,
		UpdateContext: ResourceTLSStageUpdate,
		DeleteContext: ResourceTLSStageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"backend_stage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The backend stage ID the TLS stage will be linked to",
			},
			"cache_stage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The cache stage ID the TLS stage will be linked to",
			},
			"managed_certificate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Set to true when Scaleway generates and manages a Let's Encrypt certificate for the TLS stage/custom endpoint",
			},
			"secrets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The TLS secrets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the Secret",
						},
						"region": regional.Schema(),
					},
				},
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The pipeline ID the TLS stage belongs to",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the TLS stage",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the TLS stage",
			},
			"certificate_expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TThe expiration date of the certificate",
			},
			"project_id": account.ProjectIDSchema(),
		},
	}
}

func ResourceTLSStageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, region, err := NewEdgeServicesAPIWithRegion(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	tlsStage, err := api.CreateTLSStage(&edgeservices.CreateTLSStageRequest{
		ProjectID:          d.Get("project_id").(string),
		BackendStageID:     types.ExpandStringPtr(d.Get("backend_stage_id").(string)),
		CacheStageID:       types.ExpandStringPtr(d.Get("cache_stage_id").(string)),
		ManagedCertificate: types.ExpandBoolPtr(d.Get("managed_certificate").(bool)),
		Secrets:            expandTLSSecrets(d.Get("secrets"), region),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tlsStage.ID)

	return ResourceTLSStageRead(ctx, d, m)
}

func ResourceTLSStageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := NewEdgeServicesAPI(m)

	tlsStage, err := api.GetTLSStage(&edgeservices.GetTLSStageRequest{
		TLSStageID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("backend_stage_id", types.FlattenStringPtr(tlsStage.BackendStageID))
	_ = d.Set("cache_stage_id", types.FlattenStringPtr(tlsStage.CacheStageID))
	_ = d.Set("pipeline_id", types.FlattenStringPtr(tlsStage.PipelineID))
	_ = d.Set("managed_certificate", tlsStage.ManagedCertificate)
	_ = d.Set("secrets", flattenTLSSecrets(tlsStage.Secrets))
	_ = d.Set("certificate_expires_at", types.FlattenTime(tlsStage.CertificateExpiresAt))
	_ = d.Set("created_at", types.FlattenTime(tlsStage.CreatedAt))
	_ = d.Set("updated_at", types.FlattenTime(tlsStage.UpdatedAt))
	_ = d.Set("project_id", tlsStage.ProjectID)

	return nil
}

func ResourceTLSStageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, region, err := NewEdgeServicesAPIWithRegion(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	hasChanged := false

	updateRequest := &edgeservices.UpdateTLSStageRequest{
		TLSStageID: d.Id(),
	}

	if d.HasChange("backend_stage_id") {
		updateRequest.BackendStageID = types.ExpandUpdatedStringPtr(d.Get("backend_stage_id"))
		hasChanged = true
	}

	if d.HasChange("cache_stage_id") {
		updateRequest.CacheStageID = types.ExpandUpdatedStringPtr(d.Get("cache_stage_id"))
		hasChanged = true
	}

	if d.HasChange("managed_certificate") {
		updateRequest.ManagedCertificate = types.ExpandBoolPtr(d.Get("managed_certificate"))
		hasChanged = true
	}

	if d.HasChange("secrets") {
		updateRequest.TLSSecretsConfig = wrapSecretsInConfig(expandTLSSecrets(d.Get("secrets"), region))
		hasChanged = true
	}

	if hasChanged {
		_, err = api.UpdateTLSStage(updateRequest, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceTLSStageRead(ctx, d, m)
}

func ResourceTLSStageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := NewEdgeServicesAPI(m)

	err := retry.RetryContext(ctx, defaultEdgeServicesTimeout, func() *retry.RetryError {
		err := api.DeleteTLSStage(&edgeservices.DeleteTLSStageRequest{
			TLSStageID: d.Id(),
		}, scw.WithContext(ctx))
		if err != nil && !httperrors.Is403(err) {
			if isStageUsedInPipelineError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
