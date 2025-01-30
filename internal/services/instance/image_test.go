package instance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccImage_BlockVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isImageDestroyed(tt),
			isSnapshotDestroyed(tt),
			isVolumeDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       	= "b_ssd"
						size_in_gb 	= 20
						zone 		= "nl-ams-1"
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id 	= scaleway_instance_volume.main.id
						zone 		= "nl-ams-1"
					}

					resource "scaleway_instance_image" "main" {
						name 			= "test_image_basic"
						root_volume_id 	= scaleway_instance_snapshot.main.id
						tags 			= ["tag1", "tag2", "tag3"]
						zone 			= "nl-ams-1"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.main"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "name", "test_image_basic"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "zone", "nl-ams-1"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "public", "false"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "creation_date"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "modification_date"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "modification_date", "scaleway_instance_image.main", "creation_date"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "state"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "organization_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "project_id"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       	= "b_ssd"
						size_in_gb 	= 20
						zone 		= "nl-ams-1"
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id 	= scaleway_instance_volume.main.id
						zone 		= "nl-ams-1"
					}

					resource "scaleway_instance_image" "main" {
						name 			= "test_image_renamed"
						root_volume_id 	= scaleway_instance_snapshot.main.id
						tags 			= ["new tag"]
						public 			= true
						architecture	= "arm"
						zone 			= "nl-ams-1"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.main"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "name", "test_image_renamed"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.#", "1"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.0", "new tag"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "public", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "arm"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "zone", "nl-ams-1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "creation_date"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "modification_date"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "state"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "organization_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "project_id"),
				),
			},
		},
	})
}

func TestAccImage_ExternalBlockVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isImageDestroyed(tt),
			isSnapshotDestroyed(tt),
			isVolumeDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_block_volume" "main" {
						size_in_gb = 50
						iops = 5000
					}

					resource "scaleway_block_snapshot" "main" {
						volume_id = scaleway_block_volume.main.id
					}
				`,
			},
			{
				Config: `
					resource "scaleway_block_volume" "main" {
						size_in_gb = 50
						iops = 5000
					}

					resource "scaleway_block_volume" "additional1" {
						size_in_gb = 50
						iops = 5000
					}

					resource "scaleway_block_snapshot" "main" {
						volume_id = scaleway_block_volume.main.id
					}

					resource "scaleway_block_snapshot" "additional1" {
						volume_id = scaleway_block_volume.additional1.id
					}

					resource "scaleway_instance_image" "main" {
						name 			= "tf-test-image-external-block-volume"
						root_volume_id 	= scaleway_block_snapshot.main.id
						additional_volume_ids = [scaleway_block_snapshot.additional1.id]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_block_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volume_ids.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "additional_volume_ids.0", "scaleway_block_snapshot.additional1", "id"),
				),
			},
			{
				Config: `
					resource "scaleway_block_volume" "main" {
						size_in_gb = 50
						iops = 5000
					}

					resource "scaleway_block_volume" "additional1" {
						size_in_gb = 50
						iops = 5000
					}

					resource "scaleway_block_snapshot" "main" {
						volume_id = scaleway_block_volume.main.id
					}

					resource "scaleway_block_snapshot" "additional1" {
						volume_id = scaleway_block_volume.additional1.id
					}

					resource "scaleway_instance_image" "main" {
						name 			= "tf-test-image-external-block-volume"
						root_volume_id 	= scaleway_block_snapshot.main.id
						additional_volume_ids = []
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_block_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volume_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccImage_Server(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isImageDestroyed(tt),
			isSnapshotDestroyed(tt),
			instancechecks.IsServerDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image 	= "ubuntu_focal"
						type 	= "DEV1-S"
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id	= scaleway_instance_server.main.root_volume.0.volume_id
						depends_on 	= [ scaleway_instance_server.main ]
					}

					resource "scaleway_instance_image" "main" {
						root_volume_id 			= scaleway_instance_snapshot.main.id
						tags 					= ["test_remove_tags"]
						depends_on 				= [ scaleway_instance_snapshot.main ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.0", "test_remove_tags"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "name"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "public", "false"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image 	= "ubuntu_focal"
						type 	= "DEV1-S"
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id 	= scaleway_instance_server.main.root_volume.0.volume_id
						depends_on	= [ scaleway_instance_server.main ]
					}

					resource "scaleway_instance_image" "main" {
						root_volume_id 			= scaleway_instance_snapshot.main.id
						depends_on 				= [ scaleway_instance_snapshot.main ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.main", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "tags.#", "0"),
					resource.TestCheckNoResourceAttr("scaleway_instance_image.main", "tags.0"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "name"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "public", "false"),
				),
			},
		},
	})
}

func TestAccImage_ServerWithBlockVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isImageDestroyed(tt),
			isSnapshotDestroyed(tt),
			isVolumeDestroyed(tt),
			instancechecks.IsServerDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "block01" {
						type       = "b_ssd"
						size_in_gb = 21
					}
					resource "scaleway_instance_snapshot" "block01" {
						volume_id	= scaleway_instance_volume.block01.id
						depends_on 	= [ scaleway_instance_volume.block01 ]
					}

					resource "scaleway_instance_server" "server" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.block01"),
					isServerPresent(tt, "scaleway_instance_server.server"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.block01"),
				),
			},

			{
				Config: `
					resource "scaleway_instance_volume" "block01" {
						type       = "b_ssd"
						size_in_gb = 21
					}
					resource "scaleway_instance_snapshot" "block01" {
						volume_id	= scaleway_instance_volume.block01.id
						depends_on 	= [ scaleway_instance_volume.block01 ]
					}

					resource "scaleway_instance_server" "server" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
					}
					resource "scaleway_instance_snapshot" "server" {
						volume_id 	= scaleway_instance_server.server.root_volume.0.volume_id
						depends_on 	= [ scaleway_instance_server.server ]
					}

					resource "scaleway_instance_image" "main" {
						root_volume_id 	= scaleway_instance_snapshot.server.id
						additional_volume_ids = [
							scaleway_instance_snapshot.block01.id
						]
						depends_on = [
							scaleway_instance_snapshot.block01,
							scaleway_instance_snapshot.server,
						]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.block01"),
					isServerPresent(tt, "scaleway_instance_server.server"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.block01"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.server"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.server", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "additional_volumes.0.id", "scaleway_instance_snapshot.block01", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volumes.0.volume_type", "b_ssd"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volumes.0.size", "21000000000"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "name"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "public", "false"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_volume" "block01" {
						type       = "b_ssd"
						size_in_gb = 21
					}
					resource "scaleway_instance_snapshot" "block01" {
						volume_id	= scaleway_instance_volume.block01.id
						depends_on 	= [ scaleway_instance_volume.block01 ]
					}

					resource "scaleway_instance_volume" "block02" {
						type       = "b_ssd"
						size_in_gb = 22
					}
					resource "scaleway_instance_snapshot" "block02" {
						volume_id	= scaleway_instance_volume.block02.id
						depends_on 	= [ scaleway_instance_volume.block02 ]
					}

					resource "scaleway_instance_server" "server" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
					}
					resource "scaleway_instance_snapshot" "server" {
						volume_id 	= scaleway_instance_server.server.root_volume.0.volume_id
						depends_on 	= [ scaleway_instance_server.server ]
					}

					resource "scaleway_instance_image" "main" {
						root_volume_id 	= scaleway_instance_snapshot.server.id
						additional_volume_ids = [
							scaleway_instance_snapshot.block02.id,
						]
						depends_on = [
							scaleway_instance_snapshot.block02,
							scaleway_instance_snapshot.server,
						]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.block01"),
					isVolumePresent(tt, "scaleway_instance_volume.block02"),
					isServerPresent(tt, "scaleway_instance_server.server"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.block01"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.block02"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.server"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.server", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "additional_volumes.0.id", "scaleway_instance_snapshot.block02", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volumes.0.volume_type", "b_ssd"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volumes.0.size", "22000000000"),
					resource.TestCheckResourceAttrSet("scaleway_instance_image.main", "name"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "public", "false"),
				),
			},
		},
	})
}

