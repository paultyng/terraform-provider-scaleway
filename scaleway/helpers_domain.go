package scaleway

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/exp/slices"
)

const (
	defaultDomainRecordTimeout     = 30 * time.Second
	defaultDomainZoneTimeout       = 30 * time.Second
	defaultDomainZoneRetryInterval = 5 * time.Second
)

// domainAPI returns a new domain API.
func newDomainAPI(m interface{}) *domain.API {
	meta := m.(*Meta)

	return domain.NewAPI(meta.scwClient)
}

func flattenDomainData(data string, recordType domain.RecordType) interface{} {
	switch recordType {
	case domain.RecordTypeMX: // API return this format: "{priority} {data}"
		dataSplit := strings.SplitN(data, " ", 2)
		if len(dataSplit) == 2 {
			return dataSplit[1]
		}
	case domain.RecordTypeTXT:
		return strings.Trim(data, "\"")
	}

	return data
}

func getRecordFromTypeAndData(dnsType domain.RecordType, data string, records []*domain.Record) (*domain.Record, error) {
	var currentRecord *domain.Record
	for _, r := range records {
		flattedData := flattenDomainData(strings.ToLower(r.Data), r.Type).(string)
		flattenCurrentData := flattenDomainData(strings.ToLower(data), r.Type).(string)
		if strings.HasPrefix(flattedData, flattenCurrentData) && r.Type == dnsType {
			if currentRecord != nil {
				return nil, fmt.Errorf("multiple records found with same type and data")
			}
			currentRecord = r
			break
		}
	}

	if currentRecord == nil {
		return nil, fmt.Errorf("record with type %s and data %s not found", dnsType.String(), data)
	}

	return currentRecord, nil
}

func flattenDomainGeoIP(config *domain.RecordGeoIPConfig) interface{} {
	flattened := []map[string]interface{}{}

	if config == nil {
		return flattened
	}

	flattened = []map[string]interface{}{{}}
	if config.Matches != nil && len(config.Matches) > 0 {
		matches := []map[string]interface{}{}
		for _, match := range config.Matches {
			rawMatch := map[string]interface{}{
				"data": match.Data,
			}
			if len(match.Continents) > 0 {
				rawMatch["continents"] = match.Continents
			}
			if len(match.Countries) > 0 {
				rawMatch["countries"] = match.Countries
			}

			matches = append(matches, rawMatch)
		}
		flattened[0]["matches"] = matches
	}

	return flattened
}

func expandDomainGeoIPConfig(defaultData string, i interface{}, ok bool) *domain.RecordGeoIPConfig {
	if i == nil || !ok {
		return nil
	}

	rawMap := i.([]interface{})[0].(map[string]interface{})

	config := domain.RecordGeoIPConfig{
		Default: defaultData,
	}

	rawMatches, ok := rawMap["matches"].([]interface{})
	if !ok && len(rawMatches) > 0 {
		return &config
	}

	matches := []*domain.RecordGeoIPConfigMatch{}
	for _, rawMatch := range rawMatches {
		rawMatchMap := rawMatch.(map[string]interface{})

		match := &domain.RecordGeoIPConfigMatch{
			Data: rawMatchMap["data"].(string),
		}

		rawContinents, ok := rawMatchMap["continents"].([]interface{})
		if ok {
			match.Continents = []string{}
			for _, rawContinent := range rawContinents {
				match.Continents = append(match.Continents, rawContinent.(string))
			}
		}
		rawCountries, ok := rawMatchMap["countries"].([]interface{})
		if ok {
			match.Countries = []string{}
			for _, rawCountry := range rawCountries {
				match.Countries = append(match.Countries, rawCountry.(string))
			}
		}

		matches = append(matches, match)
	}
	config.Matches = matches

	return &config
}

func flattenDomainHTTPService(config *domain.RecordHTTPServiceConfig) interface{} {
	flattened := []map[string]interface{}{}

	if config == nil {
		return flattened
	}

	ips := []interface{}{}
	if config.IPs != nil && len(config.IPs) > 0 {
		for _, ip := range config.IPs {
			ips = append(ips, ip.String())
		}
	}
	return []map[string]interface{}{
		{
			"must_contain": flattenStringPtr(config.MustContain),
			"url":          config.URL,
			"user_agent":   flattenStringPtr(config.UserAgent),
			"strategy":     config.Strategy.String(),
			"ips":          ips,
		},
	}
}

