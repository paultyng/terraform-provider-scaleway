package scaleway

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	validator "github.com/scaleway/scaleway-sdk-go/validation"
)

const (
	defaultLbLbTimeout = 10 * time.Minute
	retryLbIPInterval  = 5 * time.Second
)

// lbAPIWithZone returns an lb API WITH zone for a Create request
func lbAPIWithZone(d *schema.ResourceData, m interface{}) (*lbSDK.ZonedAPI, scw.Zone, error) {
	meta := m.(*Meta)
	lbAPI := lbSDK.NewZonedAPI(meta.scwClient)

	zone, err := extractZone(d, meta)
	if err != nil {
		return nil, "", err
	}
	return lbAPI, zone, nil
}

// lbAPIWithZoneAndID returns an lb API with zone and ID extracted from the state
func lbAPIWithZoneAndID(m interface{}, id string) (*lbSDK.ZonedAPI, scw.Zone, string, error) {
	meta := m.(*Meta)
	lbAPI := lbSDK.NewZonedAPI(meta.scwClient)

	zone, ID, err := parseZonedID(id)
	if err != nil {
		return nil, "", "", err
	}
	return lbAPI, zone, ID, nil
}

func flattenLbBackendMarkdownAction(action lbSDK.OnMarkedDownAction) interface{} {
	if action == lbSDK.OnMarkedDownActionOnMarkedDownActionNone {
		return "none"
	}
	return action.String()
}

func flattenLbACL(acl *lbSDK.ACL) interface{} {
	res := map[string]interface{}{
		"name":   acl.Name,
		"match":  flattenLbACLMatch(acl.Match),
		"action": flattenLbACLAction(acl.Action),
	}
	return res
}

// expandLbACL transforms a state acl to an api one.
func expandLbACL(i interface{}) *lbSDK.ACL {
	rawRule := i.(map[string]interface{})
	acl := &lbSDK.ACL{
		Name:        rawRule["name"].(string),
		Description: rawRule["description"].(string),
		Match:       expandLbACLMatch(rawRule["match"]),
		Action:      expandLbACLAction(rawRule["action"]),
		CreatedAt:   expandTimePtr(rawRule["created_at"]),
		UpdatedAt:   expandTimePtr(rawRule["updated_at"]),
	}

	if rawRule["index"] != nil {
		acl.Index = int32(rawRule["index"].(int))
	}
	// remove http filter values if we do not pass any http filter
	if acl.Match.HTTPFilter == "" || acl.Match.HTTPFilter == lbSDK.ACLHTTPFilterACLHTTPFilterNone {
		acl.Match.HTTPFilter = lbSDK.ACLHTTPFilterACLHTTPFilterNone
		acl.Match.HTTPFilterValue = []*string{}
	}

	return acl
}

func flattenLbACLAction(action *lbSDK.ACLAction) interface{} {
	return []map[string]interface{}{
		{
			"type":     action.Type,
			"redirect": flattenLbACLActionRedirect(action.Redirect),
		},
	}
}

func expandLbACLAction(raw interface{}) *lbSDK.ACLAction {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lbSDK.ACLAction{
		Type:     lbSDK.ACLActionType(rawMap["type"].(string)),
		Redirect: expandLbACLActionRedirect(rawMap["redirect"]),
	}
}

func flattenLbACLActionRedirect(redirect *lbSDK.ACLActionRedirect) interface{} {
	if redirect == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"type":   redirect.Type,
			"target": redirect.Target,
			"code":   redirect.Code,
		},
	}
}

func expandLbACLActionRedirect(raw interface{}) *lbSDK.ACLActionRedirect {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lbSDK.ACLActionRedirect{
		Type:   lbSDK.ACLActionRedirectRedirectType(rawMap["type"].(string)),
		Target: rawMap["target"].(string),
		Code:   expandInt32Ptr(rawMap["code"]),
	}
}

