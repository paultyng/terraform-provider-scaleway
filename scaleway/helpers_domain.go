package scaleway

import (
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2alpha2"
)

// domainAPI returns a new domain API.
func domainAPI(m interface{}) *domain.API {
	meta := m.(*Meta)

	return domain.NewAPI(meta.scwClient)
}
