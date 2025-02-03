package instancetestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
)

func CheckIPExists(tt *acctest.TestTools, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = instanceAPI.GetIP(&instanceSDK.GetIPRequest{
			IP:   ID,
			Zone: zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func IsServerDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_server" {
				continue
			}

			instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = instanceAPI.GetServer(&instanceSDK.GetServerRequest{
				ServerID: ID,
				Zone:     zone,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("server (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func IsServerRootVolumeDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_server" {
				continue
			}

			localizedRootVolumeID, exists := rs.Primary.Attributes["root_volume.0.volume_id"]
			if !exists {
				return fmt.Errorf("root_volume ID not found in resource %s", rs.Primary.ID)
			}

			zone, _, err := locality.ParseLocalizedID(rs.Primary.ID)
			if err != nil {
				return err
			}
			rootVolumeID := locality.ExpandID(localizedRootVolumeID)

			api := instance.NewBlockAndInstanceAPI(meta.ExtractScwClient(tt.Meta))

			_, err = api.GetUnknownVolume(&instance.GetUnknownVolumeRequest{
				VolumeID: rootVolumeID,
				Zone:     scw.Zone(zone),
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("server's root volume (%s) still exists", rootVolumeID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func IsIPDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "scaleway_instance_ip" {
				continue
			}

			instanceAPI, zone, id, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, errIP := instanceAPI.GetIP(&instanceSDK.GetIPRequest{
				Zone: zone,
				IP:   id,
			})

			// If no error resource still exist
			if errIP == nil {
				return fmt.Errorf("resource %s(%s) still exist", rs.Type, rs.Primary.ID)
			}

			// Unexpected api error we return it
			// We check for 403 because instanceSDK API return 403 for deleted IP
			if !httperrors.Is404(errIP) && !httperrors.Is403(errIP) {
				return errIP
			}
		}

		return nil
	}
}

func DoesImageExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		zone, ID, err := zonal.ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}
		instanceAPI := instanceSDK.NewAPI(tt.Meta.ScwClient())
		_, err = instanceAPI.GetImage(&instanceSDK.GetImageRequest{
			ImageID: ID,
			Zone:    zone,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
