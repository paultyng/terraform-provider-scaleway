package domain_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	domainSDK "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/domain"
)

func TestAccDomainRecord_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-basic." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_Basic: test dns zone: %s", testDNSZone)

	name := "tf"
	recordType := "A"
	data := "127.0.0.1"
	dataUpdated := "127.0.0.2"
	ttl := 3600
	ttlUpdated := 43200
	priority := 0
	priorityUpdated := 10

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, recordType, data, ttl, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", ""),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "fqdn", testDNSZone),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, name, recordType, data, ttl, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "fqdn", name+"."+testDNSZone),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, name, recordType, dataUpdated, ttl, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", dataUpdated),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "fqdn", name+"."+testDNSZone),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, name, recordType, dataUpdated, ttlUpdated, priorityUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", dataUpdated),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", strconv.Itoa(ttlUpdated)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", strconv.Itoa(priorityUpdated)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "fqdn", name+"."+testDNSZone),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_domain_record" "tf_A" {
						dns_zone = %[1]q
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
				}

				resource "scaleway_domain_record" "tf_MX" {
					dns_zone = %[1]q
					name     = "record_mx"
					type     = "MX"
					data     = "ASPMX.L.GOOGLE.COM."
					ttl      = 600
					priority = 1
				}
			`, testDNSZone, name, recordType, dataUpdated, ttlUpdated, priorityUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_MX"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "name", "record_mx"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "type", "MX"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "data", "ASPMX.L.GOOGLE.COM."),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "ttl", "600"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "priority", "1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_MX", "fqdn", "record_mx."+testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "fqdn", name+"."+testDNSZone),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_MX", "id"),
				),
			},
		},
	})
}

func TestAccDomainRecord_Basic2(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-basic2." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_Basic: test dns zone: %s", testDNSZone)

	recordType := "A"
	data := "127.0.0.1"
	ttl := 3600
	priority := 0

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = %[1]q
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}

					resource "scaleway_domain_record" "aws_mx" {
					  dns_zone = %[1]q
					  name     = ""
					  type     = "MX"
					  data     = "10 feedback-smtp.eu-west-1.amazonses.com."
					  ttl      = 300
					}
					
					resource "scaleway_domain_record" "mx" {
					  dns_zone = %[1]q
					  name     = ""
					  type     = "MX"
					  data     = "0 mail.scaleway.com."
					  ttl      = 300
					}
					
					resource "scaleway_domain_record" "txt_dmarc" {
					  dns_zone = %[1]q
					  name     = "_dmarc"
					  type     = "TXT"
					  data     = "v=DMARC1; p=quarantine; adkim=s"
					  ttl      = 3600
					}
				`, testDNSZone, recordType, data, ttl, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.mx"),
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.txt_dmarc"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", ""),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", strconv.Itoa(priority)),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.aws_mx"),
					resource.TestCheckResourceAttr("scaleway_domain_record.aws_mx", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.aws_mx", "name", ""),
					resource.TestCheckResourceAttr("scaleway_domain_record.aws_mx", "type", "MX"),
					resource.TestCheckResourceAttr("scaleway_domain_record.aws_mx", "data", "10 feedback-smtp.eu-west-1.amazonses.com."),
					resource.TestCheckResourceAttr("scaleway_domain_record.aws_mx", "ttl", "300"),
					resource.TestCheckResourceAttr("scaleway_domain_record.aws_mx", "priority", "10"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.aws_mx", "id"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.mx", "id"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.txt_dmarc", "id"),
				),
			},
		},
	})
}

