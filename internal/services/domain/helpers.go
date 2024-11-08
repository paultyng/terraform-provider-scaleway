package domain

import (
	"errors"
	"fmt"
	"strings"

	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

// NewDomainAPI returns a new domain API.
func NewDomainAPI(m interface{}) *domain.API {
	return domain.NewAPI(meta.ExtractScwClient(m))
}

func NewRegistrarDomainAPI(m interface{}) *domain.RegistrarAPI {
	return domain.NewRegistrarAPI(meta.ExtractScwClient(m))
}

func getRecordFromTypeAndData(dnsType domain.RecordType, data string, records []*domain.Record) (*domain.Record, error) {
	var currentRecord *domain.Record
	for _, r := range records {
		flattedData := flattenDomainData(strings.ToLower(r.Data), r.Type).(string)
		flattenCurrentData := flattenDomainData(strings.ToLower(data), r.Type).(string)
		if strings.HasPrefix(flattedData, flattenCurrentData) && r.Type == dnsType {
			if currentRecord != nil {
				return nil, errors.New("multiple records found with same type and data")
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

func FindDefaultReverse(address string) string {
	parts := strings.Split(address, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "-") + ".instances.scw.cloud"
}
