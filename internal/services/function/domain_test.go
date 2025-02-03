package function_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	containerchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/container/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/function"
)

func TestAccFunctionDomain_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "function-basic." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayContainerDomain_Basic: test dns zone: %s", testDNSZone)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource scaleway_function_namespace main {}

				resource scaleway_function main {
						namespace_id = scaleway_function_namespace.main.id
						runtime = "go122"
						privacy = "private"
						handler = "Handle"
						zip_file = "testfixture/gofunction.zip"
						deploy = true
				}
				`,
				Check: containerchecks.TestConfigContainerNamespace(tt, "scaleway_function_namespace.main"),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						namespace_id = scaleway_function_namespace.main.id
						runtime = "go122"
						privacy = "private"
						handler = "Handle"
						zip_file = "testfixture/gofunction.zip"
						deploy = true
					}

					resource scaleway_domain_record "function" {
				  		dns_zone = "%s"
				  		name     = "container"
				  		type     = "CNAME"
				  		data     = "${scaleway_function.main.domain_name}."
				  		ttl      = 60
					}

					resource "scaleway_function_domain" "main" {
					  function_id = scaleway_function.main.id
					  hostname    = "${scaleway_domain_record.function.name}.${scaleway_domain_record.function.dns_zone}"
					}
				`, testDNSZone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionDomainExists(tt, "scaleway_function_domain.main"),
				),
			},
		},
	})
}

func testAccCheckFunctionDomainExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetDomain(&functionSDK.GetDomainRequest{
			DomainID: id,
			Region:   region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFunctionDomainDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_function_domain" {
				continue
			}

			api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteDomain(&functionSDK.DeleteDomainRequest{
				DomainID: id,
				Region:   region,
			})

			if err == nil {
				return fmt.Errorf("function domain (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
