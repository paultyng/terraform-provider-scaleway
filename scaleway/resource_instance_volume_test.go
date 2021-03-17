package scaleway

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	resource.AddTestSweepers("scaleway_instance_volume", &resource.Sweeper{
		Name: "scaleway_instance_volume",
		F:    testSweepComputeInstanceVolume,
	})
}

func testSweepComputeInstanceVolume(_ string) error {
	return sweepZones(scw.AllZones, func(scwClient *scw.Client, zone scw.Zone) error {
		instanceAPI := instance.NewAPI(scwClient)
		l.Debugf("sweeper: destroying the volumes in (%s)", zone)

		listVolumesResponse, err := instanceAPI.ListVolumes(&instance.ListVolumesRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing volumes in sweeper: %s", err)
		}

		for _, volume := range listVolumesResponse.Volumes {
			if volume.Server == nil {
				err := instanceAPI.DeleteVolume(&instance.DeleteVolumeRequest{
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

func TestAccScalewayInstanceVolume_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceVolumeDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "test" {
						type       = "l_ssd"
						size_in_gb = 20
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayInstanceVolumeExists(tt, "scaleway_instance_volume.test"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "size_in_gb", "20"),
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
					testAccCheckScalewayInstanceVolumeExists(tt, "scaleway_instance_volume.test"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "name", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.test", "size_in_gb", "20"),
				),
			},
		},
	})
}

func TestAccScalewayInstanceVolume_FromVolume(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceVolumeDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "test1" {
						type       = "l_ssd"
						size_in_gb = 20
					}
			
					resource "scaleway_instance_volume" "test2" {
						type           = "l_ssd"
						from_volume_id = "${scaleway_instance_volume.test1.id}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayInstanceVolumeExists(tt, "scaleway_instance_volume.test1"),
					testAccCheckScalewayInstanceVolumeExists(tt, "scaleway_instance_volume.test2"),
				),
			},
		},
	})
}

func TestAccScalewayInstanceVolume_DifferentNameGenerated(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceVolumeDestroy(tt),
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

func TestAccScalewayInstanceVolume_ResizeBlock(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceVolumeDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						size_in_gb = 20
					}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayInstanceVolumeExists(tt, "scaleway_instance_volume.main"),
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
					testAccCheckScalewayInstanceVolumeExists(tt, "scaleway_instance_volume.main"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.main", "size_in_gb", "30"),
				),
			},
		},
	})
}

func TestAccScalewayInstanceVolume_CannotResizeBlockDown(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceVolumeDestroy(tt),
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

func testAccCheckScalewayInstanceVolumeExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		zone, id, err := parseZonedID(rs.Primary.ID)
		if err != nil {
			return err
		}

		instanceAPI := instance.NewAPI(tt.Meta.scwClient)
		_, err = instanceAPI.GetVolume(&instance.GetVolumeRequest{
			VolumeID: id,
			Zone:     zone,
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayInstanceVolumeDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		instanceAPI := instance.NewAPI(tt.Meta.scwClient)
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_volume" {
				continue
			}

			zone, id, err := parseZonedID(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = instanceAPI.GetVolume(&instance.GetVolumeRequest{
				Zone:     zone,
				VolumeID: id,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("volume (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !is404Error(err) {
				return err
			}
		}
		return nil
	}
}

func TestAccScalewayInstanceVolume_ZoneDiffersProviderZone(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceVolumeDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "scaleway" {
						zone = "fr-par-1"
						region = "fr-par"
					}

					resource "scaleway_instance_volume" "extended_disk" {
						type             = "b_ssd"
						name             = "test-instance-disk-1"
						size_in_gb       = 1
						zone             = "pl-waw-1"
					}

					resource "scaleway_instance_server" "main" {
						name  = "test-instance"
						type  = "DEV1-M"
						image = "ubuntu_focal"
						zone  = "pl-waw-1"
					
						root_volume {
							size_in_gb            = 40
							delete_on_termination = true
						}
						additional_volume_ids   = scaleway_instance_volume.extended_disk.*.id 
					}`,
			},
		},
	})
}