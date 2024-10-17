package ipam

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipam "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

const (
	defaultIPRetryInterval     = 5 * time.Second
	defaultIPReverseDNSTimeout = 10 * time.Minute
)

// newAPIWithRegion returns a new ipam API and the region
func newAPIWithRegion(d *schema.ResourceData, m interface{}) (*ipam.API, scw.Region, error) {
	ipamAPI := ipam.NewAPI(meta.ExtractScwClient(m))

	region, err := meta.ExtractRegion(d, m)
	if err != nil {
		return nil, "", err
	}

	return ipamAPI, region, nil
}

// NewAPIWithRegionAndID returns a new ipam API with locality and ID extracted from the state
func NewAPIWithRegionAndID(m interface{}, id string) (*ipam.API, scw.Region, string, error) {
	ipamAPI := ipam.NewAPI(meta.ExtractScwClient(m))

	region, ID, err := regional.ParseID(id)
	if err != nil {
		return nil, "", "", err
	}
	return ipamAPI, region, ID, err
}

func diffSuppressFuncStandaloneIPandCIDR(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	oldIP, oldNet, errOld := net.ParseCIDR(oldValue)
	if errOld != nil {
		oldIP = net.ParseIP(oldValue)
	}

	newIP, newNet, errNew := net.ParseCIDR(newValue)
	if errNew != nil {
		newIP = net.ParseIP(newValue)
	}

	if oldIP != nil && newIP != nil && oldIP.Equal(newIP) {
		return true
	}

	if oldNet != nil && newIP != nil && oldNet.Contains(newIP) {
		return true
	}

	if newNet != nil && oldIP != nil && newNet.Contains(oldIP) {
		return true
	}

	return false
}

type GetResourcePrivateIPsOptions struct {
	ResourceType     *ipam.ResourceType
	ResourceID       *string
	ResourceName     *string
	PrivateNetworkID *string
}

// GetResourcePrivateIPs fetches the private IP addresses of a resource in a private network.
func GetResourcePrivateIPs(ctx context.Context, m interface{}, region scw.Region, opts *GetResourcePrivateIPsOptions) ([]map[string]interface{}, error) {
	ipamAPI := ipam.NewAPI(meta.ExtractScwClient(m))

	req := &ipam.ListIPsRequest{
		Region: region,
	}

	if opts != nil {
		if opts.PrivateNetworkID != nil {
			req.PrivateNetworkID = opts.PrivateNetworkID
		}
		if opts.ResourceID != nil {
			req.ResourceID = opts.ResourceID
		}
		if opts.ResourceName != nil {
			req.ResourceName = opts.ResourceName
		}
		if opts.ResourceType != nil {
			req.ResourceType = *opts.ResourceType
		}
	}

	resp, err := ipamAPI.ListIPs(req, scw.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error fetching IPs from IPAM: %w", err)
	}

	if len(resp.IPs) == 0 {
		return nil, nil
	}

	ipList := make([]map[string]interface{}, 0, len(resp.IPs))
	for _, ip := range resp.IPs {
		ipNet := ip.Address
		if ipNet.IP == nil {
			continue
		}
		ipMap := map[string]interface{}{
			"id":      regional.NewIDString(region, ip.ID),
			"address": ipNet.IP.String(),
		}
		ipList = append(ipList, ipMap)
	}

	return ipList, nil
}
