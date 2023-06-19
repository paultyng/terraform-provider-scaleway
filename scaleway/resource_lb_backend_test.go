package scaleway

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func TestAccScalewayLbBackend_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayLbBackendDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_instance_ip ip01 {}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
						server_ips = [ scaleway_instance_ip.ip01.address ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayLbBackendExists(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "forward_port_algorithm", "roundrobin"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "sticky_sessions", "none"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "proxy_protocol", "none"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_server", ""),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "on_marked_down_action", "none"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_timeout", "30s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "80"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_delay", "1m0s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_max_retries", "2"),
					resource.TestCheckResourceAttrPair("scaleway_lb_backend.bkd01", "server_ips.0", "scaleway_instance_ip.ip01", "address"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ssl_bridging", "false"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ignore_ssl_server_verify", "false"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_instance_ip ip01 {}
					resource scaleway_instance_ip ip02 {}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "tcp"
						forward_port = 80
						forward_port_algorithm = "leastconn"
						sticky_sessions = "cookie"
						sticky_sessions_cookie_name = "session-id"
						server_ips = [scaleway_instance_ip.ip01.address , scaleway_instance_ip.ip02.address ]
						proxy_protocol = "none"
						timeout_server = "1s"
						timeout_connect = "2.5s"
						timeout_tunnel = "3s"
						health_check_timeout = "15s"
						health_check_delay = "10s"
						health_check_port = 81
						health_check_max_retries = 3
						on_marked_down_action = "shutdown_sessions"
						ssl_bridging = "true"
						ignore_ssl_server_verify = "true"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayLbBackendExists(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttrPair("scaleway_lb_backend.bkd01", "server_ips.0", "scaleway_instance_ip.ip01", "address"),
					resource.TestCheckResourceAttrPair("scaleway_lb_backend.bkd01", "server_ips.1", "scaleway_instance_ip.ip02", "address"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_delay", "10s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_timeout", "15s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "81"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_max_retries", "3"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "on_marked_down_action", "shutdown_sessions"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ssl_bridging", "true"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ignore_ssl_server_verify", "true"),
				),
			},
		},
	})
}

