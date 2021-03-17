package scaleway

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayInstanceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayInstanceSecurityGroupCreate,
		ReadContext:   resourceScalewayInstanceSecurityGroupRead,
		UpdateContext: resourceScalewayInstanceSecurityGroupUpdate,
		DeleteContext: resourceScalewayInstanceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultInstanceSecurityGroupTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the security group",
			},
			"stateful": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The stateful value of the security group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the security group",
			},
			"inbound_default_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "accept",
				Description: "Default inbound traffic policy for this security group",
				ValidateFunc: validation.StringInSlice([]string{
					instance.SecurityGroupPolicyAccept.String(),
					instance.SecurityGroupPolicyDrop.String(),
				}, false),
			},
			"outbound_default_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "accept",
				Description: "Default outbound traffic policy for this security group",
				ValidateFunc: validation.StringInSlice([]string{
					instance.SecurityGroupPolicyAccept.String(),
					instance.SecurityGroupPolicyDrop.String(),
				}, false),
			},
			"inbound_rule": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Inbound rules for this security group",
				Elem:          securityGroupRuleSchema(),
				ConflictsWith: []string{"external_rules"},
			},
			"outbound_rule": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Outbound rules for this security group",
				Elem:          securityGroupRuleSchema(),
				ConflictsWith: []string{"external_rules"},
			},
			"external_rules": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"inbound_rule", "outbound_rule"},
			},
			"enable_default_security": {
				Type:        schema.TypeBool,
				Description: "Enable blocking of SMTP on IPv4 and IPv6",
				Optional:    true,
				Default:     true,
			},
			"zone":            zoneSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
		},
	}
}

func resourceScalewayInstanceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &instance.CreateSecurityGroupRequest{
		Name:                  expandOrGenerateString(d.Get("name"), "sg"),
		Zone:                  zone,
		Project:               expandStringPtr(d.Get("project_id")),
		Description:           d.Get("description").(string),
		Stateful:              d.Get("stateful").(bool),
		InboundDefaultPolicy:  instance.SecurityGroupPolicy(d.Get("inbound_default_policy").(string)),
		OutboundDefaultPolicy: instance.SecurityGroupPolicy(d.Get("outbound_default_policy").(string)),
		EnableDefaultSecurity: expandBoolPtr(d.Get("enable_default_security")),
	}

	if definedZone, ok := d.GetOk("zone"); ok {
		zone = scw.Zone(definedZone.(string))
		req.Zone = zone
	}

	res, err := instanceAPI.CreateSecurityGroup(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.SecurityGroup.ID))

	if d.Get("external_rules").(bool) {
		return resourceScalewayInstanceSecurityGroupRead(ctx, d, meta)
	}
	// We call update instead of read as it will take care of creating rules.
	return resourceScalewayInstanceSecurityGroupUpdate(ctx, d, meta)
}

func resourceScalewayInstanceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := instanceAPI.GetSecurityGroup(&instance.GetSecurityGroupRequest{
		SecurityGroupID: ID,
		Zone:            zone,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("zone", zone)
	_ = d.Set("organization_id", res.SecurityGroup.Organization)
	_ = d.Set("project_id", res.SecurityGroup.Project)
	_ = d.Set("name", res.SecurityGroup.Name)
	_ = d.Set("stateful", res.SecurityGroup.Stateful)
	_ = d.Set("description", res.SecurityGroup.Description)
	_ = d.Set("inbound_default_policy", res.SecurityGroup.InboundDefaultPolicy.String())
	_ = d.Set("outbound_default_policy", res.SecurityGroup.OutboundDefaultPolicy.String())
	_ = d.Set("enable_default_security", res.SecurityGroup.EnableDefaultSecurity)

	if !d.Get("external_rules").(bool) {
		inboundRules, outboundRules, err := getSecurityGroupRules(ctx, instanceAPI, zone, ID, d)
		if err != nil {
			return diag.FromErr(err)
		}
		_ = d.Set("inbound_rule", inboundRules)
		_ = d.Set("outbound_rule", outboundRules)
	}
	return nil
}

func getSecurityGroupRules(ctx context.Context, instanceAPI *instance.API, zone scw.Zone, securityGroupID string, d *schema.ResourceData) ([]interface{}, []interface{}, error) {
	resRules, err := instanceAPI.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
		Zone:            zone,
		SecurityGroupID: expandID(securityGroupID),
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	sort.Slice(resRules.Rules, func(i, j int) bool {
		return resRules.Rules[i].Position < resRules.Rules[j].Position
	})
	apiRules := map[instance.SecurityGroupRuleDirection][]*instance.SecurityGroupRule{
		instance.SecurityGroupRuleDirectionInbound:  {},
		instance.SecurityGroupRuleDirectionOutbound: {},
	}

	stateRules := map[instance.SecurityGroupRuleDirection][]interface{}{
		instance.SecurityGroupRuleDirectionInbound:  d.Get("inbound_rule").([]interface{}),
		instance.SecurityGroupRuleDirectionOutbound: d.Get("outbound_rule").([]interface{}),
	}

	for _, apiRule := range resRules.Rules {
		if !apiRule.Editable {
			continue
		}
		apiRules[apiRule.Direction] = append(apiRules[apiRule.Direction], apiRule)
	}

	// We make sure that we keep state rule if they match their api rule.
	for direction := range apiRules {
		for index, apiRule := range apiRules[direction] {
			if index < len(stateRules[direction]) {
				stateRule := securityGroupRuleExpand(stateRules[direction][index])
				if !securityGroupRuleEquals(stateRule, apiRule) {
					stateRules[direction][index] = securityGroupRuleFlatten(apiRule)
				}
			} else {
				stateRules[direction] = append(stateRules[direction], securityGroupRuleFlatten(apiRule))
			}
		}
	}

	return stateRules[instance.SecurityGroupRuleDirectionInbound], stateRules[instance.SecurityGroupRuleDirectionOutbound], nil
}

func resourceScalewayInstanceSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, _, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	zone, ID, err := parseZonedID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	inboundDefaultPolicy := instance.SecurityGroupPolicy("")
	if d.Get("inbound_default_policy") != nil {
		inboundDefaultPolicy = instance.SecurityGroupPolicy(d.Get("inbound_default_policy").(string))
	}
	outboundDefaultPolicy := instance.SecurityGroupPolicy("")
	if d.Get("outbound_default_policy") != nil {
		outboundDefaultPolicy = instance.SecurityGroupPolicy(d.Get("outbound_default_policy").(string))
	}

	description := ""
	if d.Get("description") != nil {
		description = d.Get("description").(string)
	}
	updateReq := &instance.UpdateSecurityGroupRequest{
		Zone:                  zone,
		SecurityGroupID:       ID,
		Stateful:              scw.BoolPtr(d.Get("stateful").(bool)),
		Description:           expandStringPtr(description),
		InboundDefaultPolicy:  &inboundDefaultPolicy,
		OutboundDefaultPolicy: &outboundDefaultPolicy,
	}

	if d.HasChange("enable_default_security") {
		updateReq.EnableDefaultSecurity = expandBoolPtr(d.Get("enable_default_security"))
	}

	// Only update name if one is provided in the state
	if d.Get("name") != nil && d.Get("name").(string) != "" {
		updateReq.Name = expandStringPtr(d.Get("name"))
	}

	_, err = instanceAPI.UpdateSecurityGroup(updateReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.Get("external_rules").(bool) {
		err = updateSecurityGroupeRules(ctx, d, zone, ID, instanceAPI)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayInstanceSecurityGroupRead(ctx, d, meta)
}

// updateSecurityGroupeRules handles updating SecurityGroupRules
//
// It works as followed:
//   1) Creates 2 map[direction][]rule: one for rules in state and one for rules in API
//   2) For each direction we:
//     A) Loop for each rule in state for this direction
//       a) Compare with api rule in this direction at the same index
//          if different update / if equals do nothing / if no more api rules to compare create new api rule
//     B) If there is more rule in the API we remove them
func updateSecurityGroupeRules(ctx context.Context, d *schema.ResourceData, zone scw.Zone, securityGroupID string, instanceAPI *instance.API) error {
	apiRules := map[instance.SecurityGroupRuleDirection][]*instance.SecurityGroupRule{
		instance.SecurityGroupRuleDirectionInbound:  {},
		instance.SecurityGroupRuleDirectionOutbound: {},
	}
	stateRules := map[instance.SecurityGroupRuleDirection][]interface{}{
		instance.SecurityGroupRuleDirectionInbound:  d.Get("inbound_rule").([]interface{}),
		instance.SecurityGroupRuleDirectionOutbound: d.Get("outbound_rule").([]interface{}),
	}

	// Fill apiRules with data from API
	resRules, err := instanceAPI.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
		Zone:            zone,
		SecurityGroupID: expandID(securityGroupID),
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return err
	}
	sort.Slice(resRules.Rules, func(i, j int) bool {
		return resRules.Rules[i].Position < resRules.Rules[j].Position
	})
	for _, apiRule := range resRules.Rules {
		if !apiRule.Editable {
			continue
		}
		apiRules[apiRule.Direction] = append(apiRules[apiRule.Direction], apiRule)
	}

	// Loop through all directions
	for direction := range stateRules {
		// Loop for all state rules in this direction
		for index, rawStateRule := range stateRules[direction] {
			apiRule := (*instance.SecurityGroupRule)(nil)
			stateRule := securityGroupRuleExpand(rawStateRule)

			// This happen when there is more rule in state than in the api. We create more rule in API.
			if index >= len(apiRules[direction]) {
				_, err = instanceAPI.CreateSecurityGroupRule(&instance.CreateSecurityGroupRuleRequest{
					Zone:            zone,
					SecurityGroupID: securityGroupID,
					Protocol:        stateRule.Protocol,
					IPRange:         stateRule.IPRange,
					Action:          stateRule.Action,
					DestPortTo:      stateRule.DestPortTo,
					DestPortFrom:    stateRule.DestPortFrom,
					Direction:       direction,
				}, scw.WithContext(ctx))
				if err != nil {
					return err
				}
				continue
			}

			// We compare rule stateRule[index] and apiRule[index]. If they are different we update api rule to match state.
			apiRule = apiRules[direction][index]
			if !securityGroupRuleEquals(stateRule, apiRule) {
				destPortFrom := stateRule.DestPortFrom
				destPortTo := stateRule.DestPortTo
				if destPortFrom == nil {
					destPortFrom = scw.Uint32Ptr(0)
				}
				if destPortTo == nil {
					destPortTo = scw.Uint32Ptr(0)
				}

				_, err = instanceAPI.UpdateSecurityGroupRule(&instance.UpdateSecurityGroupRuleRequest{
					Zone:                zone,
					SecurityGroupID:     securityGroupID,
					SecurityGroupRuleID: apiRule.ID,
					Protocol:            &stateRule.Protocol,
					IPRange:             &stateRule.IPRange,
					Action:              &stateRule.Action,
					DestPortTo:          destPortTo,
					DestPortFrom:        destPortFrom,
					Direction:           &direction,
				}, scw.WithContext(ctx))
				if err != nil {
					return err
				}
			}
		}

		// We loop through remaining API rules and delete them as they are no longer in the state.
		for index := len(stateRules[direction]); index < len(apiRules[direction]); index++ {
			err = instanceAPI.DeleteSecurityGroupRule(&instance.DeleteSecurityGroupRuleRequest{
				Zone:                zone,
				SecurityGroupID:     securityGroupID,
				SecurityGroupRuleID: apiRules[direction][index].ID,
			}, scw.WithContext(ctx))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceScalewayInstanceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, _, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	zone, ID, err := parseZonedID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = instanceAPI.DeleteSecurityGroup(&instance.DeleteSecurityGroupRequest{
		SecurityGroupID: ID,
		Zone:            zone,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}

// securityGroupRuleSchema returns schema for inbound/outbound rule in security group
func securityGroupRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					instance.SecurityGroupRuleActionAccept.String(),
					instance.SecurityGroupRuleActionDrop.String(),
				}, false),
				Description: "Action when rule match request (drop or accept)",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  instance.SecurityGroupRuleProtocolTCP.String(),
				ValidateFunc: validation.StringInSlice([]string{
					instance.SecurityGroupRuleProtocolICMP.String(),
					instance.SecurityGroupRuleProtocolTCP.String(),
					instance.SecurityGroupRuleProtocolUDP.String(),
					instance.SecurityGroupRuleProtocolANY.String(),
				}, false),
				Description: "Protocol for this rule (TCP, UDP, ICMP or ANY)",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Network port for this rule",
			},
			"port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Computed port range for this rule (e.g: 1-1024, 22-22)",
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Ip address for this rule (e.g: 1.1.1.1). Only one of ip or ip_range should be provided",
			},
			"ip_range": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 128),
				Description:  "Ip range for this rule (e.g: 192.168.1.0/24). Only one of ip or ip_range should be provided",
			},
		},
	}
}

