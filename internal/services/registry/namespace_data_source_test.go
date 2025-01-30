package registry_test

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
					resource "scaleway_registry_namespace" "reg00" {
						region = "pl-waw"
						name = "test-cr-data"
					}
					
					data "scaleway_registry_namespace" "regData01" {
						region = "pl-waw"
						name = scaleway_registry_namespace.reg00.name
					}
					
					data "scaleway_registry_namespace" "regData02" {
						region = "pl-waw"
						namespace_id = scaleway_registry_namespace.reg00.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isNamespacePresent(tt, "scaleway_registry_namespace.reg00"),
					resource.TestCheckResourceAttr("scaleway_registry_namespace.reg00", "name", "test-cr-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData01", "id"),
					resource.TestCheckResourceAttr("data.scaleway_registry_namespace.regData01", "is_public", "false"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData01", "endpoint"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData01", "namespace_id"),

					resource.TestCheckResourceAttr("data.scaleway_registry_namespace.regData02", "name", "test-cr-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData02", "id"),
					resource.TestCheckResourceAttr("data.scaleway_registry_namespace.regData02", "is_public", "false"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData02", "endpoint"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData02", "namespace_id"),
				),
			},
			{
				Config: `
					resource "scaleway_registry_namespace" "reg00" {
						region = "pl-waw"
						name = "test-cr-data"
						is_public = "true"
					}
					
					data "scaleway_registry_namespace" "regData01" {
						region = "pl-waw"
						name = scaleway_registry_namespace.reg00.name
					}
					
					data "scaleway_registry_namespace" "regData02" {
						region = "pl-waw"
						namespace_id = scaleway_registry_namespace.reg00.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isNamespacePresent(tt, "scaleway_registry_namespace.reg00"),

					resource.TestCheckResourceAttr("scaleway_registry_namespace.reg00", "name", "test-cr-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData01", "id"),
					resource.TestCheckResourceAttr("data.scaleway_registry_namespace.regData01", "is_public", "true"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData01", "endpoint"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData01", "namespace_id"),

					resource.TestCheckResourceAttr("data.scaleway_registry_namespace.regData02", "name", "test-cr-data"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData02", "id"),
					resource.TestCheckResourceAttr("data.scaleway_registry_namespace.regData02", "is_public", "true"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData02", "endpoint"),
					resource.TestCheckResourceAttrSet("data.scaleway_registry_namespace.regData02", "namespace_id"),
				),
			},
		},
	})
}