func expandPrivateNetworks(data interface{}) ([]*lbSDK.PrivateNetwork, error) {
	if data == nil {
		return nil, nil
	}

	pns := []*lbSDK.PrivateNetwork(nil)

	for _, pn := range data.(*schema.Set).List() {
		rawPn := pn.(map[string]interface{})
		privateNetwork := &lbSDK.PrivateNetwork{}
		privateNetwork.PrivateNetworkID = expandID(rawPn["private_network_id"].(string))
		if staticConfig, hasStaticConfig := rawPn["static_config"]; hasStaticConfig {
			privateNetwork.StaticConfig = expandLbPrivateNetworkStaticConfig(staticConfig)
		} else {
			privateNetwork.DHCPConfig = expandLbPrivateNetworkDHCPConfig(rawPn["dhcp_config"])
		}

		pns = append(pns, privateNetwork)
	}
	return pns, nil
}

func isPrivateNetworkEqual(a, b interface{}) bool {
	// Find out the diff Private Network or not
	if _, ok := a.(*lbSDK.PrivateNetwork); ok {
		if _, ok := b.(*lbSDK.PrivateNetwork); ok {
			if a.(*lbSDK.PrivateNetwork).PrivateNetworkID == b.(*lbSDK.PrivateNetwork).PrivateNetworkID {
				// if both has dhcp config should not update
				if a.(*lbSDK.PrivateNetwork).DHCPConfig != nil && b.(*lbSDK.PrivateNetwork).DHCPConfig != nil {
					return true
				}
				// check static config
				aConfig := a.(*lbSDK.PrivateNetwork).StaticConfig
				bConfig := b.(*lbSDK.PrivateNetwork).StaticConfig
				if aConfig != nil && bConfig != nil {
					// check if static config is different
					return reflect.DeepEqual(aConfig.IPAddress, bConfig.IPAddress)
				}
			}
		}
	}
	return false
}

func privateNetworksCompare(slice1, slice2 []*lbSDK.PrivateNetwork) []*lbSDK.PrivateNetwork {
	var diff []*lbSDK.PrivateNetwork

	m := make(map[string]struct{}, len(slice1))
	for _, pn := range slice1 {
		m[pn.PrivateNetworkID] = struct{}{}
	}
	// find the differences
	for _, pn := range slice2 {
		if _, foundID := m[pn.PrivateNetworkID]; !foundID || (foundID && !isPrivateNetworkEqual(slice1, slice2)) {
			diff = append(diff, pn)
		}
	}
	return diff
}

func flattenPrivateNetworkConfigs(privateNetworks []*lbSDK.PrivateNetwork) interface{} {
	if len(privateNetworks) == 0 || privateNetworks == nil {
		return nil
	}

	pnI := []map[string]interface{}(nil)
	var dhcpConfigExist bool
	for _, pn := range privateNetworks {
		if pn.DHCPConfig != nil {
			dhcpConfigExist = true
		}
		pnZonedID := newZonedIDString(pn.LB.Zone, pn.PrivateNetworkID)
		pnI = append(pnI, map[string]interface{}{
			"private_network_id": pnZonedID,
			"dhcp_config":        dhcpConfigExist,
			"status":             pn.Status.String(),
			"zone":               pn.LB.Zone.String(),
			"static_config":      flattenLbPrivateNetworkStaticConfig(pn.StaticConfig),
		})
	}
	return pnI
}

func flattenLbACLMatch(match *lbSDK.ACLMatch) interface{} {
	return []map[string]interface{}{
		{
			"ip_subnet":          flattenSliceStringPtr(match.IPSubnet),
			"http_filter":        match.HTTPFilter.String(),
			"http_filter_value":  flattenSliceStringPtr(match.HTTPFilterValue),
			"http_filter_option": match.HTTPFilterOption,
			"invert":             match.Invert,
		},
	}
}

func expandLbACLMatch(raw interface{}) *lbSDK.ACLMatch {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})

	// scaleway api require ip subnet, so if we did not specify one, just put 0.0.0.0/0 instead
	ipSubnet := expandSliceStringPtr(rawMap["ip_subnet"].([]interface{}))
	if len(ipSubnet) == 0 {
		ipSubnet = []*string{expandStringPtr("0.0.0.0/0")}
	}

	return &lbSDK.ACLMatch{
		IPSubnet:         ipSubnet,
		HTTPFilter:       lbSDK.ACLHTTPFilter(rawMap["http_filter"].(string)),
		HTTPFilterValue:  expandSliceStringPtr(rawMap["http_filter_value"].([]interface{})),
		HTTPFilterOption: expandStringPtr(rawMap["http_filter_option"].(string)),
		Invert:           rawMap["invert"].(bool),
	}
}

