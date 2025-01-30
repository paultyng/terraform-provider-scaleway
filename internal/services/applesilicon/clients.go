package applesilicon

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

// newAPIWithZone returns a new apple silicon API and the zone
func newAPIWithZone(d *schema.ResourceData, m interface{}) (*applesilicon.API, scw.Zone, error) {
	asAPI := applesilicon.NewAPI(meta.ExtractScwClient(m))

	zone, err := meta.ExtractZone(d, m)
	if err != nil {
		return nil, "", err
	}
	return asAPI, zone, nil
}

// NewAPIWithZoneAndID returns an apple silicon API with zone and ID extracted from the state
func NewAPIWithZoneAndID(m interface{}, id string) (*applesilicon.API, scw.Zone, string, error) {
	asAPI := applesilicon.NewAPI(meta.ExtractScwClient(m))

	zone, ID, err := zonal.ParseID(id)
	if err != nil {
		return nil, "", "", err
	}
	return asAPI, zone, ID, nil
}

func newPrivateNetworkAPIWithZone(d *schema.ResourceData, m interface{}) (*applesilicon.PrivateNetworkAPI, scw.Zone, error) {
	privateNetworkAPI := applesilicon.NewPrivateNetworkAPI(meta.ExtractScwClient(m))

	zone, err := meta.ExtractZone(d, m)
	if err != nil {
		return nil, "", err
	}
	return privateNetworkAPI, zone, nil
}
