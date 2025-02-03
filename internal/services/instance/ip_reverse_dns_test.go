package instance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccIPReverseDns_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testDNSZone := "tf-reverse-instance." + acctest.TestDomain
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_ip" "main" {}
					
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = %[1]q
						name     = ""
						type     = "A"
						data     = "${scaleway_instance_ip.main.address}"
						ttl      = 3600
						priority = 1
					}
				`, testDNSZone),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_ip" "main" {}
					
					resource "scaleway_domain_record" "tf_A" {
					  dns_zone = %[1]q
					  name     = ""
					  type     = "A"
					  data     = "${scaleway_instance_ip.main.address}"
					  ttl      = 3600
					  priority = 1
					}

					resource "scaleway_instance_ip_reverse_dns" "base" {
					  ip_id      = scaleway_instance_ip.main.id
					  reverse    = %[1]q
					  depends_on = [scaleway_domain_record.tf_A]
					}
				`, testDNSZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_ip_reverse_dns.base", "reverse", testDNSZone),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "ip" {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_ip.ip", "reverse", ""),
				),
			},
		},
	})
}
