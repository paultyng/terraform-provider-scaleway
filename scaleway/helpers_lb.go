package scaleway

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	lb "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	LbWaitForTimeout = 10 * time.Minute
)

// lbAPI returns a new lb API
func lbAPI(m interface{}) *lb.API {
	meta := m.(*Meta)
	return lb.NewAPI(meta.scwClient)
}

// lbAPIWithRegion returns a new lb API and the region for a Create request
func lbAPIWithRegion(d *schema.ResourceData, m interface{}) (*lb.API, scw.Region, error) {
	meta := m.(*Meta)
	lbApi := lb.NewAPI(meta.scwClient)

	region, err := extractRegion(d, meta)
	return lbApi, region, err
}

// lbAPIWithRegionAndID returns an lb API with region and ID extracted from the state
func lbAPIWithRegionAndID(m interface{}, id string) (*lb.API, scw.Region, string, error) {
	meta := m.(*Meta)
	lbApi := lb.NewAPI(meta.scwClient)

	region, ID, err := parseRegionalID(id)
	return lbApi, region, ID, err
}

func flattenLbBackendMarkdownAction(action lb.OnMarkedDownAction) interface{} {
	if action == lb.OnMarkedDownActionOnMarkedDownActionNone {
		return "none"
	}
	return action.String()
}

func flattenLbAcl(acl *lb.ACL) interface{} {
	res := map[string]interface{}{
		"name":   acl.Name,
		"match":  flattenLbAclMatch(acl.Match),
		"action": flattenLbAclAction(acl.Action),
	}
	return res
}

// expandLbAcl transforms a state acl to an api one.
func expandLbAcl(i interface{}) *lb.ACL {
	rawRule := i.(map[string]interface{})
	acl := &lb.ACL{
		Name:   rawRule["name"].(string),
		Match:  expandLbAclMatch(rawRule["match"]),
		Action: expandLbAclAction(rawRule["action"]),
	}

	//remove http filter values if we do not pass any http filter
	if acl.Match.HTTPFilter == "" || acl.Match.HTTPFilter == lb.ACLHTTPFilterACLHTTPFilterNone {
		acl.Match.HTTPFilter = lb.ACLHTTPFilterACLHTTPFilterNone
		acl.Match.HTTPFilterValue = []*string{}
	}

	return acl
}
func flattenLbAclAction(action *lb.ACLAction) interface{} {
	return map[string]interface{}{
		"type": action.Type,
	}
}
func expandLbAclAction(raw interface{}) *lb.ACLAction {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lb.ACLAction{
		Type: lb.ACLActionType(rawMap["type"].(string)),
	}
}
func flattenLbAclMatch(match *lb.ACLMatch) interface{} {
	return map[string]interface{}{
		"ip_subnet":         flattenSliceStringPtr(match.IPSubnet),
		"http_filter":       match.HTTPFilter.String(),
		"http_filter_value": flattenSliceStringPtr(match.HTTPFilterValue),
		"invert":            match.Invert,
	}
}
func expandLbAclMatch(raw interface{}) *lb.ACLMatch {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})

	//scaleway api require ip subnet, so if we did not specify one, just put 0.0.0.0/0 instead
	ipSubnet := expandSliceStringPtr(rawMap["ip_subnet"].([]interface{}))
	if len(ipSubnet) == 0 {
		ipSubnet = []*string{scw.StringPtr("0.0.0.0/0")}
	}

	return &lb.ACLMatch{
		IPSubnet:        ipSubnet,
		HTTPFilter:      lb.ACLHTTPFilter(rawMap["http_filter"].(string)),
		HTTPFilterValue: expandSliceStringPtr(rawMap["http_filter_value"].([]interface{})),
		Invert:          rawMap["invert"].(bool),
	}
}

func expandLbBackendMarkdownAction(raw interface{}) lb.OnMarkedDownAction {
	if raw == "none" {
		return lb.OnMarkedDownActionOnMarkedDownActionNone
	}
	return lb.OnMarkedDownAction(raw.(string))
}

func flattenLbProtocol(protocol lb.Protocol) interface{} {
	return protocol.String()
}

func expandLbProtocol(raw interface{}) lb.Protocol {
	return lb.Protocol(raw.(string))
}

func flattenLbForwardPortAlgorithm(algo lb.ForwardPortAlgorithm) interface{} {
	return algo.String()
}

func expandLbForwardPortAlgorithm(raw interface{}) lb.ForwardPortAlgorithm {
	return lb.ForwardPortAlgorithm(raw.(string))
}

func flattenLbStickySessionsType(t lb.StickySessionsType) interface{} {
	return t.String()
}

func expandLbStickySessionsType(raw interface{}) lb.StickySessionsType {
	return lb.StickySessionsType(raw.(string))
}

func flattenLbHCTCP(config *lb.HealthCheckTCPConfig) interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{},
	}
}

func expandLbHCTCP(raw interface{}) *lb.HealthCheckTCPConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	return &lb.HealthCheckTCPConfig{}
}

func flattenLbHCHTTP(config *lb.HealthCheckHTTPConfig) interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"uri":    config.URI,
			"method": config.Method,
			"code":   flattenInt32Ptr(config.Code),
		},
	}
}

func expandLbHCHTTP(raw interface{}) *lb.HealthCheckHTTPConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lb.HealthCheckHTTPConfig{
		URI:    rawMap["uri"].(string),
		Method: rawMap["method"].(string),
		Code:   expandInt32Ptr(rawMap["code"]),
	}
}

func flattenLbHCHTTPS(config *lb.HealthCheckHTTPSConfig) interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"uri":    config.URI,
			"method": config.Method,
			"code":   flattenInt32Ptr(config.Code),
		},
	}
}

func expandLbHCHTTPS(raw interface{}) *lb.HealthCheckHTTPSConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lb.HealthCheckHTTPSConfig{
		URI:    rawMap["uri"].(string),
		Method: rawMap["method"].(string),
		Code:   expandInt32Ptr(rawMap["code"]),
	}
}

func expandLbLetsEncrypt(raw interface{}) *lb.CreateCertificateRequestLetsencryptConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	alternativeNames := rawMap["subject_alternative_name"].([]interface{})
	config := &lb.CreateCertificateRequestLetsencryptConfig{
		CommonName: rawMap["common_name"].(string),
	}
	for _, alternativeName := range alternativeNames {
		config.SubjectAlternativeName = append(config.SubjectAlternativeName, alternativeName.(string))
	}
	return config
}

func expandLbCustomCertificate(raw interface{}) *lb.CreateCertificateRequestCustomCertificate {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	config := &lb.CreateCertificateRequestCustomCertificate{
		CertificateChain: rawMap["certificate_chain"].(string),
	}
	return config
}

func expandLbSubscriberEmailConfig(raw interface{}) *lb.SubscriberEmailConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	config := &lb.SubscriberEmailConfig{
		Email: rawMap["email"].(string),
	}
	return config
}

func expandLbSubscriberWebhookConfig(raw interface{}) *lb.SubscriberWebhookConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	config := &lb.SubscriberWebhookConfig{
		URI: rawMap["uri"].(string),
	}
	return config
}
