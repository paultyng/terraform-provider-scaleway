package scaleway

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	defaultSecretTimeout = 5 * time.Minute
)

// secretAPIWithRegion returns a new Secret API and the region for a Create request
func secretAPIWithRegion(d *schema.ResourceData, m interface{}) (*secret.API, scw.Region, error) {
	meta := m.(*Meta)
	api := secret.NewAPI(meta.scwClient)

	region, err := extractRegion(d, meta)
	if err != nil {
		return nil, "", err
	}
	return api, region, nil
}

// secretAPIWithRegionAndID returns a Secret API with locality and ID extracted from the state
func secretAPIWithRegionAndID(m interface{}, id string) (*secret.API, scw.Region, string, error) {
	meta := m.(*Meta)
	api := secret.NewAPI(meta.scwClient)

	region, id, err := parseRegionalID(id)
	if err != nil {
		return nil, "", "", err
	}
	return api, region, id, nil
}

// secretVersionAPIWithRegionAndID returns a Secret API with locality and Nested ID extracted from the state
func secretVersionAPIWithRegionAndID(m interface{}, id string) (*secret.API, scw.Region, string, string, error) {
	meta := m.(*Meta)
	api := secret.NewAPI(meta.scwClient)

	region, id, revision, err := parseLocalizedNestedID(id)
	if err != nil {
		return nil, "", "", "", err
	}
	return api, scw.Region(region), id, revision, nil
}
