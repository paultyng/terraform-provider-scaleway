package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceGroup_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamGroupDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_iam_group" "main_ds_basic" {
					  name = "tf_test_data_source_basic"
					}
					
					data "scaleway_iam_group" "find_by_id_basic" {
					  group_id        = scaleway_iam_group.main_ds_basic.id
					}
					
					data "scaleway_iam_group" "find_by_name_basic" {
					  name            = scaleway_iam_group.main_ds_basic.name
					  organization_id = "105bdce1-64c0-48ab-899d-868455867ecf"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamGroupExists(tt, "scaleway_iam_group.main_ds_basic"),
					resource.TestCheckResourceAttr("data.scaleway_iam_group.find_by_id_basic", "name", "tf_test_data_source_basic"),
					resource.TestCheckResourceAttr("data.scaleway_iam_group.find_by_name_basic", "name", "tf_test_data_source_basic"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_group.find_by_id_basic", "id", "scaleway_iam_group.main_ds_basic", "id"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_group.find_by_name_basic", "id", "scaleway_iam_group.main_ds_basic", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceGroup_UsersAndApplications(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIamGroupDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_iam_application" "app00" {
					  name = "tf_tests_iam_group_ds_app"
					}
					
					data "scaleway_iam_user" "user00" {
					  user_id         = "ef29ce05-3f2b-4fa0-a259-d76110850d57"
					}
					data "scaleway_iam_user" "user01" {
					  user_id         = "84d20ae1-9650-419a-ab74-7ab09b6262e0"
					}
					
					resource "scaleway_iam_group" "main_ds_mix" {
					  name = "tf_test_data_source_mix"
					  application_ids = [
						scaleway_iam_application.app00.id,
					  ]
					  user_ids = [
						data.scaleway_iam_user.user00.user_id,
						data.scaleway_iam_user.user01.user_id,
					  ]
					}
					
					data "scaleway_iam_group" "find_by_id_mix" {
					  group_id        = scaleway_iam_group.main_ds_mix.id
					  organization_id = "105bdce1-64c0-48ab-899d-868455867ecf"
					}
					
					data "scaleway_iam_group" "find_by_name_mix" {
					  name            = scaleway_iam_group.main_ds_mix.name
					  organization_id = "105bdce1-64c0-48ab-899d-868455867ecf"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIamGroupExists(tt, "scaleway_iam_group.main_ds_mix"),
					resource.TestCheckResourceAttr("data.scaleway_iam_group.find_by_id_mix", "name", "tf_test_data_source_mix"),
					resource.TestCheckResourceAttr("data.scaleway_iam_group.find_by_name_mix", "name", "tf_test_data_source_mix"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_group.find_by_id_mix", "id", "scaleway_iam_group.main_ds_mix", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.scaleway_iam_group.find_by_name_mix", "id", "scaleway_iam_group.main_ds_mix", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.scaleway_iam_group.find_by_id_mix", "user_ids.*", "data.scaleway_iam_user.user00", "user_id"),
					resource.TestCheckTypeSetElemAttrPair("data.scaleway_iam_group.find_by_id_mix", "user_ids.*", "data.scaleway_iam_user.user01", "user_id"),
					resource.TestCheckTypeSetElemAttrPair("data.scaleway_iam_group.find_by_name_mix", "user_ids.*", "data.scaleway_iam_user.user00", "user_id"),
					resource.TestCheckTypeSetElemAttrPair("data.scaleway_iam_group.find_by_name_mix", "user_ids.*", "data.scaleway_iam_user.user01", "user_id"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_group.find_by_id_mix", "application_ids.0", "scaleway_iam_application.app00", "id"),
					resource.TestCheckResourceAttrPair("data.scaleway_iam_group.find_by_name_mix", "application_ids.0", "scaleway_iam_application.app00", "id"),
				),
			},
		},
	})
}
