package scaleway

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	defaultRdbInstanceTimeout   = 15 * time.Minute
	defaultWaitRDBRetryInterval = 30 * time.Second
)

// newRdbAPI returns a new RDB API
func newRdbAPI(m interface{}) *rdb.API {
	meta := m.(*Meta)
	return rdb.NewAPI(meta.scwClient)
}

// rdbAPIWithRegion returns a new lb API and the region for a Create request
func rdbAPIWithRegion(d *schema.ResourceData, m interface{}) (*rdb.API, scw.Region, error) {
	meta := m.(*Meta)

	region, err := extractRegion(d, meta)
	if err != nil {
		return nil, "", err
	}
	return newRdbAPI(m), region, nil
}

// rdbAPIWithRegionAndID returns an lb API with region and ID extracted from the state
func rdbAPIWithRegionAndID(m interface{}, id string) (*rdb.API, scw.Region, string, error) {
	region, ID, err := parseRegionalID(id)
	if err != nil {
		return nil, "", "", err
	}
	return newRdbAPI(m), region, ID, nil
}

func flattenInstanceSettings(settings []*rdb.InstanceSetting) interface{} {
	res := make(map[string]string)
	for _, value := range settings {
		res[value.Name] = value.Value
	}

	return res
}

func expandInstanceSettings(i interface{}) []*rdb.InstanceSetting {
	rawRule := i.(map[string]interface{})
	var res []*rdb.InstanceSetting
	for key, value := range rawRule {
		res = append(res, &rdb.InstanceSetting{
			Name:  key,
			Value: value.(string),
		})
	}

	return res
}

func waitForRDBInstance(ctx context.Context, api *rdb.API, region scw.Region, id string, timeout time.Duration) (*rdb.Instance, error) {
	retryInterval := defaultWaitRDBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	return api.WaitForInstance(&rdb.WaitForInstanceRequest{
		Region:        region,
		Timeout:       scw.TimeDurationPtr(timeout),
		InstanceID:    id,
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))
}

func waitForRDBDatabaseBackup(ctx context.Context, api *rdb.API, region scw.Region, id string, timeout time.Duration) (*rdb.DatabaseBackup, error) {
	retryInterval := defaultWaitRDBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
		Region:           region,
		Timeout:          scw.TimeDurationPtr(timeout),
		DatabaseBackupID: id,
		RetryInterval:    &retryInterval,
	}, scw.WithContext(ctx))
}

func waitForRDBReadReplica(ctx context.Context, api *rdb.API, region scw.Region, id string, timeout time.Duration) (*rdb.ReadReplica, error) {
	retryInterval := defaultWaitRDBRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	return api.WaitForReadReplica(&rdb.WaitForReadReplicaRequest{
		Region:        region,
		Timeout:       scw.TimeDurationPtr(timeout),
		ReadReplicaID: id,
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))
}

func expandPrivateNetwork(data interface{}, exist bool) ([]*rdb.EndpointSpec, error) {
	if data == nil || !exist {
		return nil, nil
	}

	var res []*rdb.EndpointSpec
	for _, pn := range data.([]interface{}) {
		r := pn.(map[string]interface{})
		ipNet := r["ip_net"].(string)
		spec := &rdb.EndpointSpec{
			PrivateNetwork: &rdb.EndpointSpecPrivateNetwork{
				PrivateNetworkID: expandID(r["pn_id"].(string)),
			},
		}
		if len(ipNet) > 0 {
			ip, err := expandIPNet(r["ip_net"].(string))
			if err != nil {
				return res, err
			}
			spec.PrivateNetwork.ServiceIP = &ip
		} else {
			spec.PrivateNetwork.IpamConfig = &rdb.EndpointSpecPrivateNetworkIpamConfig{}
		}
		res = append(res, spec)
	}

	return res, nil
}

func expandLoadBalancer() []*rdb.EndpointSpec {
	var res []*rdb.EndpointSpec

	res = append(res, &rdb.EndpointSpec{
		LoadBalancer: &rdb.EndpointSpecLoadBalancer{},
	})

	return res
}

