package applesilicon

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
)

func detachAllPrivateNetworkFromServer(ctx context.Context, d *schema.ResourceData, m interface{}, serverID string) error {
	privateNetworkAPI, zone, err := newPrivateNetworkAPIWithZone(d, m)
	if err != nil {
		return err
	}
	listPrivateNetwork, err := privateNetworkAPI.ListServerPrivateNetworks(&applesilicon.PrivateNetworkAPIListServerPrivateNetworksRequest{
		Zone:     zone,
		ServerID: &serverID,
	}, scw.WithContext(ctx))
	if err != nil {
		return err
	}

	for _, pn := range listPrivateNetwork.ServerPrivateNetworks {
		err := privateNetworkAPI.DeleteServerPrivateNetwork(&applesilicon.PrivateNetworkAPIDeleteServerPrivateNetworkRequest{
			Zone:             zone,
			ServerID:         serverID,
			PrivateNetworkID: pn.PrivateNetworkID,
		}, scw.WithContext(ctx))
		if err != nil {
			return err
		}
	}

	_, err = waitForAppleSiliconPrivateNetworkServer(ctx, privateNetworkAPI, zone, serverID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return err
	}
	return nil
}