// securityGroupRuleExpand transform a state rule to an api one.
func securityGroupRuleExpand(i interface{}) *instance.SecurityGroupRule {
	rawRule := i.(map[string]interface{})

	portFrom, portTo := uint32(0), uint32(0)

	portRange := rawRule["port_range"].(string)
	if portRange != "" {
		_, _ = fmt.Sscanf(portRange, "%d-%d", &portFrom, &portTo)
	} else {
		portFrom = uint32(rawRule["port"].(int))
		portTo = portFrom
	}

	action, _ := rawRule["action"].(string)
	ipRange := rawRule["ip_range"].(string)
	if ipRange == "" {
		ipRange = rawRule["ip"].(string) + "/32"
	}
	if ipRange == "/32" {
		ipRange = "0.0.0.0/0"
	}

	rule := &instance.SecurityGroupRule{
		DestPortFrom: &portFrom,
		DestPortTo:   &portTo,
		Protocol:     instance.SecurityGroupRuleProtocol(rawRule["protocol"].(string)),
		IPRange:      expandIPNet(ipRange),
		Action:       instance.SecurityGroupRuleAction(action),
	}

	if *rule.DestPortFrom == *rule.DestPortTo {
		rule.DestPortTo = nil
	}

	// Handle when no port is specified.
	if portFrom == 0 && portTo == 0 {
		rule.DestPortFrom = nil
		rule.DestPortTo = nil
	}

	return rule
}