func endpointsToRemove(endPoints []*rdb.Endpoint, updates interface{}) (map[string]bool, error) {
	actions := make(map[string]bool)
	endpoints := make(map[string]*rdb.Endpoint)
	for _, e := range endPoints {
		// skip load balancer
		if e.PrivateNetwork != nil {
			actions[e.ID] = true
			endpoints[newZonedIDString(e.PrivateNetwork.Zone, e.PrivateNetwork.PrivateNetworkID)] = e
		}
	}

	// compare if private networks are persisted
	for _, raw := range updates.([]interface{}) {
		r := raw.(map[string]interface{})
		pnZonedID := r["pn_id"].(string)
		locality, id, err := parseLocalizedID(pnZonedID)
		if err != nil {
			return nil, err
		}

		pnUpdated, err := newEndPointPrivateNetworkDetails(id, r["ip_net"].(string), locality)
		if err != nil {
			return nil, err
		}
		endpoint, exist := endpoints[pnZonedID]
		if !exist {
			continue
		}
		// match the endpoint id for a private network
		actions[endpoint.ID] = !isEndPointEqual(endpoints[pnZonedID].PrivateNetwork, pnUpdated)
	}

	return actions, nil
}

func newEndPointPrivateNetworkDetails(id, ip, locality string) (*rdb.EndpointPrivateNetworkDetails, error) {
	serviceIP, err := expandIPNet(ip)
	if err != nil {
		return nil, err
	}
	return &rdb.EndpointPrivateNetworkDetails{
		PrivateNetworkID: id,
		ServiceIP:        serviceIP,
		Zone:             scw.Zone(locality),
	}, nil
}

func isEndPointEqual(a, b interface{}) bool {
	// Find out the diff Private Network or not
	if _, ok := a.(*rdb.EndpointPrivateNetworkDetails); ok {
		if _, ok := b.(*rdb.EndpointPrivateNetworkDetails); ok {
			detailsA := a.(*rdb.EndpointPrivateNetworkDetails)
			detailsB := b.(*rdb.EndpointPrivateNetworkDetails)
			return reflect.DeepEqual(detailsA, detailsB)
		}
	}
	return false
}

func flattenPrivateNetwork(endpoints []*rdb.Endpoint) (interface{}, bool) {
	pnI := []map[string]interface{}(nil)
	for _, endpoint := range endpoints {
		if endpoint.PrivateNetwork != nil {
			pn := endpoint.PrivateNetwork
			serviceIP, err := flattenIPNet(pn.ServiceIP)
			if err != nil {
				return pnI, false
			}
			pnI = append(pnI, map[string]interface{}{
				"endpoint_id": endpoint.ID,
				"ip":          flattenIPPtr(endpoint.IP),
				"port":        int(endpoint.Port),
				"name":        endpoint.Name,
				"ip_net":      serviceIP,
				"pn_id":       pn.PrivateNetworkID,
				"hostname":    flattenStringPtr(endpoint.Hostname),
			})
			return pnI, true
		}
	}

	return pnI, false
}

func flattenLoadBalancer(endpoints []*rdb.Endpoint) interface{} {
	flat := []map[string]interface{}(nil)
	for _, endpoint := range endpoints {
		if endpoint.LoadBalancer != nil {
			flat = append(flat, map[string]interface{}{
				"endpoint_id": endpoint.ID,
				"ip":          flattenIPPtr(endpoint.IP),
				"port":        int(endpoint.Port),
				"name":        endpoint.Name,
				"hostname":    flattenStringPtr(endpoint.Hostname),
			})
			return flat
		}
	}

	return flat
}

// expandTimePtr returns a time pointer for an RFC3339 time.
// It returns nil if time is not valid, you should use validateDate to validate field.
func expandTimePtr(i interface{}) *time.Time {
	rawTime := expandStringPtr(i)
	if rawTime == nil {
		return nil
	}
	parsedTime, err := time.Parse(time.RFC3339, *rawTime)
	if err != nil {
		return nil
	}
	return &parsedTime
}