func TestAccDomainRecord_Arobase(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-arobase." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_Arobase: test dns zone: %s", testDNSZone)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "main" {
						dns_zone = %[1]q
						name     = "@"
						type     = "TXT"
						data     = "this-is-a-test"
					}
				`, testDNSZone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.main"),
					resource.TestCheckResourceAttr("scaleway_domain_record.main", "name", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "main" {
						dns_zone = %[1]q
						name     = ""
						type     = "TXT"
						data     = "this-is-a-test"
					}
				`, testDNSZone),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDomainRecord_GeoIP(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-geoip." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_GeoIP: test dns zone: %s", testDNSZone)

	name := "tf_geo_ip"
	recordType := "A"
	data := "127.0.0.2"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A_geo_ip" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
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
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_geo_ip"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.continents.0", "EU"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.countries.0", "FR"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.continents.0", "NA"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.data", "1.2.3.5"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_geo_ip", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_domain_record" "tf_A_geo_ip" {
					dns_zone = "%s"
					name     = "%s"
					type     = "%s"
					data     = "%s"
					geo_ip {
						matches {
							continents = ["EU","AS"]
							countries  = ["FR","AE"]
							data       = "1.2.3.4"
						}
						matches {
							countries  = ["CI"]
							data       = "1.2.3.5"
						}
					}
				}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_geo_ip"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.continents.0", "EU"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.continents.1", "AS"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.countries.0", "FR"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.countries.1", "AE"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.countries.0", "CI"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.data", "1.2.3.5"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_geo_ip", "id"),
				),
			},
		},
	})
}

func TestAccDomainRecord_HTTPService(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-httpservice." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_HTTPService: test dns zone: %s", testDNSZone)

	name := "tf_http_service"
	recordType := "A"
	data := "127.0.0.3"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A_http_service" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						http_service {
							ips          = ["1.2.3.4", "4.3.2.1"]
							must_contain = "up"
							url          = "http://mywebsite.com/health"
							user_agent   = "scw_service_up"
							strategy     = "hashed"
						}
					}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_http_service"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.ips.0", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.ips.1", "4.3.2.1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.must_contain", "up"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.url", "http://mywebsite.com/health"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.user_agent", "scw_service_up"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.strategy", "hashed"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_http_service", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_domain_record" "tf_A_http_service" {
					dns_zone = "%s"
					name     = "%s"
					type     = "%s"
					data     = "%s"
					http_service {
						ips          = ["5.6.7.8"]
						must_contain = "online"
						url          = "http://mywebsite.com/healthcheck"
						user_agent   = "scw_service_online"
						strategy     = "random"
					}
				}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_http_service"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.ips.0", "5.6.7.8"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.must_contain", "online"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.url", "http://mywebsite.com/healthcheck"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.user_agent", "scw_service_online"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.strategy", "random"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_http_service", "id"),
				),
			},
		},
	})
}

func TestAccDomainRecord_View(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-view." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_View: test dns zone: %s", testDNSZone)

	name := "tf_view"
	recordType := "A"
	data := "127.0.0.4"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A_view" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						view {
							subnet = "100.0.0.0/16"
							data   = "1.2.3.4"
						}
						view {
							subnet = "100.1.0.0/16"
							data   = "4.3.2.1"
						}
					}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_view"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.subnet", "100.0.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.subnet", "100.1.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.data", "4.3.2.1"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_view", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_domain_record" "tf_A_view" {
					dns_zone = "%s"
					name     = "%s"
					type     = "%s"
					data     = "%s"
					view {
						subnet = "100.0.0.0/16"
						data   = "1.2.3.4"
					}
					view {
						subnet = "90.1.0.0/32"
						data   = "4.3.2.2"
					}
					view {
						subnet = "1.1.1.1/16"
						data   = "2.2.2.2"
					}
				}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_view"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.subnet", "100.0.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.subnet", "90.1.0.0/32"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.data", "4.3.2.2"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.2.subnet", "1.1.1.1/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.2.data", "2.2.2.2"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_view", "id"),
				),
			},
		},
	})
}

func TestAccDomainRecord_Weighted(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-weighted." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_Weighted: test dns zone: %s", testDNSZone)

	name := "tf_weighted"
	recordType := "A"
	data := "127.0.0.5"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_A_weighted" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						weighted {
							ip     = "1.2.3.4"
							weight = 1
						}
						weighted {
							ip     = "4.3.2.1"
							weight = 2
						}
					}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_weighted"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.weight", "1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.ip", "4.3.2.1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.weight", "2"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_weighted", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_domain_record" "tf_A_weighted" {
					dns_zone = "%s"
					name     = "%s"
					type     = "%s"
					data     = "%s"
					weighted {
						ip     = "1.2.3.4"
						weight = 2
					}
					weighted {
						ip     = "4.3.2.1"
						weight = 1
					}
					weighted {
						ip     = "5.6.7.8"
						weight = 999
					}
				}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_A_weighted"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "priority", strconv.Itoa(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.weight", "2"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.ip", "4.3.2.1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.weight", "1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.2.ip", "5.6.7.8"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.2.weight", "999"),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_A_weighted", "id"),
				),
			},
		},
	})
}

