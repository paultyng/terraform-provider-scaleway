package lb_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
)

func TestAccBackend_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isBackendDestroyed(tt),
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
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "forward_port_algorithm", "roundrobin"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "sticky_sessions", "none"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "proxy_protocol", "none"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_server", "5m0s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_connect", "5s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_tunnel", "15m0s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "on_marked_down_action", "none"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_timeout", "30s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "80"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_delay", "1m0s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_max_retries", "2"),
					resource.TestCheckResourceAttrPair("scaleway_lb_backend.bkd01", "server_ips.0", "scaleway_instance_ip.ip01", "address"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ssl_bridging", "false"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ignore_ssl_server_verify", "false"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "redispatch_attempt_count", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "max_retries", "3"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_transient_delay", "500ms"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_send_proxy", "false"),
					resource.TestCheckResourceAttrSet("scaleway_lb_backend.bkd01", "send_proxy_v2"), // Deprecated attribute
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
						max_connections = 42
						timeout_queue = "4s"
						redispatch_attempt_count = 1
						max_retries = 6
						health_check_transient_delay = "0.2s"
						health_check_send_proxy = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttrPair("scaleway_lb_backend.bkd01", "server_ips.0", "scaleway_instance_ip.ip01", "address"),
					resource.TestCheckResourceAttrPair("scaleway_lb_backend.bkd01", "server_ips.1", "scaleway_instance_ip.ip02", "address"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_server", "1s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_connect", "2.5s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_tunnel", "3s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_delay", "10s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_timeout", "15s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "81"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_max_retries", "3"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "on_marked_down_action", "shutdown_sessions"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ssl_bridging", "true"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "ignore_ssl_server_verify", "true"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "max_connections", "42"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "timeout_queue", "4s"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "redispatch_attempt_count", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "max_retries", "6"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_transient_delay", "200ms"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_send_proxy", "true"),
				),
			},
		},
	})
}

func TestAccBackend_HealthCheck(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isBackendDestroyed(tt),
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

func TestAccBackend_WithFailoverHost(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	resourceName := "scaleway_object_bucket_website_configuration.test"

	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isBackendDestroyed(tt),
			objectchecks.IsObjectDestroyed(tt),
			objectchecks.IsBucketDestroyed(tt),
			objectchecks.IsWebsiteConfigurationDestroyed(tt),
		),
		Steps: []resource.TestStep{
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
				`, rName),
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
					}
				`, rName),
				Check: isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
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
					objectchecks.IsWebsiteConfigurationPresent(tt, resourceName),
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr(resourceName, "website_endpoint", rName+".s3-website.fr-par.scw.cloud"),
					resource.TestCheckResourceAttrSet("scaleway_lb_backend.bkd01", "failover_host"),
				),
				ExpectNonEmptyPlan: !*acctest.UpdateCassettes,
			},
		},
	})
}

func TestAccBackend_HealthCheck_Port(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isBackendDestroyed(tt),
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
					forward_port = "3333"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "forward_port", "3333"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "3333"),
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
					forward_port = "4444"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "forward_port", "4444"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "4444"),
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
					forward_port = "4444"
					health_check_port = "4444"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "forward_port", "4444"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "4444"),
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
					forward_port = "5555"
					health_check_port = "4444"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					isBackendPresent(tt, "scaleway_lb_backend.bkd01"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "forward_port", "5555"),
					resource.TestCheckResourceAttr("scaleway_lb_backend.bkd01", "health_check_port", "4444"),
				),
			},
		},
	})
}

func isBackendPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
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

func isBackendDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_backend" {
				continue
			}

			lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
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
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