func expandDomainHTTPService(i interface{}, ok bool) *domain.RecordHTTPServiceConfig {
	if i == nil || !ok {
		return nil
	}

	rawMap := i.([]interface{})[0].(map[string]interface{})

	ips := []net.IP{}
	rawIPs, ok := rawMap["ips"].([]interface{})
	if ok {
		for _, rawIP := range rawIPs {
			ips = append(ips, net.ParseIP(rawIP.(string)))
		}
	}

	return &domain.RecordHTTPServiceConfig{
		MustContain: expandStringPtr(rawMap["must_contain"]),
		URL:         rawMap["url"].(string),
		UserAgent:   expandStringPtr(rawMap["user_agent"]),
		Strategy:    domain.RecordHTTPServiceConfigStrategy(rawMap["strategy"].(string)),
		IPs:         ips,
	}
}

func flattenDomainWeighted(config *domain.RecordWeightedConfig) interface{} {
	flattened := []map[string]interface{}{}

	if config == nil {
		return flattened
	}

	if config.WeightedIPs != nil && len(config.WeightedIPs) > 0 {
		for _, weightedIPs := range config.WeightedIPs {
			flattened = append(flattened, map[string]interface{}{
				"ip":     weightedIPs.IP.String(),
				"weight": int(weightedIPs.Weight),
			})
		}
	}

	return flattened
}

func expandDomainWeighted(i interface{}, ok bool) *domain.RecordWeightedConfig {
	if i == nil || !ok {
		return nil
	}

	weightedIPs := []*domain.RecordWeightedConfigWeightedIP{}
	if raw := i.([]interface{}); len(raw) > 0 {
		for _, rawWeighted := range raw {
			rawMap := rawWeighted.(map[string]interface{})
			weightedIPs = append(weightedIPs, &domain.RecordWeightedConfigWeightedIP{
				IP:     net.ParseIP(rawMap["ip"].(string)),
				Weight: uint32(rawMap["weight"].(int)),
			})
		}
	}
	return &domain.RecordWeightedConfig{
		WeightedIPs: weightedIPs,
	}
}

func flattenDomainView(config *domain.RecordViewConfig) interface{} {
	flattened := []map[string]interface{}{}

	if config == nil {
		return flattened
	}

	if config.Views != nil && len(config.Views) > 0 {
		for _, view := range config.Views {
			flattened = append(flattened, map[string]interface{}{
				"subnet": view.Subnet,
				"data":   view.Data,
			})
		}
	}

	return flattened
}

func expandDomainView(i interface{}, ok bool) *domain.RecordViewConfig {
	if i == nil || !ok {
		return nil
	}

	views := []*domain.RecordViewConfigView{}
	if raw := i.([]interface{}); len(raw) > 0 {
		for _, rawWeighted := range raw {
			rawMap := rawWeighted.(map[string]interface{})
			views = append(views, &domain.RecordViewConfigView{
				Subnet: rawMap["subnet"].(string),
				Data:   rawMap["data"].(string),
			})
		}
	}

	return &domain.RecordViewConfig{
		Views: views,
	}
}

func waitForDNSZone(ctx context.Context, domainAPI *domain.API, dnsZone string, timeout time.Duration) (*domain.DNSZone, error) {
	retryInterval := defaultDomainZoneRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	return domainAPI.WaitForDNSZone(&domain.WaitForDNSZoneRequest{
		DNSZone:       dnsZone,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: scw.TimeDurationPtr(retryInterval),
	}, scw.WithContext(ctx))
}

var disableHostResolver bool

func hostResolver(ctx context.Context, timeout time.Duration, reverse, ip string) bool {
	if disableHostResolver {
		return true
	}
	ticker := time.Tick(time.Millisecond * 500)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(500),
			}
			conn, err := d.DialContext(ctx, network, "ns0.dom.scw.cloud:53")
			if err != nil {
				conn, err = d.DialContext(ctx, network, "ns1.dom.scw.cloud:53")
			}
			return conn, err
		},
	}

	for range ticker {
		address, err := r.LookupHost(ctx, reverse)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return false
			}
			continue
		}
		if slices.Contains(address, ip) {
			return true
		}
	}

	return false
}

func findDefaultReverse(address string) string {
	parts := strings.Split(address, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "-") + ".instances.scw.cloud"
}