func expandLbBackendMarkdownAction(raw interface{}) lbSDK.OnMarkedDownAction {
	if raw == "none" {
		return lbSDK.OnMarkedDownActionOnMarkedDownActionNone
	}
	return lbSDK.OnMarkedDownAction(raw.(string))
}

func flattenLbProtocol(protocol lbSDK.Protocol) interface{} {
	return protocol.String()
}

func expandLbProtocol(raw interface{}) lbSDK.Protocol {
	return lbSDK.Protocol(raw.(string))
}

func flattenLbForwardPortAlgorithm(algo lbSDK.ForwardPortAlgorithm) interface{} {
	return algo.String()
}

func expandLbForwardPortAlgorithm(raw interface{}) lbSDK.ForwardPortAlgorithm {
	return lbSDK.ForwardPortAlgorithm(raw.(string))
}

func flattenLbStickySessionsType(t lbSDK.StickySessionsType) interface{} {
	return t.String()
}

func expandLbStickySessionsType(raw interface{}) lbSDK.StickySessionsType {
	return lbSDK.StickySessionsType(raw.(string))
}

func flattenLbHCTCP(config *lbSDK.HealthCheckTCPConfig) interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{},
	}
}

func expandLbHCTCP(raw interface{}) *lbSDK.HealthCheckTCPConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	return &lbSDK.HealthCheckTCPConfig{}
}

func flattenLbHCHTTP(config *lbSDK.HealthCheckHTTPConfig) interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"uri":         config.URI,
			"method":      config.Method,
			"code":        flattenInt32Ptr(config.Code),
			"host_header": config.HostHeader,
		},
	}
}

func expandLbHCHTTP(raw interface{}) *lbSDK.HealthCheckHTTPConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lbSDK.HealthCheckHTTPConfig{
		URI:        rawMap["uri"].(string),
		Method:     rawMap["method"].(string),
		Code:       expandInt32Ptr(rawMap["code"]),
		HostHeader: rawMap["host_header"].(string),
	}
}

func flattenLbHCHTTPS(config *lbSDK.HealthCheckHTTPSConfig) interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"uri":         config.URI,
			"method":      config.Method,
			"code":        flattenInt32Ptr(config.Code),
			"host_header": config.HostHeader,
			"sni":         config.Sni,
		},
	}
}

func expandLbHCHTTPS(raw interface{}) *lbSDK.HealthCheckHTTPSConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lbSDK.HealthCheckHTTPSConfig{
		URI:        rawMap["uri"].(string),
		Method:     rawMap["method"].(string),
		Code:       expandInt32Ptr(rawMap["code"]),
		HostHeader: rawMap["host_header"].(string),
		Sni:        rawMap["sni"].(string),
	}
}

func expandLbHCMySQL(raw interface{}) *lbSDK.HealthCheckMysqlConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lbSDK.HealthCheckMysqlConfig{
		User: rawMap["database_user"].(string),
	}
}

func expandLbHCPgSQL(raw interface{}) *lbSDK.HealthCheckPgsqlConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &lbSDK.HealthCheckPgsqlConfig{
		User: rawMap["database_user"].(string),
	}
}

func expandLbHCLDAP(raw interface{}) *lbSDK.HealthCheckLdapConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	return &lbSDK.HealthCheckLdapConfig{}
}

