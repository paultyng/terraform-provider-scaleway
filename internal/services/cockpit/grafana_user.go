package cockpit

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceCockpitGrafanaUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCockpitGrafanaUserCreate,
		ReadContext:   ResourceCockpitGrafanaUserRead,
		DeleteContext: ResourceCockpitGrafanaUserDelete,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(DefaultCockpitTimeout),
			Read:    schema.DefaultTimeout(DefaultCockpitTimeout),
			Delete:  schema.DefaultTimeout(DefaultCockpitTimeout),
			Default: schema.DefaultTimeout(DefaultCockpitTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"login": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The login of the Grafana user",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_-]{2,24}$`), "must have between 2 and 24 alphanumeric characters"),
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password of the Grafana user",
				Sensitive:   true,
			},
			"role": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The role of the Grafana user",
				ValidateDiagFunc: verify.ValidateEnum[cockpit.GrafanaUserRole](),
			},
			"grafana_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The grafana URL",
			},

			"project_id": account.ProjectIDSchema(),
		},
	}
}

func ResourceCockpitGrafanaUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, err := NewGlobalAPI(m)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get("project_id").(string)
	login := d.Get("login").(string)
	role := cockpit.GrafanaUserRole(d.Get("role").(string))

	grafanaUser, err := api.CreateGrafanaUser(&cockpit.GlobalAPICreateGrafanaUserRequest{
		ProjectID: projectID,
		Login:     login,
		Role:      role,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("password", grafanaUser.Password)
	d.SetId(cockpitIDWithProjectID(projectID, strconv.FormatUint(uint64(grafanaUser.ID), 10)))
	return ResourceCockpitGrafanaUserRead(ctx, d, m)
}

func ResourceCockpitGrafanaUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, projectID, grafanaUserID, err := NewAPIGrafanaUserID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := api.ListGrafanaUsers(&cockpit.GlobalAPIListGrafanaUsersRequest{
		ProjectID: projectID,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var grafanaUser *cockpit.GrafanaUser
	for _, user := range res.GrafanaUsers {
		if user.ID == grafanaUserID {
			grafanaUser = user
			break
		}
	}

	if grafanaUser == nil {
		d.SetId("")
		return nil
	}

	grafana, err := api.GetGrafana(&cockpit.GlobalAPIGetGrafanaRequest{
		ProjectID: projectID,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("login", grafanaUser.Login)
	_ = d.Set("role", grafanaUser.Role)
	_ = d.Set("project_id", projectID)
	_ = d.Set("grafana_url", grafana.GrafanaURL)

	return nil
}

func ResourceCockpitGrafanaUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, projectID, grafanaUserID, err := NewAPIGrafanaUserID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = api.DeleteGrafanaUser(&cockpit.GlobalAPIDeleteGrafanaUserRequest{
		ProjectID:     projectID,
		GrafanaUserID: grafanaUserID,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