func TestAccImage_ServerWithLocalVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isImageDestroyed(tt),
			isSnapshotDestroyed(tt),
			isVolumeDestroyed(tt),
			instancechecks.IsServerDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "server01" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
						root_volume {
							size_in_gb = 15
							volume_type = "l_ssd"
						}
					}
					resource "scaleway_instance_snapshot" "local01" {
						volume_id = scaleway_instance_server.server01.root_volume.0.volume_id
						depends_on = [ scaleway_instance_server.server01 ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.server01"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.local01"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "server01" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
						root_volume {
							size_in_gb = 15
							volume_type = "l_ssd"
						}
					}
					resource "scaleway_instance_server" "server02" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
						root_volume {
							size_in_gb = 10
							volume_type = "l_ssd"
						}
					}

					resource "scaleway_instance_snapshot" "local01" {
						volume_id = scaleway_instance_server.server01.root_volume.0.volume_id
						depends_on = [ scaleway_instance_server.server01 ]
					}
					resource "scaleway_instance_snapshot" "local02" {
						volume_id = scaleway_instance_server.server02.root_volume.0.volume_id
						depends_on = [ scaleway_instance_server.server02 ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.server01"),
					isServerPresent(tt, "scaleway_instance_server.server02"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.local01"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.local02"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "server01" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
						root_volume {
							size_in_gb = 15
							volume_type = "l_ssd"
						}
					}
					resource "scaleway_instance_server" "server02" {
						image	= "ubuntu_focal"
						type 	= "DEV1-S"
						root_volume {
							size_in_gb = 10
							volume_type = "l_ssd"
						}
					}

					resource "scaleway_instance_snapshot" "local01" {
						volume_id = scaleway_instance_server.server01.root_volume.0.volume_id
						depends_on = [ scaleway_instance_server.server01 ]
					}
					resource "scaleway_instance_snapshot" "local02" {
						volume_id = scaleway_instance_server.server02.root_volume.0.volume_id
						depends_on = [ scaleway_instance_server.server02 ]
					}

					resource "scaleway_instance_image" "main" {
						root_volume_id 	= scaleway_instance_snapshot.local01.id
						additional_volume_ids = [ scaleway_instance_snapshot.local02.id ]
						depends_on = [
							scaleway_instance_snapshot.local01,
							scaleway_instance_snapshot.local02,
						]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.server01"),
					isServerPresent(tt, "scaleway_instance_server.server02"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.local01"),
					isSnapshotPresent(tt, "scaleway_instance_snapshot.local02"),
					instancechecks.DoesImageExists(tt, "scaleway_instance_image.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "root_volume_id", "scaleway_instance_snapshot.local01", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_image.main", "additional_volumes.0.id", "scaleway_instance_snapshot.local02", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volumes.0.volume_type", "l_ssd"),
					resource.TestCheckResourceAttr("scaleway_instance_image.main", "additional_volumes.0.size", "10000000000"),
				),
			},
		},
	})
}

func isImageDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_image" {
				continue
			}
			instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}
			_, err = instanceAPI.GetImage(&instanceSDK.GetImageRequest{
				ImageID: ID,
				Zone:    zone,
			})
			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("image (%s) still exists", rs.Primary.ID)
			}
			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