func expandLbHealthCheck(raw interface{}) (*lbSDK.HealthCheck, error) {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil, nil
	}

	var err error

	rawMap := raw.([]interface{})[0].(map[string]interface{})

	protocol := rawMap["protocol"]
	hc := &lbSDK.HealthCheck{}
	hc.CheckMaxRetries = int32(rawMap["max_retries"].(int))
	hc.Port = int32(rawMap["port"].(int))
	hc.CheckDelay, err = expandDuration(rawMap["check_delay"])
	if err != nil {
		return nil, err
	}
	hc.CheckTimeout, err = expandDuration(rawMap["check_timeout"])
	if err != nil {
		return nil, err
	}
	hc.CheckSendProxy = rawMap["check_send_proxy"].(bool)
	hc.TransientCheckDelay = &scw.Duration{Seconds: int64(rawMap["transient_check_delay"].(int))}

	switch protocol {
	case "tcp":
		hc.TCPConfig = &lbSDK.HealthCheckTCPConfig{}
	case "http":
		hc.HTTPConfig = expandLbHCHTTP(raw)
	case "https":
		hc.HTTPSConfig = expandLbHCHTTPS(raw)
	case "mysql":
		hc.MysqlConfig = expandLbHCMySQL(raw)
	case "pgsql":
		hc.PgsqlConfig = expandLbHCPgSQL(raw)
	case "ldap":
		hc.LdapConfig = expandLbHCLDAP(raw)
	}

	return hc, err
}

func updateHealthCheckChanges(hc *lbSDK.HealthCheck, updateHCRequest *lbSDK.ZonedAPIUpdateHealthCheckRequest) {
	updateHCRequest.TransientCheckDelay = hc.TransientCheckDelay
	updateHCRequest.CheckSendProxy = hc.CheckSendProxy
	updateHCRequest.CheckTimeout = hc.CheckTimeout
	updateHCRequest.CheckMaxRetries = hc.CheckMaxRetries
	updateHCRequest.CheckDelay = hc.CheckDelay
	updateHCRequest.Port = hc.Port

	switch {
	case hc.TCPConfig != nil:
		updateHCRequest.TCPConfig = hc.TCPConfig

	case hc.PgsqlConfig != nil:
		updateHCRequest.PgsqlConfig = hc.PgsqlConfig

	case hc.MysqlConfig != nil:
		updateHCRequest.MysqlConfig = hc.MysqlConfig

	case hc.HTTPConfig != nil:
		updateHCRequest.HTTPConfig = hc.HTTPConfig

	case hc.HTTPSConfig != nil:
		updateHCRequest.HTTPSConfig = hc.HTTPSConfig

	case hc.LdapConfig != nil:
		updateHCRequest.LdapConfig = hc.LdapConfig

	case hc.RedisConfig != nil:
		updateHCRequest.RedisConfig = hc.RedisConfig
	}
}

func flattenLbHealthCheck(hc *lbSDK.HealthCheck) interface{} {
	if hc == nil {
		return nil
	}

	raw := make(map[string]interface{})
	raw["port"] = hc.Port
	raw["check_timeout"] = flattenDuration(hc.CheckTimeout)
	raw["check_delay"] = flattenDuration(hc.CheckDelay)
	raw["transient_check_delay"] = hc.TransientCheckDelay.Seconds
	raw["check_send_proxy"] = flattenBoolPtr(&hc.CheckSendProxy)
	raw["max_retries"] = hc.CheckMaxRetries

	if hc.HTTPConfig != nil {
		raw["protocol"] = "http"
		raw["uri"] = hc.HTTPConfig.URI
		raw["method"] = hc.HTTPConfig.Method
		raw["code"] = flattenInt32Ptr(hc.HTTPConfig.Code)
		raw["host_header"] = hc.HTTPConfig.HostHeader
	}

	if hc.HTTPSConfig != nil {
		raw["protocol"] = "https"
		raw["uri"] = hc.HTTPSConfig.URI
		raw["method"] = hc.HTTPSConfig.Method
		raw["code"] = flattenInt32Ptr(hc.HTTPSConfig.Code)
		raw["host_header"] = hc.HTTPSConfig.HostHeader
		raw["sni"] = hc.HTTPSConfig.Sni
	}

	if hc.LdapConfig != nil {
		raw["protocol"] = "ldap"
	}
	if hc.TCPConfig != nil {
		raw["protocol"] = "tcp"
	}
	if hc.MysqlConfig != nil {
		raw["protocol"] = "mysql"
		raw["database_user"] = hc.MysqlConfig.User
	}
	if hc.PgsqlConfig != nil {
		raw["protocol"] = "pgsql"
		raw["database_user"] = hc.PgsqlConfig.User
	}
	if hc.RedisConfig != nil {
		raw["protocol"] = "redis"
	}

	return []map[string]interface{}{raw}
}

