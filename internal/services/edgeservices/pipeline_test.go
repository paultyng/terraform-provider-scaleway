package edgeservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	edgeservicestestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/edgeservices/testfuncs"
)

func TestAccEdgeServicesPipeline_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      edgeservicestestfuncs.CheckEdgeServicesCacheDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_edge_services_pipeline" "main" {
					  name        = "tf-tests-pipeline-name"
					  description = "a description"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					edgeservicestestfuncs.CheckEdgeServicesPipelineExists(tt, "scaleway_edge_services_pipeline.main"),
					resource.TestCheckResourceAttr("scaleway_edge_services_pipeline.main", "name", "tf-tests-pipeline-name"),
					resource.TestCheckResourceAttr("scaleway_edge_services_pipeline.main", "description", "a description"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_pipeline.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_pipeline.main", "updated_at"),
				),
			},
		},
	})
}
