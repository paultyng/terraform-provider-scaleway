package edgeservices

import (
	"context"
	"time"

	edgeservices "github.com/scaleway/scaleway-sdk-go/api/edge_services/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
)

const (
	defaultEdgeServicesTimeout = 5 * time.Minute
)

func waitForPipeline(ctx context.Context, edgeServicesapi *edgeservices.API, id string, timeout time.Duration) (*edgeservices.Pipeline, error) {
	retryInterval := defaultEdgeServicesTimeout
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}

	pipeline, err := edgeServicesapi.WaitForPipeline(&edgeservices.WaitForPipelineRequest{
		PipelineID:    id,
		RetryInterval: &retryInterval,
		Timeout:       scw.TimeDurationPtr(timeout),
	}, scw.WithContext(ctx))

	return pipeline, err
}

func waitForPurge(ctx context.Context, edgeServicesapi *edgeservices.API, id string, timeout time.Duration) (*edgeservices.PurgeRequest, error) {
	retryInterval := defaultEdgeServicesTimeout
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}

	purgeRequest, err := edgeServicesapi.WaitForPurgeRequest(&edgeservices.WaitForPurgeRequestRequest{
		PurgeRequestID: id,
		RetryInterval:  &retryInterval,
		Timeout:        scw.TimeDurationPtr(timeout),
	}, scw.WithContext(ctx))

	return purgeRequest, err
}
