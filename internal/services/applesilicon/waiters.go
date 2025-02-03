package applesilicon

import (
	"context"
	"time"

	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
)

const (
	defaultAppleSiliconServerTimeout       = 20 * time.Minute
	defaultAppleSiliconServerRetryInterval = 5 * time.Second
)

func waitForAppleSiliconServer(ctx context.Context, api *applesilicon.API, zone scw.Zone, serverID string, timeout time.Duration) (*applesilicon.Server, error) {
	retryInterval := defaultAppleSiliconServerRetryInterval
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}

	server, err := api.WaitForServer(&applesilicon.WaitForServerRequest{
		ServerID:      serverID,
		Zone:          zone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))

	return server, err
}

func waitForAppleSiliconPrivateNetworkServer(ctx context.Context, api *applesilicon.PrivateNetworkAPI, zone scw.Zone, serverID string, timeout time.Duration) ([]*applesilicon.ServerPrivateNetwork, error) {
	retryInterval := defaultAppleSiliconServerRetryInterval
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}

	privateNetwork, err := api.WaitForServerPrivateNetworks(&applesilicon.WaitForServerRequest{
		ServerID:      serverID,
		Zone:          zone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))

	return privateNetwork, err
}
