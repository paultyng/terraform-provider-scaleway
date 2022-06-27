package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayIamGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayIamGroupCreate,
		ReadContext:   resourceScalewayIamGroupRead,
		UpdateContext: resourceScalewayIamGroupUpdate,
		DeleteContext: resourceScalewayIamGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the iam group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the iam group",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the group",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the group",
			},
			"user_ids": {
				Type:        schema.TypeList,
				Description: "List of IDs of the users attached to the group",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validationUUID(),
				},
			},
			"application_ids": {
				Type:        schema.TypeList,
				Description: "List of IDs of the applications attached to the group",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validationUUID(),
				},
			},
			"organization_id": organizationIDSchema(),
		},
	}
}

func resourceScalewayIamGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := iamAPI(meta)
	req := &iam.CreateGroupRequest{
		OrganizationID: d.Get("organization_id").(string),
		Name:           expandOrGenerateString(d.Get("name"), "group-"),
		Description:    d.Get("description").(string),
	}
	group, err := api.CreateGroup(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if appIds := d.Get("application_ids").([]interface{}); len(appIds) > 0 {
		appIdsStr := []string(nil)
		for _, id := range appIds {
			appIdsStr = append(appIdsStr, id.(string))
		}
		_, err := api.SetGroupPrincipals(&iam.SetGroupPrincipalsRequest{
			ApplicationIDs: appIdsStr,
			GroupID:        group.ID,
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if userIds := d.Get("user_ids").([]interface{}); len(userIds) > 0 {
		userIdsStr := []string(nil)
		for _, id := range userIds {
			userIdsStr = append(userIdsStr, id.(string))
		}
		_, err := api.SetGroupPrincipals(&iam.SetGroupPrincipalsRequest{
			UserIDs: userIdsStr,
			GroupID: group.ID,
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(group.ID)

	return resourceScalewayIamGroupRead(ctx, d, meta)
}

func resourceScalewayIamGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := iamAPI(meta)
	group, err := api.GetGroup(&iam.GetGroupRequest{
		GroupID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", group.Name)
	_ = d.Set("description", group.Description)
	_ = d.Set("created_at", flattenTime(group.CreatedAt))
	_ = d.Set("updated_at", flattenTime(group.UpdatedAt))
	_ = d.Set("organization_id", group.OrganizationID)
	_ = d.Set("user_ids", group.UserIDs)
	_ = d.Set("application_ids", group.ApplicationIDs)

	return nil
}

func resourceScalewayIamGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := iamAPI(meta)

	group, err := api.GetGroup(&iam.GetGroupRequest{
		GroupID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &iam.UpdateGroupRequest{
		GroupID: group.ID,
	}

	if d.HasChange("name") {
		req.Name = expandStringPtr(d.Get("name").(string))
	} else {
		req.Name = &group.Name
	}

	if d.HasChange("description") {
		req.Description = expandStringPtr(d.Get("description").(string))
	} else if group.Description != "" {
		req.Description = &group.Description
	} else {
		req.Description = nil
	}

	if d.HasChange("application_ids") {
		appIdsStr := []string(nil)
		if appIds := d.Get("application_ids").([]interface{}); len(appIds) > 0 {
			for _, id := range appIds {
				appIdsStr = append(appIdsStr, id.(string))
			}
			_, err = api.SetGroupPrincipals(&iam.SetGroupPrincipalsRequest{
				ApplicationIDs: appIdsStr,
				GroupID:        group.ID,
			}, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			for _, toRemove := range group.ApplicationIDs {
				_, err = api.DeletePrincipalFromGroup(&iam.DeletePrincipalFromGroupRequest{
					GroupID:     group.ID,
					PrincipalID: toRemove,
				}, scw.WithContext(ctx))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("user_ids") {
		userIdsStr := []string(nil)
		if userIds := d.Get("user_ids").([]interface{}); len(userIds) > 0 {
			for _, id := range userIds {
				userIdsStr = append(userIdsStr, id.(string))
			}
			_, err = api.SetGroupPrincipals(&iam.SetGroupPrincipalsRequest{
				UserIDs: userIdsStr,
				GroupID: group.ID,
			}, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			for _, toRemove := range group.UserIDs {
				_, err = api.DeletePrincipalFromGroup(&iam.DeletePrincipalFromGroupRequest{
					GroupID:     group.ID,
					PrincipalID: toRemove,
				}, scw.WithContext(ctx))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	_, err = api.UpdateGroup(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayIamGroupRead(ctx, d, meta)
}

func resourceScalewayIamGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := iamAPI(meta)

	err := api.DeleteGroup(&iam.DeleteGroupRequest{
		GroupID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
