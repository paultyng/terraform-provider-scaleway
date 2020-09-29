package scaleway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// registryAPIWithRegion returns a new container registry API and the region.
func registryAPIWithRegion(d *schema.ResourceData, m interface{}) (*registry.API, scw.Region, error) {
	meta := m.(*Meta)
	api := registry.NewAPI(meta.scwClient)

	region, err := extractRegion(d, meta)
	return api, region, err
}

// registryAPIWithRegionAndID returns a new container registry API, region and ID.
func registryAPIWithRegionAndID(m interface{}, id string) (*registry.API, scw.Region, string, error) {
	meta := m.(*Meta)
	api := registry.NewAPI(meta.scwClient)

	region, id, err := parseRegionalID(id)
	return api, region, id, err
}