func expandLbLetsEncrypt(raw interface{}) *lbSDK.CreateCertificateRequestLetsencryptConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	alternativeNames := rawMap["subject_alternative_name"].([]interface{})
	config := &lbSDK.CreateCertificateRequestLetsencryptConfig{
		CommonName: rawMap["common_name"].(string),
	}
	for _, alternativeName := range alternativeNames {
		config.SubjectAlternativeName = append(config.SubjectAlternativeName, alternativeName.(string))
	}
	return config
}

func expandLbCustomCertificate(raw interface{}) *lbSDK.CreateCertificateRequestCustomCertificate {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}

	rawMap := raw.([]interface{})[0].(map[string]interface{})
	config := &lbSDK.CreateCertificateRequestCustomCertificate{
		CertificateChain: rawMap["certificate_chain"].(string),
	}
	return config
}

func expandLbProxyProtocol(raw interface{}) lbSDK.ProxyProtocol {
	return lbSDK.ProxyProtocol("proxy_protocol_" + raw.(string))
}

func flattenLbProxyProtocol(pp lbSDK.ProxyProtocol) interface{} {
	return strings.TrimPrefix(pp.String(), "proxy_protocol_")
}

func lbUpgradeV1SchemaType() cty.Type {
	return cty.Object(map[string]cty.Type{
		"id": cty.String,
	})
}

// lbUpgradeV1UpgradeFunc allow upgrade the from regional to a zoned resource.
func lbUpgradeV1SchemaUpgradeFunc(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	var err error
	// element id: upgrade
	ID, exist := rawState["id"]
	if !exist {
		return nil, fmt.Errorf("upgrade: id not exist")
	}
	rawState["id"], err = lbUpgradeV1RegionalToZonedID(ID.(string))
	if err != nil {
		return nil, err
	}
	// return rawState updated
	return rawState, nil
}

func lbUpgradeV1RegionalToZonedID(element string) (string, error) {
	locality, id, err := parseLocalizedID(element)
	// return error if can't parse
	if err != nil {
		return "", fmt.Errorf("upgrade: could not retrieve the locality from `%s`", element)
	}
	// if locality is already zoned return
	if validator.IsZone(locality) {
		return element, nil
	}
	//  append zone 1 as default: e.g. fr-par-1
	return fmt.Sprintf("%s-1/%s", locality, id), nil
}

func expandLbPrivateNetworkStaticConfig(raw interface{}) *lbSDK.PrivateNetworkStaticConfig {
	if raw == nil || len(raw.([]interface{})) < 1 {
		return nil
	}
	return &lbSDK.PrivateNetworkStaticConfig{
		IPAddress: expandStrings(raw),
	}
}

func flattenLbPrivateNetworkStaticConfig(cfg *lbSDK.PrivateNetworkStaticConfig) []string {
	if cfg == nil {
		return nil
	}

	return cfg.IPAddress
}

func expandLbPrivateNetworkDHCPConfig(raw interface{}) *lbSDK.PrivateNetworkDHCPConfig {
	if raw == nil || !raw.(bool) {
		return nil
	}
	return &lbSDK.PrivateNetworkDHCPConfig{}
}

func waitForLB(ctx context.Context, lbAPI *lbSDK.ZonedAPI, zone scw.Zone, lbID string, timeout time.Duration) (*lbSDK.LB, error) {
	retryInterval := defaultWaitLBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	loadBalancer, err := lbAPI.WaitForLb(&lbSDK.ZonedAPIWaitForLBRequest{
		LBID:          lbID,
		Zone:          zone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))

	return loadBalancer, err
}

func waitForLbInstances(ctx context.Context, lbAPI *lbSDK.ZonedAPI, zone scw.Zone, lbID string, timeout time.Duration) (*lbSDK.LB, error) {
	retryInterval := defaultWaitLBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	loadBalancer, err := lbAPI.WaitForLbInstances(&lbSDK.ZonedAPIWaitForLBInstancesRequest{
		Zone:          zone,
		LBID:          lbID,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))

	return loadBalancer, err
}

