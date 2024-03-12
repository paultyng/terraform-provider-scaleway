package scaleway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

// marketplaceAPIWithZone returns a new marketplace API and the zone for a Create request
func marketplaceAPIWithZone(d *schema.ResourceData, m interface{}) (*marketplace.API, scw.Zone, error) {
	meta := m.(*meta.Meta)
	marketplaceAPI := marketplace.NewAPI(meta.ScwClient())

	zone, err := extractZone(d, meta)
	if err != nil {
		return nil, "", err
	}
	return marketplaceAPI, zone, nil
}
