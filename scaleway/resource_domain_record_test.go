package scaleway

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	// prevent using production domain for testing
	reservedDomains = []*regexp.Regexp{
		regexp.MustCompile(`.*iliad.*`),
		regexp.MustCompile(`.*\.free\..*`),
		regexp.MustCompile(`.*\.online\..*`),
		regexp.MustCompile(`.*scaleway.*`),
		regexp.MustCompile(`.*dedibox.*`),
	}
	testDomain = ""
)

func init() {
	// if domain is set, use it
	testDomainPtr := flag.String("test-domain", os.Getenv("TF_TEST_DOMAIN"), "Test domain")
	if testDomainPtr != nil && *testDomainPtr != "" {
		testDomain = *testDomainPtr
		l.Infof("start domain record test with set domain: %s", testDomain)

		return
	}

	// retrieve the first domain is found
	scwClient, err := sharedClientForZone(scw.Zone(""))
	if err != nil {
		l.Errorf("error when retrieve scwClient: %s", err)
		return
	}

	domainRegistrarAPI := domain.NewRegistrarAPI(scwClient)
	listDomains, err := domainRegistrarAPI.ListDomains(&domain.RegistrarAPIListDomainsRequest{
		Status:     domain.DomainStatusActive,
		OrderBy:    domain.ListDomainsRequestOrderByDomainAsc,
		IsExternal: scw.BoolPtr(true),
	})
	if err != nil {
		l.Errorf("error listing domains: %s", err)
		return
	}

	for _, domain := range listDomains.Domains {
		// skipping reserved domains
		isReserved := false
		for _, reservedDomain := range reservedDomains {
			if reservedDomain.MatchString(domain.Domain) {
				isReserved = true
				break
			}
		}

		if !isReserved {
			testDomain = domain.Domain
			break
		}
	}

	if testDomain == "" {
		l.Errorf("your account doesn't have any domain to test with")
		return
	}

	l.Infof("start domain record test with domain: %s", testDomain)
}

func TestAccScalewayDomainRecord_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := fmt.Sprintf("test-basic.%s", testDomain)
	l.Debugf("TestAccScalewayDomainRecord_Basic: test dns zone: %s", testDNSZone)

	name := "tf"
	recordType := "A"
	data := "127.0.0.1"
	dataUpdated := "127.0.0.2"
	ttl := 3600
	ttlUpdated := 43200
	priority := 0
	priorityUpdated := 10

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayDomainRecordDestroy(tt),
		Steps: []resource.TestStep{
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", fmt.Sprint(priority)),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", dataUpdated),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", fmt.Sprint(priority)),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "data", dataUpdated),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "ttl", fmt.Sprint(ttlUpdated)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A", "priority", fmt.Sprint(priorityUpdated)),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A", "id"),
				),
			},
		},
	})
}

func TestAccScalewayDomainRecord_GeoIP(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := fmt.Sprintf("test-geoip.%s", testDomain)
	l.Debugf("TestAccScalewayDomainRecord_GeoIP: test dns zone: %s", testDNSZone)

	name := "tf_geo_ip"
	recordType := "A"
	data := "127.0.0.2"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayDomainRecordDestroy(tt),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_geo_ip"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.continents.0", "EU"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.countries.0", "FR"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.continents.0", "NA"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.data", "1.2.3.5"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_geo_ip", "id"),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_geo_ip"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.continents.0", "EU"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.continents.1", "AS"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.countries.0", "FR"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.countries.1", "AE"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.countries.0", "CI"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_geo_ip", "geo_ip.0.matches.1.data", "1.2.3.5"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_geo_ip", "id"),
				),
			},
		},
	})
}

func TestAccScalewayDomainRecord_HTTPService(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := fmt.Sprintf("test-httpservice.%s", testDomain)
	l.Debugf("TestAccScalewayDomainRecord_HTTPService: test dns zone: %s", testDNSZone)

	name := "tf_http_service"
	recordType := "A"
	data := "127.0.0.3"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayDomainRecordDestroy(tt),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_http_service"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.ips.0", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.ips.1", "4.3.2.1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.must_contain", "up"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.url", "http://mywebsite.com/health"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.user_agent", "scw_service_up"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.strategy", "hashed"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_http_service", "id"),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_http_service"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.ips.0", "5.6.7.8"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.must_contain", "online"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.url", "http://mywebsite.com/healthcheck"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.user_agent", "scw_service_online"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_http_service", "http_service.0.strategy", "random"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_http_service", "id"),
				),
			},
		},
	})
}

func TestAccScalewayDomainRecord_View(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := fmt.Sprintf("test-view.%s", testDomain)
	l.Debugf("TestAccScalewayDomainRecord_View: test dns zone: %s", testDNSZone)

	name := "tf_view"
	recordType := "A"
	data := "127.0.0.4"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayDomainRecordDestroy(tt),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_view"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.subnet", "100.0.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.subnet", "100.1.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.data", "4.3.2.1"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_view", "id"),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_view"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.subnet", "100.0.0.0/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.0.data", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.subnet", "90.1.0.0/32"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.1.data", "4.3.2.2"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.2.subnet", "1.1.1.1/16"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_view", "view.2.data", "2.2.2.2"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_view", "id"),
				),
			},
		},
	})
}

func TestAccScalewayDomainRecord_Weighted(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := fmt.Sprintf("test-weighted.%s", testDomain)
	l.Debugf("TestAccScalewayDomainRecord_Weighted: test dns zone: %s", testDNSZone)

	name := "tf_weighted"
	recordType := "A"
	data := "127.0.0.5"
	ttl := 3600   // default value
	priority := 0 // default value

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayDomainRecordDestroy(tt),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_weighted"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.weight", "1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.ip", "4.3.2.1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.weight", "2"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_weighted", "id"),
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
					testAccCheckScalewayDomainRecordExists(tt, "scaleway_domain_record.tf_A_weighted"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "dns_zone", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "name", name),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "type", recordType),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "data", data),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "ttl", fmt.Sprint(ttl)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "priority", fmt.Sprint(priority)),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.0.weight", "2"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.ip", "4.3.2.1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.1.weight", "1"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.2.ip", "5.6.7.8"),
					resource.TestCheckResourceAttr("scaleway_domain_record.tf_A_weighted", "weighted.2.weight", "999"),
					testCheckResourceAttrUUID("scaleway_domain_record.tf_A_weighted", "id"),
				),
			},
		},
	})
}

func testAccCheckScalewayDomainRecordExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		domainAPI := newDomainAPI(tt.Meta)
		listDNSZones, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
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

func testAccCheckScalewayDomainRecordDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_domain_record" {
				continue
			}

			// check if the zone still exists
			domainAPI := newDomainAPI(tt.Meta)
			listDNSZones, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
				DNSZone: rs.Primary.Attributes["dns_zone"],
			})
			if is403Error(err) { // forbidden: subdomain not found
				return nil
			}

			if err != nil {
				return err
			}

			if listDNSZones.TotalCount > 0 {
				return fmt.Errorf("zone %s still exist", rs.Primary.Attributes["dns_zone"])
			}
			return nil
		}

		return nil
	}
}
