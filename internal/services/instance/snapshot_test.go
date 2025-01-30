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

func TestAccSnapshot_BlockVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						size_in_gb = 20
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id = scaleway_instance_volume.main.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
				),
			},
		},
	})
}

func TestAccSnapshot_Unified(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "l_ssd"
						size_in_gb = 10
					}

					resource "scaleway_instance_server" "main" {
						image    = "ubuntu_jammy"
						type     = "DEV1-S"
						root_volume {
							size_in_gb = 10
							volume_type = "l_ssd"
						}
						additional_volume_ids = [
							scaleway_instance_volume.main.id
						]
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id = scaleway_instance_volume.main.id
						type = "unified"
						depends_on = [scaleway_instance_server.main]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					resource.TestCheckResourceAttr("scaleway_instance_snapshot.main", "type", "unified"),
				),
			},
		},
	})
}

func TestAccSnapshot_Server(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_focal"
						type = "DEV1-S"
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id = scaleway_instance_server.main.root_volume.0.volume_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
				),
			},
		},
	})
}

func TestAccSnapshot_ServerWithBlockVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isVolumeDestroyed(tt),
			instancechecks.IsServerDestroyed(tt),
			isSnapshotDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" main {
						type       = "b_ssd"
						size_in_gb = 10
					}

					resource "scaleway_instance_server" main {
						image = "ubuntu_focal"
						type = "DEV1-S"
						root_volume {
							size_in_gb = 10
							volume_type = "l_ssd"
						}
						additional_volume_ids = [
							scaleway_instance_volume.main.id
						]
					}

					resource "scaleway_instance_snapshot" main {
						volume_id = scaleway_instance_volume.main.id
					}`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
				),
			},
		},
	})
}

func TestAccSnapshot_RenameSnapshot(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						size_in_gb = 20
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id = scaleway_instance_volume.main.id
						name = "first_name"
						tags = ["test-terraform"]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					resource.TestCheckResourceAttr("scaleway_instance_snapshot.main", "tags.0", "test-terraform"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						size_in_gb = 20
					}

					resource "scaleway_instance_snapshot" "main" {
						volume_id = scaleway_instance_volume.main.id
						name = "second_name"
					}`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.main"),
					resource.TestCheckResourceAttr("scaleway_instance_snapshot.main", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccSnapshot_FromObject(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_object_bucket" "bucket" {
						name = "test-instance-snapshot-import-from-object"
					}

					resource "scaleway_object" "image" {
						bucket = scaleway_object_bucket.bucket.name
						key    = "image.qcow"
						file   = "testfixture/empty.qcow2"
					}

					resource "scaleway_instance_snapshot" "snapshot" {
						name = "test-instance-snapshot-import-from-object"
						type = "b_ssd"
						import {
							bucket = scaleway_object.image.bucket
							key    = scaleway_object.image.key
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					isSnapshotPresent(tt, "scaleway_instance_snapshot.snapshot"),
				),
			},
		},
	})
}

func isSnapshotPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = instanceAPI.GetSnapshot(&instanceSDK.GetSnapshotRequest{
			Zone:       zone,
			SnapshotID: ID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isSnapshotDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_snapshot" {
				continue
			}

			instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = instanceAPI.GetSnapshot(&instanceSDK.GetSnapshotRequest{
				SnapshotID: ID,
				Zone:       zone,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("snapshot (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