// securityGroupRuleFlatten transform a api rule to an state one.
func securityGroupRuleFlatten(rule *instance.SecurityGroupRule) map[string]interface{} {
	portFrom, portTo := uint32(0), uint32(0)

	if rule.DestPortFrom != nil {
		portFrom = *rule.DestPortFrom
	}

	if rule.DestPortTo != nil {
		portTo = *rule.DestPortTo
	}

	res := map[string]interface{}{
		"protocol":   rule.Protocol.String(),
		"ip_range":   flattenIPNet(rule.IPRange),
		"port_range": fmt.Sprintf("%d-%d", portFrom, portTo),
		"action":     rule.Action.String(),
	}
	return res
}

// securityGroupRuleEquals compares two security group rule.
func securityGroupRuleEquals(ruleA, ruleB *instance.SecurityGroupRule) bool {
	zeroIfNil := func(v *uint32) uint32 {
		if v == nil {
			return 0
		}
		return *v
	}
	portFromEqual := zeroIfNil(ruleA.DestPortFrom) == zeroIfNil(ruleB.DestPortFrom)
	portToEqual := zeroIfNil(ruleA.DestPortTo) == zeroIfNil(ruleB.DestPortTo)
	ipEqual := flattenIPNet(ruleA.IPRange) == flattenIPNet(ruleB.IPRange)

	return ruleA.Action == ruleB.Action &&
		portFromEqual &&
		portToEqual &&
		ipEqual &&
		ruleA.Protocol == ruleB.Protocol
}
