package domain

import (
	"context"
	"time"

	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
)

const (
	defaultDomainRecordTimeout     = 5 * time.Minute
	defaultDomainZoneTimeout       = 5 * time.Minute
	defaultDomainZoneRetryInterval = 5 * time.Second
)

func waitForDNSZone(ctx context.Context, domainAPI *domain.API, dnsZone string, timeout time.Duration) (*domain.DNSZone, error) {
	retryInterval := defaultDomainZoneRetryInterval
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}

	return domainAPI.WaitForDNSZone(&domain.WaitForDNSZoneRequest{
		DNSZone:       dnsZone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: scw.TimeDurationPtr(retryInterval),
	}, scw.WithContext(ctx))
}

func waitForDNSRecordExist(ctx context.Context, domainAPI *domain.API, dnsZone, recordName string, recordType domain.RecordType, timeout time.Duration) (*domain.Record, error) {
	retryInterval := defaultDomainZoneRetryInterval
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}

	return domainAPI.WaitForDNSRecordExist(&domain.WaitForDNSRecordExistRequest{
		DNSZone:       dnsZone,
		RecordName:    recordName,
		RecordType:    recordType,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: scw.TimeDurationPtr(retryInterval),
	}, scw.WithContext(ctx))
}

func waitForDomainsRegistration(ctx context.Context, api *domain.RegistrarAPI, domainName string, timeout time.Duration) (*domain.Domain, error) {
	retryInterval := defaultWaitDomainsRegistrationRetryInterval
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}
	return api.WaitForOrderDomain(&domain.WaitForOrderDomainRequest{
		Domain:        domainName,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))
}
