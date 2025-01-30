package container_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceContainer_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isNamespaceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_container_namespace main {
					}

					resource scaleway_container main {
						name = "test-container-data"
						namespace_id = scaleway_container_namespace.main.id
					}

					data "scaleway_container" "by_name" {
						namespace_id = scaleway_container_namespace.main.id
						name = scaleway_container.main.name
					}
					
					data "scaleway_container" "by_id" {
						namespace_id = scaleway_container_namespace.main.id
						container_id = scaleway_container.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isContainerPresent(tt, "scaleway_container.main"),
					resource.TestCheckResourceAttr("scaleway_container.main", "name", "test-container-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_container.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.scaleway_container.by_name", "name"),

					resource.TestCheckResourceAttr("data.scaleway_container.by_id", "name", "test-container-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_container.by_id", "id"),
				),
			},
		},
	})
}
