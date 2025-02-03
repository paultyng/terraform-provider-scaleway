package domain_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceDomainRecord_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_domain_record main {
						dns_zone = "test-data-source.%s"
						name     = "www"
						type     = "A"
						data     = "1.2.3.4"
						ttl      = 3600
						priority = 10
					}

					data scaleway_domain_record test {
						dns_zone  = "${scaleway_domain_record.main.dns_zone}"
						record_id = "${scaleway_domain_record.main.id}"
					}

					data scaleway_domain_record test2 {
						dns_zone = "${scaleway_domain_record.main.dns_zone}"
						name     = "${scaleway_domain_record.main.name}"
						type     = "${scaleway_domain_record.main.type}"
						data     = "${scaleway_domain_record.main.data}"
					}
				`, acctest.TestDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "data.scaleway_domain_record.test"),
					testAccCheckDomainRecordExists(tt, "data.scaleway_domain_record.test2"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test", "dns_zone",
						"scaleway_domain_record.main", "dns_zone"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "id",
						"scaleway_domain_record.main", "id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "dns_zone",
						"scaleway_domain_record.main", "dns_zone"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "name",
						"scaleway_domain_record.main", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "type",
						"scaleway_domain_record.main", "type"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "data",
						"scaleway_domain_record.main", "data"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "ttl",
						"scaleway_domain_record.main", "ttl"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test2", "priority",
						"scaleway_domain_record.main", "priority"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_domain_record geo_ip {
						dns_zone = "test-data-source-geo-ip.%s"
						name     = "tf_geo_ip"
						type     = "A"
						data     = "1.2.3.4"
						geo_ip {
							matches {
								continents = ["EU"]
								countries  = ["FR"]
								data       = "1.2.3.4"
							}
							matches {
								continents = ["NA"]
								data       = "1.2.3.5"
							}
						}
					}

					data scaleway_domain_record test_geo_ip {
						dns_zone  = "${scaleway_domain_record.geo_ip.dns_zone}"
						record_id = "${scaleway_domain_record.geo_ip.id}"
					}
				`, acctest.TestDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "data.scaleway_domain_record.test_geo_ip"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "id",
						"scaleway_domain_record.geo_ip", "id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "dns_zone",
						"scaleway_domain_record.geo_ip", "dns_zone"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "name",
						"scaleway_domain_record.geo_ip", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "type",
						"scaleway_domain_record.geo_ip", "type"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "data",
						"scaleway_domain_record.geo_ip", "data"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "ttl",
						"scaleway_domain_record.geo_ip", "ttl"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "priority",
						"scaleway_domain_record.geo_ip", "priority"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_geo_ip", "geo_ip",
						"scaleway_domain_record.geo_ip", "geo_ip"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_domain_record http_service {
						dns_zone = "test-data-source-http-service.%s"
						name     = "tf_http_service"
						type     = "A"
						data     = "1.2.3.4"
						http_service {
							ips          = ["1.2.3.4", "4.3.2.1"]
							must_contain = "up"
							url          = "http://mywebsite.com/health"
							user_agent   = "scw_service_up"
							strategy     = "hashed"
						}
					}

					data scaleway_domain_record test_http_service {
						dns_zone  = "${scaleway_domain_record.http_service.dns_zone}"
						record_id = "${scaleway_domain_record.http_service.id}"
					}
				`, acctest.TestDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "data.scaleway_domain_record.test_http_service"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "id",
						"scaleway_domain_record.http_service", "id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "dns_zone",
						"scaleway_domain_record.http_service", "dns_zone"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "name",
						"scaleway_domain_record.http_service", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "type",
						"scaleway_domain_record.http_service", "type"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "data",
						"scaleway_domain_record.http_service", "data"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "ttl",
						"scaleway_domain_record.http_service", "ttl"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "priority",
						"scaleway_domain_record.http_service", "priority"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_http_service", "http_service",
						"scaleway_domain_record.http_service", "http_service"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_domain_record view {
						dns_zone = "test-data-source-view.%s"
						name     = "tf_view"
						type     = "A"
						data     = "1.2.3.4"
						view {
							subnet = "100.0.0.0/16"
							data   = "1.2.3.4"
						}
						view {
							subnet = "100.1.0.0/16"
							data   = "4.3.2.1"
						}
					}

					data scaleway_domain_record test_view {
						dns_zone  = "${scaleway_domain_record.view.dns_zone}"
						record_id = "${scaleway_domain_record.view.id}"
					}
				`, acctest.TestDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "data.scaleway_domain_record.test_view"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "id",
						"scaleway_domain_record.view", "id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "dns_zone",
						"scaleway_domain_record.view", "dns_zone"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "name",
						"scaleway_domain_record.view", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "type",
						"scaleway_domain_record.view", "type"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "data",
						"scaleway_domain_record.view", "data"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "ttl",
						"scaleway_domain_record.view", "ttl"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "priority",
						"scaleway_domain_record.view", "priority"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_view", "view",
						"scaleway_domain_record.view", "view"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource scaleway_domain_record weighted {
						dns_zone = "test-data-source-weighted.%s"
						name     = "tf_weighted"
						type     = "A"
						data     = "1.2.3.4"
						weighted {
							ip     = "1.2.3.4"
							weight = 1
						}
						weighted {
							ip     = "4.3.2.1"
							weight = 2
						}
					}

					data scaleway_domain_record test_weighted {
						dns_zone  = "${scaleway_domain_record.weighted.dns_zone}"
						record_id = "${scaleway_domain_record.weighted.id}"
					}
				`, acctest.TestDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "data.scaleway_domain_record.test_weighted"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "id",
						"scaleway_domain_record.weighted", "id"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "dns_zone",
						"scaleway_domain_record.weighted", "dns_zone"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "name",
						"scaleway_domain_record.weighted", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "type",
						"scaleway_domain_record.weighted", "type"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "data",
						"scaleway_domain_record.weighted", "data"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "ttl",
						"scaleway_domain_record.weighted", "ttl"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "priority",
						"scaleway_domain_record.weighted", "priority"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_domain_record.test_weighted", "weighted",
						"scaleway_domain_record.weighted", "weighted"),
				),
			},
		},
	})
}
