package edgeservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	edgeservicestestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/edgeservices/testfuncs"
)

func TestAccEdgeServicesBackend_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      edgeservicestestfuncs.CheckEdgeServicesBackendDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_object_bucket" "main" {
						name = "test-acc-scaleway-object-bucket-basic-es"
						tags = {
							foo = "bar"
						}
					}

					resource "scaleway_edge_services_backend_stage" "main" {
					  s3_backend_config {
						bucket_name   = scaleway_object_bucket.main.name
						bucket_region = "fr-par"
					  }
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					edgeservicestestfuncs.CheckEdgeServicesBackendExists(tt, "scaleway_edge_services_backend_stage.main"),
					resource.TestCheckResourceAttr("scaleway_edge_services_backend_stage.main", "s3_backend_config.0.is_website", "false"),
					resource.TestCheckResourceAttr("scaleway_edge_services_backend_stage.main", "s3_backend_config.0.bucket_region", "fr-par"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_backend_stage.main", "s3_backend_config.0.bucket_name"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_backend_stage.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_backend_stage.main", "updated_at"),
				),
			},
		},
	})
}
