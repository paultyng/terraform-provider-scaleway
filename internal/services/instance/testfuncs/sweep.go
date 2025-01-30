package instancetestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
)

func AddTestSweepers() {
	resource.AddTestSweepers("scaleway_instance_image", &resource.Sweeper{
		Name:         "scaleway_instance_image",
		Dependencies: []string{"scaleway_instance_server"},
		F:            testSweepImage,
	})
	resource.AddTestSweepers("scaleway_instance_ip", &resource.Sweeper{
		Name: "scaleway_instance_ip",
		F:    testSweepIP,
	})
	resource.AddTestSweepers("scaleway_instance_placement_group", &resource.Sweeper{
		Name: "scaleway_instance_placement_group",
		F:    testSweepPlacementGroup,
	})
	resource.AddTestSweepers("scaleway_instance_security_group", &resource.Sweeper{
		Name: "scaleway_instance_security_group",
		F:    testSweepSecurityGroup,
	})
	resource.AddTestSweepers("scaleway_instance_server", &resource.Sweeper{
		Name: "scaleway_instance_server",
		F:    testSweepServer,
	})
	resource.AddTestSweepers("scaleway_instance_snapshot", &resource.Sweeper{
		Name: "scaleway_instance_snapshot",
		F:    testSweepSnapshot,
	})
	resource.AddTestSweepers("scaleway_instance_volume", &resource.Sweeper{
		Name: "scaleway_instance_volume",
		F:    testSweepVolume,
	})
}

func testSweepVolume(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instanceSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the volumes in (%s)", zone)

		listVolumesResponse, err := instanceAPI.ListVolumes(&instanceSDK.ListVolumesRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing volumes in sweeper: %s", err)
		}

		for _, volume := range listVolumesResponse.Volumes {
			if volume.Server == nil {
				err := instanceAPI.DeleteVolume(&instanceSDK.DeleteVolumeRequest{
					Zone:     zone,
					VolumeID: volume.ID,
				})
				if err != nil {
					return fmt.Errorf("error deleting volume in sweeper: %s", err)
				}
			}
		}
		return nil
	})
}

func testSweepSnapshot(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		api := instanceSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying instance snapshots in (%+v)", zone)

		listSnapshotsResponse, err := api.ListSnapshots(&instanceSDK.ListSnapshotsRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing instance snapshots in sweeper: %w", err)
		}

		for _, snapshot := range listSnapshotsResponse.Snapshots {
			err := api.DeleteSnapshot(&instanceSDK.DeleteSnapshotRequest{
				Zone:       zone,
				SnapshotID: snapshot.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting instance snapshot in sweeper: %w", err)
			}
		}

		return nil
	})
}

func testSweepServer(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instanceSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the instanceSDK server in (%s)", zone)
		listServers, err := instanceAPI.ListServers(&instanceSDK.ListServersRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing servers in (%s) in sweeper: %s", zone, err)
			return nil
		}

		for _, srv := range listServers.Servers {
			if srv.State == instanceSDK.ServerStateStopped || srv.State == instanceSDK.ServerStateStoppedInPlace {
				err := instanceAPI.DeleteServer(&instanceSDK.DeleteServerRequest{
					Zone:     zone,
					ServerID: srv.ID,
				})
				if err != nil {
					return fmt.Errorf("error deleting server in sweeper: %s", err)
				}
			} else if srv.State == instanceSDK.ServerStateRunning {
				_, err := instanceAPI.ServerAction(&instanceSDK.ServerActionRequest{
					Zone:     zone,
					ServerID: srv.ID,
					Action:   instanceSDK.ServerActionTerminate,
				})
				if err != nil {
					return fmt.Errorf("error terminating server in sweeper: %s", err)
				}
			}
		}

		return nil
	})
}

func testSweepSecurityGroup(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instanceSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the security groups in (%s)", zone)

		listResp, err := instanceAPI.ListSecurityGroups(&instanceSDK.ListSecurityGroupsRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing security groups in sweeper: %s", err)
			return nil
		}

		for _, securityGroup := range listResp.SecurityGroups {
			// Can't delete default security group.
			if securityGroup.ProjectDefault {
				continue
			}
			err = instanceAPI.DeleteSecurityGroup(&instanceSDK.DeleteSecurityGroupRequest{
				Zone:            zone,
				SecurityGroupID: securityGroup.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting security groups in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepPlacementGroup(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instanceSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying the instance placement group in (%s)", zone)
		listPlacementGroups, err := instanceAPI.ListPlacementGroups(&instanceSDK.ListPlacementGroupsRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing placement groups in (%s) in sweeper: %s", zone, err)
			return nil
		}

		for _, pg := range listPlacementGroups.PlacementGroups {
			err := instanceAPI.DeletePlacementGroup(&instanceSDK.DeletePlacementGroupRequest{
				Zone:             zone,
				PlacementGroupID: pg.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting placement group in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepIP(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instanceSDK.NewAPI(scwClient)

		listIPs, err := instanceAPI.ListIPs(&instanceSDK.ListIPsRequest{Zone: zone}, scw.WithAllPages())
		if err != nil {
			logging.L.Warningf("error listing ips in (%s) in sweeper: %s", zone, err)
			return nil
		}

		for _, ip := range listIPs.IPs {
			err := instanceAPI.DeleteIP(&instanceSDK.DeleteIPRequest{
				IP:   ip.ID,
				Zone: zone,
			})
			if err != nil {
				return fmt.Errorf("error deleting ip in sweeper: %s", err)
			}
		}

		return nil
	})
}

func testSweepImage(_ string) error {
	return acctest.SweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		api := instanceSDK.NewAPI(scwClient)
		logging.L.Debugf("sweeper: destroying instance images in (%+v)", zone)

		listImagesResponse, err := api.ListImages(&instanceSDK.ListImagesRequest{
			Zone:   zone,
			Public: scw.BoolPtr(false),
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing instance images in sweeper: %w", err)
		}

		for _, image := range listImagesResponse.Images {
			err := api.DeleteImage(&instanceSDK.DeleteImageRequest{
				Zone:    zone,
				ImageID: image.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting instance image in sweeper: %w", err)
			}
		}

		return nil
	})
}
