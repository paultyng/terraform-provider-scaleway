package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayRdbACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayRdbACLCreate,
		ReadContext:   resourceScalewayRdbACLRead,
		UpdateContext: resourceScalewayRdbACLUpdate,
		DeleteContext: resourceScalewayRdbACLDelete,
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
				Description:  "Instance on which the ACL is applied",
			},
			"acl_rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of ACL rules to apply",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:         schema.TypeString,
							ValidateFunc: validation.IsCIDR,
							Required:     true,
							Description:  "Target IP of the rules",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the rule",
						},
					},
				},
			},
			// Common
			"region": regionSchema(),
		},
	}
}

func resourceScalewayRdbACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceID := d.Get("instance_id").(string)
	rdbAPI, region, ID, err := rdbAPIWithRegionAndID(meta, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	//InstanceStatus.READY,
	//	InstanceStatus.CONFIGURING,
	//	InstanceStatus.BACKUPING,
	//	InstanceStatus.SNAPSHOTTING,
	_ = rdb.WaitForInstanceRequest{
		InstanceID: ID,
		Region:     region,
	}

	createReq := &rdb.SetInstanceACLRulesRequest{
		Region:     region,
		InstanceID: ID,
		Rules:      rdbACLExpand(d.Get("acl_rules")),
	}

	_, err = rdbAPI.SetInstanceACLRules(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID)

	return resourceScalewayRdbACLRead(ctx, d, meta)
}

func resourceScalewayRdbACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := rdbAPIWithRegionAndID(meta, d.Get("instance_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := rdbAPI.ListInstanceACLRules(&rdb.ListInstanceACLRulesRequest{
		Region:     region,
		InstanceID: ID,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	id := newRegionalID(region, ID).String()
	d.SetId(id)
	_ = d.Set("instance_id", id)
	_ = d.Set("acl_rules", rdbACLRulesFlatten(res.Rules))

	return nil
}

func resourceScalewayRdbACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := rdbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("acl_rules") {
		//InstanceStatus.READY,
		//	InstanceStatus.CONFIGURING,
		//	InstanceStatus.BACKUPING,
		//	InstanceStatus.SNAPSHOTTING,
		_ = rdb.WaitForInstanceRequest{
			InstanceID: ID,
			Region:     region,
		}

		req := &rdb.SetInstanceACLRulesRequest{
			Region:     region,
			InstanceID: ID,
			Rules:      rdbACLExpand(d.Get("acl_rules")),
		}

		_, err = rdbAPI.SetInstanceACLRules(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayRdbACLRead(ctx, d, meta)
}

func resourceScalewayRdbACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rdbAPI, region, ID, err := rdbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	aclruleips := make([]string, 0)
	for _, acl := range rdbACLExpand(d.Get("acl_rules")) {
		aclruleips = append(aclruleips, acl.IP.String())
	}
	_, err = rdbAPI.DeleteInstanceACLRules(&rdb.DeleteInstanceACLRulesRequest{
		Region:     region,
		InstanceID: ID,
		ACLRuleIPs: aclruleips,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}

func rdbACLExpand(data interface{}) []*rdb.ACLRuleRequest {
	type aclRule struct {
		IP          string
		Description string
	}
	var res []*rdb.ACLRuleRequest
	for _, rule := range data.([]interface{}) {
		r := rule.(map[string]interface{})
		res = append(res, &rdb.ACLRuleRequest{
			IP:          expandIPNet(r["ip"].(string)),
			Description: r["description"].(string),
		})
	}

	return res
}

func rdbACLRulesFlatten(rules []*rdb.ACLRule) []map[string]interface{} {
	var res []map[string]interface{}
	for _, rule := range rules {
		r := map[string]interface{}{
			"ip":          rule.IP.String(),
			"description": rule.Description,
		}
		res = append(res, r)
	}
	return res
}