func TestAccScalewayLbBackend_HealthCheck_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayLbBackendDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_lb_ip" "ip01" {}

				resource "scaleway_lb" "lb01" {
				  ip_id = scaleway_lb_ip.ip01.id
				  name  = "test-lb"
				  type  = "lb-s"
				}

				resource "scaleway_lb_backend" "bkd01" {
				  lb_id            = scaleway_lb.lb01.id
				  name             = "bkd01"
				  forward_protocol = "tcp"
				  forward_port     = 80
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_tcp.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.#", "0"),
				),
			},
			{
				Config: `
						resource scaleway_lb_ip ip01 {}

						resource scaleway_lb lb01 {
							ip_id = scaleway_lb_ip.ip01.id
							name = "test-lb"
							type = "lb-s"
						}

						resource scaleway_lb_backend bkd01 {
							lb_id = scaleway_lb.lb01.id
							name = "bkd01"
							forward_protocol = "tcp"
							forward_port = 80

							health_check {
								protocol = "http"
								uri = "http://test.com/health"
								method = "POST"
								code = 404
								host_header = "test.com"
							}
						}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_tcp.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.uri", "http://test.com/health"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.method", "POST"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.code", "404"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.host_header", "test.com"),
				),
			},
			{
				Config: `
					resource "scaleway_lb_ip" "ip01" {}

					resource "scaleway_lb" "lb01" {
					  ip_id = scaleway_lb_ip.ip01.id
					  name  = "test-lb"
					  type  = "lb-s"
					}

					resource "scaleway_lb_backend" "bkd01" {
					  lb_id            = scaleway_lb.lb01.id
					  name             = "bkd01"
					  forward_protocol = "tcp"
					  forward_port     = 80

					  health_check {
						protocol    = "https"
						uri         = "http://test.com/health"
						method      = "POST"
						code        = 404
						host_header = "test.com"
						sni         = "sni.test.com"
					  }
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_tcp.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.uri", "http://test.com/health"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.method", "POST"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.protocol", "https"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.code", "404"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.host_header", "test.com"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.sni", "sni.test.com"),
				),
			},
			{
				Config: `
					resource "scaleway_lb_ip" "ip01" {}
					
					resource "scaleway_lb" "lb01" {
					  ip_id = scaleway_lb_ip.ip01.id
					  name  = "test-lb"
					  type  = "lb-s"
					}
					
					resource "scaleway_lb_backend" "bkd01" {
					  lb_id            = scaleway_lb.lb01.id
					  name             = "bkd01"
					  forward_protocol = "tcp"
					  forward_port     = 80
					
					  health_check {
						protocol      = "mysql"
						database_user = "devtools"
					  }
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.protocol", "mysql"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check.0.code", "0"),
				),
			},
		},
	})
}

func TestAccScalewayLbBackend_HealthCheck(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayLbBackendDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "tcp"
						forward_port = 80
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_tcp.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.#", "0"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "tcp"
						forward_port = 80

						health_check_http {
							uri = "http://test.com/health"
							method = "POST"
							code = 404
							host_header = "test.com"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_tcp.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.0.uri", "http://test.com/health"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.0.method", "POST"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.0.code", "404"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.0.host_header", "test.com"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "tcp"
						forward_port = 80

						health_check_https {
							uri = "http://test.com/health"
							method = "POST"
							code = 404
							host_header = "test.com"
							sni = "sni.test.com"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_tcp.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_http.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.0.uri", "http://test.com/health"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.0.method", "POST"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.0.code", "404"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.0.host_header", "test.com"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_https.0.sni", "sni.test.com"),
				),
			},
		},
	})
}

func TestAccScalewayLbBackend_WithFailoverHost(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "scaleway_object_bucket_website_configuration.test"

	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckScalewayLbBackendDestroy(tt),
			testAccCheckScalewayObjectDestroy(tt),
			testAccCheckScalewayObjectBucketDestroy(tt),
			testAccCheckBucketWebsiteConfigurationDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_instance_ip ip01 {}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
						server_ips = [ scaleway_instance_ip.ip01.address ]
					}
				`,
				Check: testAccCheckScalewayLbBackendExists(tt, "scaleway_lb_backend.bkd01"),
			},
			{
				Config: fmt.Sprintf(`
			  		resource "scaleway_object_bucket" "test" {
						name = %[1]q
						acl  = "public-read"
						tags = {
							TestName = "TestAccSCW_WebsiteConfig_basic"
						}
					}

					resource scaleway_object "some_file" {
						bucket = scaleway_object_bucket.test.name
						key = "index.html"
						file = "testfixture/index.html"
						visibility = "public-read"
					}
				
				  	resource "scaleway_object_bucket_website_configuration" "test" {
						bucket = scaleway_object_bucket.test.name
						index_document {
							suffix = "index.html"
						}
						error_document {
							key = "error.html"
						}
				  	}

					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_instance_ip ip01 {}

					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						name = "bkd01"
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
						server_ips = [ scaleway_instance_ip.ip01.address ]
						failover_host = scaleway_object_bucket_website_configuration.test.website_endpoint
					}
				`, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketWebsiteConfigurationExists(tt, resourceName),
					testAccCheckScalewayLbBackendExists(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", rName+".s3-website.fr-par.scw.cloud"),
					resource.TestCheckResourceAttrSet("scaleway_lb_backend.bkd01", "failover_host"),
				),
				ExpectNonEmptyPlan: !*UpdateCassettes,
			},
		},
	})
}

func testAccCheckScalewayLbBackendExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, zone, ID, err := lbAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = lbAPI.GetBackend(&lbSDK.ZonedAPIGetBackendRequest{
			BackendID: ID,
			Zone:      zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayLbBackendDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_backend" {
				continue
			}

			lbAPI, zone, ID, err := lbAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = lbAPI.GetBackend(&lbSDK.ZonedAPIGetBackendRequest{
				Zone:      zone,
				BackendID: ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("LB Backend (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}
