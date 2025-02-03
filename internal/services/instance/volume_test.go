package instance_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
)

func TestAccVolume_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "test" {
						type       = "l_ssd"
						size_in_gb = 20
						tags = ["test-terraform"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.test"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "size_in_gb", "20"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "tags.0", "test-terraform"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_volume" "test" {
						type       = "l_ssd"
						name       = "terraform-test"
						size_in_gb = 20
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.test"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "name", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "size_in_gb", "20"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccVolume_DifferentNameGenerated(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "test" {
						type       = "l_ssd"
						size_in_gb = 20
					}
				`,
			},
			{
				Config: `
					resource "scaleway_instance_volume" "test" {
						type       = "l_ssd"
						size_in_gb = 20
					}
				`,
			},
		},
	})
}

func TestAccVolume_ResizeBlock(t *testing.T) {
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
					}`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.main"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.main", "size_in_gb", "20"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						size_in_gb = 30
					}`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.main"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.main", "size_in_gb", "30"),
				),
			},
		},
	})
}

func TestAccVolume_ResizeNotBlock(t *testing.T) {
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
						type       = "l_ssd"
						size_in_gb = 20
					}`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.main"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.main", "size_in_gb", "20"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "l_ssd"
						size_in_gb = 30
					}`,
				ExpectError: regexp.MustCompile("only block volume can be resized"),
			},
		},
	})
}

func TestAccVolume_CannotResizeBlockDown(t *testing.T) {
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
					}`,
			},
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						size_in_gb = 10
					}`,
				ExpectError: regexp.MustCompile("block volumes cannot be resized down"),
			},
		},
	})
}

func TestAccVolume_Scratch(t *testing.T) {
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
						type       = "scratch"
						size_in_gb = 20
					}`,
			},
		},
	})
}

func isVolumePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		zone, id, err := zonal.ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		instanceAPI := instanceSDK.NewAPI(tt.Meta.ScwClient())
		_, err = instanceAPI.GetVolume(&instanceSDK.GetVolumeRequest{
			VolumeID: id,
			Zone:     zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isVolumeDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instanceAPI := instanceSDK.NewAPI(tt.Meta.ScwClient())
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_volume" {
				continue
			}

			zone, id, err := zonal.ParseID(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = instanceAPI.GetVolume(&instanceSDK.GetVolumeRequest{
				Zone:     zone,
				VolumeID: id,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("volume (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}
		return nil
	}
}
