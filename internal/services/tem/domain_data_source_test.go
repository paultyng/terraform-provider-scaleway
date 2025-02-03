package tem_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceDomain_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName := "terraform-ds.test.local"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isDomainDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_tem_domain" "main" {
						name 	   = "%s"
						accept_tos = true
					}
					
					data "scaleway_tem_domain" "prod" {
						name = "${scaleway_tem_domain.main.name}"
					}
					
					data "scaleway_tem_domain" "stg" {
						domain_id = "${scaleway_tem_domain.main.id}"
					}
				`, domainName),
				Check: resource.ComposeTestCheckFunc(
					isDomainPresent(tt, "data.scaleway_tem_domain.prod"),
					resource.TestCheckResourceAttr("data.scaleway_tem_domain.prod", "name", domainName),

					isDomainPresent(tt, "data.scaleway_tem_domain.stg"),
					resource.TestCheckResourceAttr("data.scaleway_tem_domain.stg", "name", domainName),
				),
			},
		},
	})
}

func TestAccDataSourceDomain_Reputation(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName := "test.scaleway-terraform.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_tem_domain" "test" {
						name = "%s"
					}
				`, domainName),
				Check: resource.ComposeTestCheckFunc(
					isDomainPresent(tt, "data.scaleway_tem_domain.test"),
					resource.TestCheckResourceAttr("data.scaleway_tem_domain.test", "name", domainName),
					resource.TestCheckResourceAttrSet("data.scaleway_tem_domain.test", "reputation.0.status"),
					resource.TestCheckResourceAttrSet("data.scaleway_tem_domain.test", "reputation.0.score"),
					resource.TestCheckResourceAttrSet("data.scaleway_tem_domain.test", "reputation.0.scored_at"),
				),
			},
		},
	})
}
