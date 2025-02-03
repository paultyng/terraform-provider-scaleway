package edgeservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	edgeservicestestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/edgeservices/testfuncs"
)

func TestAccEdgeServicesCache_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      edgeservicestestfuncs.CheckEdgeServicesCacheDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_edge_services_cache_stage" "main" {}
				`,
				Check: resource.ComposeTestCheckFunc(
					edgeservicestestfuncs.CheckEdgeServicesCacheExists(tt, "scaleway_edge_services_cache_stage.main"),
					resource.TestCheckResourceAttr("scaleway_edge_services_cache_stage.main", "fallback_ttl", "3600"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_cache_stage.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_cache_stage.main", "updated_at"),
				),
			},
		},
	})
}
