package registry_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry"
)

func TestAccDataSourceImage_Basic(t *testing.T) {
	t.Skip("It is difficult to test this datasource as we cannot create registry images with Terraform.")
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	ubuntuImageID := "4b5a47c0-6fbf-4388-8783-c07c28d3c2eb"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isNamespaceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_registry_image" "ubuntu" {
						image_id = "` + ubuntuImageID + `"
					}

					data "scaleway_registry_image" "ubuntu2" {
						name = "ubuntu"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isImagePresent(tt, "data.scaleway_registry_image.ubuntu"),

					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "name", "ubuntu"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "id", "fr-par/"+ubuntuImageID),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "image_id", ubuntuImageID),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_image.ubuntu", "size"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "visibility", "inherit"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "tags.0", "bar"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "tags.1", "foo"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu", "tags.2", "latest"),

					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "name", "ubuntu"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "id", "fr-par/"+ubuntuImageID),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "image_id", ubuntuImageID),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_image.ubuntu2", "size"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "visibility", "inherit"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "tags.0", "bar"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "tags.1", "foo"),
					resource.TestCheckResourceAttr("data.scaleway_registry_image.ubuntu2", "tags.2", "latest"),
				),
			},
		},
	})
}

func isImagePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := registry.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetImage(&registrySDK.GetImageRequest{
			ImageID: id,
			Region:  region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
