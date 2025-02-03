package edgeservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	edgeservicestestfuncs "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/edgeservices/testfuncs"
)

func TestAccEdgeServicesDNS_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      edgeservicestestfuncs.CheckEdgeServicesDNSDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_edge_services_dns_stage" "main" {
					  fqdns        = ["subodomain.example.fr"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					edgeservicestestfuncs.CheckEdgeServicesDNSExists(tt, "scaleway_edge_services_dns_stage.main"),
					resource.TestCheckResourceAttr("scaleway_edge_services_dns_stage.main", "fqdns.0", "subodomain.example.fr"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_dns_stage.main", "type"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_dns_stage.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_edge_services_dns_stage.main", "updated_at"),
				),
			},
		},
	})
}
