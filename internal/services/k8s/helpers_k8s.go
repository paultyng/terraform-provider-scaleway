package k8s

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

const (
	defaultK8SClusterTimeout = 15 * time.Minute
	defaultK8SPoolTimeout    = 30 * time.Minute
	defaultK8SRetryInterval  = 5 * time.Second
)

func newAPIWithRegion(d *schema.ResourceData, m interface{}) (*k8s.API, scw.Region, error) {
	k8sAPI := k8s.NewAPI(meta.ExtractScwClient(m))

	region, err := meta.ExtractRegion(d, m)
	if err != nil {
		return nil, "", err
	}
	return k8sAPI, region, nil
}

func NewAPIWithRegionAndID(m interface{}, id string) (*k8s.API, scw.Region, string, error) {
	k8sAPI := k8s.NewAPI(meta.ExtractScwClient(m))

	region, ID, err := regional.ParseID(id)
	if err != nil {
		return nil, "", "", err
	}
	return k8sAPI, region, ID, nil
}

func GetMinorVersionFromFull(version string) (string, error) {
	versionSplit := strings.Split(version, ".")
	if len(versionSplit) != 3 {
		return "", errors.New("version is not a full x.y.z version") // shoud never happen
	}

	return versionSplit[0] + "." + versionSplit[1], nil
}

// k8sGetLatestVersionFromMinor returns the latest full version (x.y.z) for a given minor version (x.y)
func k8sGetLatestVersionFromMinor(ctx context.Context, k8sAPI *k8s.API, region scw.Region, version string) (string, error) {
	versionSplit := strings.Split(version, ".")
	if len(versionSplit) != 2 {
		return "", fmt.Errorf("minor version should be like x.y not %s", version)
	}

	versionsResp, err := k8sAPI.ListVersions(&k8s.ListVersionsRequest{
		Region: region,
	}, scw.WithContext(ctx))
	if err != nil {
		return "", err
	}

	for _, v := range versionsResp.Versions {
		vSplit := strings.Split(v.Name, ".")
		if len(vSplit) != 3 {
			return "", fmt.Errorf("upstream version %s is not correctly formatted", v.Name) // should never happen
		}
		if versionSplit[0] == vSplit[0] && versionSplit[1] == vSplit[1] {
			return v.Name, nil
		}
	}
	return "", fmt.Errorf("no available upstream version found for %s", version)
}

// convert a list of nodes to a list of map
func convertNodes(res *k8s.ListNodesResponse) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(res.Nodes))
	for _, node := range res.Nodes {
		n := make(map[string]interface{})
		n["id"] = node.ID
		n["name"] = node.Name
		n["status"] = node.Status.String()
		if node.PublicIPV4 != nil && node.PublicIPV4.String() != types.NetIPNil { //nolint:staticcheck
			n["public_ip"] = node.PublicIPV4.String() //nolint:staticcheck
		}
		if node.PublicIPV6 != nil && node.PublicIPV6.String() != types.NetIPNil { //nolint:staticcheck
			n["public_ip_v6"] = node.PublicIPV6.String() //nolint:staticcheck
		}
		result = append(result, n)
	}
	return result
}

func getNodes(ctx context.Context, k8sAPI *k8s.API, pool *k8s.Pool) ([]map[string]interface{}, error) {
	req := &k8s.ListNodesRequest{
		Region:    pool.Region,
		ClusterID: pool.ClusterID,
		PoolID:    &pool.ID,
	}

	nodes, err := k8sAPI.ListNodes(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return convertNodes(nodes), nil
}
