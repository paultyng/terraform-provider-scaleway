package instance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceVolume_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	volumeName := "tf-volume"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isVolumeDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_volume" "test" {
						name = "%s"
						size_in_gb = 2
						type = "l_ssd"
					}`, volumeName),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_volume" "test" {
						name = "%s"
						size_in_gb = 2
						type = "l_ssd"
					}
					
					data "scaleway_instance_volume" "test" {
						name = "${scaleway_instance_volume.test.name}"
					}
					
					data "scaleway_instance_volume" "test2" {
						volume_id = "${scaleway_instance_volume.test.id}"
					}
				`, volumeName),
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "data.scaleway_instance_volume.test"),
					resource.TestCheckResourceAttr("data.scaleway_instance_volume.test", "size_in_gb", "2"),
					resource.TestCheckResourceAttr("data.scaleway_instance_volume.test", "type", "l_ssd"),
				),
			},
		},
	})
}