func expandReadReplicaEndpointsSpecDirectAccess(data interface{}) *rdb.ReadReplicaEndpointSpec {
	if data == nil || len(data.([]interface{})) == 0 {
		return nil
	}

	return &rdb.ReadReplicaEndpointSpec{
		DirectAccess: new(rdb.ReadReplicaEndpointSpecDirectAccess),
	}
}

// expandReadReplicaEndpointsSpecPrivateNetwork expand read-replica private network endpoints from schema to specs
func expandReadReplicaEndpointsSpecPrivateNetwork(data interface{}) (*rdb.ReadReplicaEndpointSpec, error) {
	if data == nil || len(data.([]interface{})) == 0 {
		return nil, nil
	}
	// private_network is a list of size 1
	data = data.([]interface{})[0]

	rawEndpoint := data.(map[string]interface{})

	endpoint := new(rdb.ReadReplicaEndpointSpec)

	serviceIP := rawEndpoint["service_ip"].(string)
	endpoint.PrivateNetwork = &rdb.ReadReplicaEndpointSpecPrivateNetwork{
		PrivateNetworkID: expandID(rawEndpoint["private_network_id"]),
	}
	if len(serviceIP) > 0 {
		ipNet, err := expandIPNet(serviceIP)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private_network service_ip (%s): %w", rawEndpoint["service_ip"], err)
		}
		endpoint.PrivateNetwork.ServiceIP = &ipNet
	} else {
		endpoint.PrivateNetwork.IpamConfig = &rdb.ReadReplicaEndpointSpecPrivateNetworkIpamConfig{}
	}

	return endpoint, nil
}

// flattenReadReplicaEndpoints flatten read-replica endpoints to directAccess and privateNetwork
func flattenReadReplicaEndpoints(endpoints []*rdb.Endpoint) (directAccess, privateNetwork interface{}) {
	for _, endpoint := range endpoints {
		rawEndpoint := map[string]interface{}{
			"endpoint_id": endpoint.ID,
			"ip":          flattenIPPtr(endpoint.IP),
			"port":        int(endpoint.Port),
			"name":        endpoint.Name,
			"hostname":    flattenStringPtr(endpoint.Hostname),
		}
		if endpoint.DirectAccess != nil {
			directAccess = rawEndpoint
		}
		if endpoint.PrivateNetwork != nil {
			rawEndpoint["private_network_id"] = endpoint.PrivateNetwork.PrivateNetworkID
			rawEndpoint["service_ip"] = endpoint.PrivateNetwork.ServiceIP.String()
			rawEndpoint["zone"] = endpoint.PrivateNetwork.Zone
			privateNetwork = rawEndpoint
		}
	}

	// direct_access and private_network are lists

	if directAccess != nil {
		directAccess = []interface{}{directAccess}
	}
	if privateNetwork != nil {
		privateNetwork = []interface{}{privateNetwork}
	}

	return directAccess, privateNetwork
}

// rdbPrivilegeV1SchemaUpgradeFunc allow upgrade the privilege ID on schema V1
func rdbPrivilegeV1SchemaUpgradeFunc(_ context.Context, rawState map[string]interface{}, m interface{}) (map[string]interface{}, error) {
	idRaw, exist := rawState["id"]
	if !exist {
		return nil, fmt.Errorf("upgrade: id not exist")
	}

	idParts := strings.Split(idRaw.(string), "/")
	if len(idParts) == 4 {
		return rawState, nil
	}

	region, idStr, err := parseRegionalID(idRaw.(string))
	if err != nil {
		// force the default region
		meta := m.(*Meta)
		defaultRegion, exist := meta.scwClient.GetDefaultRegion()
		if exist {
			region = defaultRegion
		}
	}

	databaseName := rawState["database_name"].(string)
	userName := rawState["user_name"].(string)
	rawState["id"] = resourceScalewayRdbUserPrivilegeID(region, idStr, databaseName, userName)
	rawState["region"] = region.String()

	return rawState, nil
}

func rdbPrivilegeUpgradeV1SchemaType() cty.Type {
	return cty.Object(map[string]cty.Type{
		"id": cty.String,
	})
}
