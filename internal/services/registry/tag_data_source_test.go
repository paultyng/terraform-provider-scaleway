package registry_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry"
	registrytestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry/testfuncs"
)

func TestAccDataSourceImageTag_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	expectedTagName := "test"
	namespaceName := "test-namespace-2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isNamespaceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_registry_namespace" "test" {
						name        = "%s"
						description = "Test namespace for Docker image"
						is_public   = false
					}
				`, namespaceName),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["scaleway_registry_namespace.test"]
						if !ok {
							return errors.New("not found: scaleway_registry_namespace.test")
						}

						endpoint := rs.Primary.Attributes["endpoint"]
						if endpoint == "" {
							return errors.New("no endpoint found for scaleway_registry_namespace.test")
						}

						return registrytestfuncs.PushImageToRegistry(tt, endpoint, expectedTagName)(s)
					},
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_registry_namespace" "test" {
						name        = "%s"
						description = "Test namespace for Docker image"
						is_public   = false
					}

					data "scaleway_registry_namespace" "test" {
						name = "%s"
					}

					data "scaleway_registry_image" "image" {
  						name = "alpine"
					}

					data "scaleway_registry_image_tag" "tag" {
						name  = "%s"
						image_id = "${data.scaleway_registry_image.image.id}"
					}
				`, namespaceName, namespaceName, expectedTagName),
				Check: resource.ComposeTestCheckFunc(
					isTagPresent(tt, "data.scaleway_registry_image_tag.tag"),
					isNamespacePresent(tt, "scaleway_registry_namespace.test"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image_tag.tag", "name", expectedTagName),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_image_tag.tag", "digest"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_image_tag.tag", "created_at"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_image_tag.tag", "updated_at"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_registry_namespace" "test" {
						name        = "%s"
						description = "Test namespace for Docker image"
						is_public   = false
					}

					data "scaleway_registry_namespace" "test" {
						name = "%s"
					}

					data "scaleway_registry_image" "image" {
  						name = "alpine"
					}
				`, namespaceName, namespaceName),
				Check: resource.ComposeTestCheckFunc(
					deleteImage(tt, "data.scaleway_registry_image.image"),
				),
			},
		},
	})
}

func isTagPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, _, err := registry.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetTag(&registrySDK.GetTagRequest{
			TagID:  locality.ExpandID(rs.Primary.ID),
			Region: region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func deleteImage(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}
		api, region, _, err := registry.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}
		_, err = api.DeleteImage(&registrySDK.DeleteImageRequest{
			Region:  region,
			ImageID: locality.ExpandID(rs.Primary.ID),
		})
		if err != nil {
			return err
		}

		return nil
	}
}
