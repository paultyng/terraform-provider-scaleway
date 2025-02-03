package container_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceNamespace_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isNamespaceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_container_namespace" "main" {
						name = "test-cr-data"
					}
					
					data "scaleway_container_namespace" "by_name" {
						name = scaleway_container_namespace.main.name
					}
					
					data "scaleway_container_namespace" "by_id" {
						namespace_id = scaleway_container_namespace.main.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isNamespacePresent(tt, "scaleway_container_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "name", "test-cr-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_container_namespace.by_name", "id"),

					resource.TestCheckResourceAttr("data.scaleway_container_namespace.by_id", "name", "test-cr-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_container_namespace.by_id", "id"),
				),
			},
		},
	})
}