func waitForLBPN(ctx context.Context, lbAPI *lbSDK.ZonedAPI, zone scw.Zone, lbID string, timeout time.Duration) ([]*lbSDK.PrivateNetwork, error) {
	retryInterval := defaultWaitLBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	privateNetworks, err := lbAPI.WaitForLBPN(&lbSDK.ZonedAPIWaitForLBPNRequest{
		LBID:          lbID,
		Zone:          zone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))

	return privateNetworks, err
}

func waitForLBCertificate(ctx context.Context, lbAPI *lbSDK.ZonedAPI, zone scw.Zone, id string, timeout time.Duration) (*lbSDK.Certificate, error) {
	retryInterval := defaultWaitLBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	certificate, err := lbAPI.WaitForLBCertificate(&lbSDK.ZonedAPIWaitForLBCertificateRequest{
		CertID:        id,
		Zone:          zone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))

	return certificate, err
}

func attachLBPrivateNetworks(ctx context.Context, lbAPI *lbSDK.ZonedAPI, zone scw.Zone, pnConfigs []*lbSDK.PrivateNetwork, lbID string, timeout time.Duration) ([]*lbSDK.PrivateNetwork, error) {
	var privateNetworks []*lbSDK.PrivateNetwork

	for i := range pnConfigs {
		pn, err := lbAPI.AttachPrivateNetwork(&lbSDK.ZonedAPIAttachPrivateNetworkRequest{
			Zone:             zone,
			LBID:             lbID,
			PrivateNetworkID: pnConfigs[i].PrivateNetworkID,
			StaticConfig:     pnConfigs[i].StaticConfig,
			DHCPConfig:       pnConfigs[i].DHCPConfig,
		}, scw.WithContext(ctx))
		if err != nil && !is404Error(err) {
			return nil, err
		}

		privateNetworks, err = waitForLBPN(ctx, lbAPI, zone, pn.LB.ID, timeout)
		if err != nil && !is404Error(err) {
			return nil, err
		}

		for _, pn := range privateNetworks {
			if pn.Status == lbSDK.PrivateNetworkStatusError {
				err = lbAPI.DetachPrivateNetwork(&lbSDK.ZonedAPIDetachPrivateNetworkRequest{
					Zone:             zone,
					LBID:             pn.LB.ID,
					PrivateNetworkID: pn.PrivateNetworkID,
				}, scw.WithContext(ctx))
				if err != nil && !is404Error(err) {
					return nil, err
				}
				tflog.Debug(ctx, fmt.Sprintf("DHCP config: %v", pn.DHCPConfig))
				tflog.Debug(ctx, fmt.Sprintf("Static config: %v", pn.StaticConfig))
				return nil, fmt.Errorf("attaching private network with id: %s on error state. please check your config", pn.PrivateNetworkID)
			}
		}
	}

	return privateNetworks, nil
}

func flattenLbInstances(instances []*lbSDK.Instance) interface{} {
	if instances == nil {
		return nil
	}
	flattenedInstances := []map[string]interface{}(nil)
	for _, instance := range instances {
		flattenedInstances = append(flattenedInstances, map[string]interface{}{
			"id":         instance.ID,
			"status":     instance.Status.String(),
			"ip_address": instance.IPAddress,
			"created_at": flattenTime(instance.CreatedAt),
			"updated_at": flattenTime(instance.UpdatedAt),
			"zone":       instance.Zone,
		})
	}
	return flattenedInstances
}

func flattenLbIPs(ips []*lbSDK.IP) interface{} {
	if ips == nil {
		return nil
	}
	flattenedIPs := []map[string]interface{}(nil)
	for _, ip := range ips {
		flattenedIPs = append(flattenedIPs, map[string]interface{}{
			"id":              ip.ID,
			"ip_address":      ip.IPAddress,
			"reverse":         ip.Reverse,
			"organization_id": ip.OrganizationID,
			"project_id":      ip.ProjectID,
			"zone":            ip.Zone,
			"lb_id":           flattenStringPtr(ip.LBID),
		})
	}
	return flattenedIPs
}

func ipv4Match(cidr, ipStr string) bool {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(ipStr)

	return cidrNet.Contains(ip)
}
