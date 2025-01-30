package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccDataSourceImage_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_instance_image" "test" {
						image_id = "fr-par-1/cf44b8f5-77e2-42ed-8f1e-09ed5bb028fc"
					}

					data "scaleway_instance_image" "test2" {
						image_id = "cf44b8f5-77e2-42ed-8f1e-09ed5bb028fc"
					}`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.DoesImageExists(tt, "data.scaleway_instance_image.test"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test", "name", "Ubuntu 20.04 Focal Fossa"),

					instancechecks.DoesImageExists(tt, "data.scaleway_instance_image.test2"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "image_id", "fr-par-1/cf44b8f5-77e2-42ed-8f1e-09ed5bb028fc"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "name", "Ubuntu 20.04 Focal Fossa"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "creation_date", "2020-07-09T10:38:16Z"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "modification_date", "2020-07-09T10:38:16Z"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "latest", "true"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "public", "true"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "from_server_id", ""),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "state", "available"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "root_volume_id", "6e66445c-e52e-4cfa-bf4c-f36e291e2c30"),
					resource.TestCheckResourceAttr("data.scaleway_instance_image.test2", "additional_volume_ids.#", "0"),
				),
			},
		},
	})
}
