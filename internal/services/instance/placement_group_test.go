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

func TestAccPlacementGroup_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isPlacementGroupDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_placement_group" "base" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.base"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_mode", "optional"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_type", "max_availability"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_placement_group" "base" {
						policy_mode = "enforced"
						policy_type = "low_latency"
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.base"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_mode", "enforced"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_type", "low_latency"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_respected", "true"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_placement_group" "base" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.base"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_mode", "optional"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "policy_type", "max_availability"),
				),
			},
		},
	})
}

func TestAccPlacementGroup_Rename(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isPlacementGroupDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_placement_group" "base" {
						name        = "foo"
						policy_mode = "enforced"
						policy_type = "low_latency"
					}

					resource "scaleway_instance_server" "base" {
						type  = "DEV1-S"
						image = "ubuntu_focal"
						placement_group_id = "${scaleway_instance_placement_group.base.id}"
					}`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.base"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "name", "foo"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_placement_group" "base" {
						name        = "bar"
						policy_mode = "enforced"
						policy_type = "low_latency"
					}`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.base"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.base", "name", "bar"),
				),
			},
		},
	})
}

func TestAccPlacementGroup_Tags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_instance_placement_group" "main" {}
					`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.main"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.main", "tags.#", "0"),
				),
			},
			{
				Config: `
						resource "scaleway_instance_placement_group" "main" {
							tags = ["foo", "bar"]
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.main"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.main", "tags.0", "foo"),
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.main", "tags.1", "bar"),
				),
			},
			{
				Config: `
						resource "scaleway_instance_placement_group" "main" {
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_placement_group.main", "tags.#", "0"),
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.main"),
				),
			},
		},
	})
}

func isPlacementGroupPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = instanceAPI.GetPlacementGroup(&instanceSDK.GetPlacementGroupRequest{
			Zone:             zone,
			PlacementGroupID: ID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isPlacementGroupDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_instance_placement_group" {
				continue
			}

			instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = instanceAPI.GetPlacementGroup(&instanceSDK.GetPlacementGroupRequest{
				Zone:             zone,
				PlacementGroupID: ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("placement group (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}
		return nil
	}
}
