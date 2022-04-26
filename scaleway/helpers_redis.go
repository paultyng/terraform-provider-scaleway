package scaleway

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/scaleway/scaleway-sdk-go/api/redis/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	defaultRedisClusterTimeout           = 15 * time.Minute
	defaultWaitRedisClusterRetryInterval = 5 * time.Second
)

// newRedisApi returns a new Redis API
func newRedisAPI(m interface{}) *redis.API {
	meta := m.(*Meta)
	return redis.NewAPI(meta.scwClient)
}

// redisAPIWithZone returns a new Redis API and the zone for a Create request
func redisAPIWithZone(d *schema.ResourceData, m interface{}) (*redis.API, scw.Zone, error) {
	meta := m.(*Meta)

	zone, err := extractZone(d, meta)
	if err != nil {
		return nil, "", err
	}
	return newRedisAPI(m), zone, nil
}

// redisAPIWithZoneAndID returns a Redis API with zone and ID extracted from the state
func redisAPIWithZoneAndID(m interface{}, id string) (*redis.API, scw.Zone, string, error) {
	zone, ID, err := parseZonedID(id)
	if err != nil {
		return nil, "", "", err
	}
	return newRedisAPI(m), zone, ID, nil
}

func waitForRedisCluster(ctx context.Context, api *redis.API, zone scw.Zone, id string, timeout time.Duration) (*redis.Cluster, error) {
	retryInterval := defaultWaitRedisClusterRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	return api.WaitForCluster(&redis.WaitForClusterRequest{
		Zone:          zone,
		Timeout:       scw.TimeDurationPtr(timeout),
		ClusterID:     id,
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))
}