func TestAccDomainRecord_SRVZone(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-srv." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_SRVZone: test dns zone: %s", testDNSZone)

	name := "_proxy-preproduction._tcp"
	recordType := "SRV"
	data := "100 1 3128 bigbox.example.com."
	ttl := 3600     // default value
	priority := 100 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "proxy_srv" {
						dns_zone        = "%[1]s"
						name            = "%[2]s"
						type            = "%[3]s"
						data            = "%[4]s"
						keep_empty_zone = false
					}
				`, testDNSZone, name, recordType, data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.proxy_srv"),
					resource.TestCheckResourceAttr("scaleway_domain_record.proxy_srv", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.proxy_srv", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.proxy_srv", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.proxy_srv", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.proxy_srv", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.proxy_srv", "priority", strconv.Itoa(priority)),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.proxy_srv", "id"),
				),
			},
		},
	})
}

func testAccCheckDomainRecordExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		domainAPI := domain.NewDomainAPI(tt.Meta)
		listDNSZones, err := domainAPI.ListDNSZoneRecords(&domainSDK.ListDNSZoneRecordsRequest{
			DNSZone: rs.Primary.Attributes["dns_zone"],
		})
		if err != nil {
			return err
		}

		for _, record := range listDNSZones.Records {
			if record.ID == rs.Primary.ID {
				// record found
				return nil
			}
		}

		return fmt.Errorf("record (%s) not found in: %s", rs.Primary.ID, rs.Primary.Attributes["dns_zone"])
	}
}

func TestAccDomainRecord_CNAME(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-basic-cname." + acctest.TestDomain
	logging.L.Debugf("TestAccScalewayDomainRecord_Basic: test dns zone: %s", testDNSZone)

	name := "tf"
	recordType := "CNAME"
	data := "xxx.scw.cloud"
	dataUpdated := "yyy.scw.cloud"
	ttl := 3600
	ttlUpdated := 43200
	priority := 0
	priorityUpdated := 10

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_CNAME" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, name, recordType, data, ttl, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_CNAME"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "priority", strconv.Itoa(priority)),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_CNAME", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_CNAME" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, name, recordType, dataUpdated, ttl, priority),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_CNAME"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "data", dataUpdated),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "ttl", strconv.Itoa(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "priority", strconv.Itoa(priority)),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_CNAME", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_record" "tf_CNAME" {
						dns_zone = "%s"
						name     = "%s"
						type     = "%s"
						data     = "%s"
						ttl      = %d
						priority = %d
					}
				`, testDNSZone, name, recordType, dataUpdated, ttlUpdated, priorityUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainRecordExists(tt, "scaleway_domain_record.tf_CNAME"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "data", dataUpdated),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "ttl", strconv.Itoa(ttlUpdated)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_CNAME", "priority", strconv.Itoa(priorityUpdated)),
					acctest.CheckResourceAttrUUID("scaleway_domain_record.tf_CNAME", "id"),
				),
			},
		},
	})
}

func testAccCheckDomainRecordDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_domain_record" {
				continue
			}

			// check if the zone still exists
			domainAPI := domain.NewDomainAPI(tt.Meta)
			listDNSZones, err := domainAPI.ListDNSZoneRecords(&domainSDK.ListDNSZoneRecordsRequest{
				DNSZone: rs.Primary.Attributes["dns_zone"],
			})
			if httperrors.Is403(err) { // forbidden: subdomain not found
				return nil
			}

			if err != nil {
				return fmt.Errorf("failed to check if domain zone exists: %w", err)
			}

			if listDNSZones.TotalCount > 0 {
				return fmt.Errorf("zone %s still exist", rs.Primary.Attributes["dns_zone"])
			}
			return nil
		}

		return nil
	}
}
