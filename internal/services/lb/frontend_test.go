package lb_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestAccFrontend_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isFrontendDestroyed(tt),
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
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						inbound_port = 80
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isFrontendPresent(tt, "scaleway_lb_frontend.frt01"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "inbound_port", "80"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "timeout_client", ""),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "enable_http3", "false"),
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
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name = "tf-test"
						inbound_port = 443
						timeout_client = "30s"
						enable_http3 = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isFrontendPresent(tt, "scaleway_lb_frontend.frt01"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "name", "tf-test"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "inbound_port", "443"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "timeout_client", "30s"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "enable_http3", "true"),
				),
			},
		},
	})
}

// TODO: Refactor this test to enable testing of several custom domain names
//
// Let's encrypt currently has a limit of 50 certificates per week and a limit of 5 certificates per week per set of domains (including alternative names).
// So we need to change the list of alternative domain names to be able to test more than one domain name.
// One possible way to circumvent this limitation is to generate for a random set of alternative domain names that are all subdomains of the main test domain.
// For instance: *.test.scaleway-terraform.com which is a wildcard domain name.
// And we generate certificate for foo.test.scaleway-terraform.com, bar.test.scaleway-terraform.com, baz.test.scaleway-terraform.com, etc.
// Even changing one alternative domain name is enough to count as a new certificate (which is rate limited by the 50 certificates per week limit and not the 5 duplicate certificates per week limit).
// The only limitation is that all subdomains must resolve to the same IP address.
func TestAccFrontend_Certificate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isFrontendDestroyed(tt),
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
						forward_protocol = "http"
						forward_port = 443
						proxy_protocol = "none"
					}

					resource scaleway_lb_certificate cert01 {
						lb_id = scaleway_lb.lb01.id
						name = "test-cert-front-end"
					  	letsencrypt {
							common_name = "${replace(scaleway_lb_ip.ip01.ip_address,".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud"
					  	}
					}

					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						inbound_port = 443
						certificate_ids = [scaleway_lb_certificate.cert01.id]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isFrontendPresent(tt, "scaleway_lb_frontend.frt01"),
					isFrontendCertificatePresent(tt, "scaleway_lb_frontend.frt01", "scaleway_lb_certificate.cert01"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "certificate_ids.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_lb_frontend.frt01", "certificate_id"),
				),
			},
		},
	})
}

func isFrontendCertificatePresent(tt *acctest.TestTools, f, c string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[f]
		if !ok {
			return fmt.Errorf("resource not found: %s", f)
		}

		cs, ok := s.RootModule().Resources[c]
		if !ok {
			return fmt.Errorf("resource not found: %s", c)
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		frEnd, err := lbAPI.GetFrontend(&lbSDK.ZonedAPIGetFrontendRequest{
			FrontendID: ID,
			Zone:       zone,
		})
		if err != nil {
			return err
		}

		for _, id := range frEnd.CertificateIDs {
			if locality.ExpandID(cs.Primary.ID) == id {
				return nil
			}
		}

		return fmt.Errorf("certificate not found: %s", c)
	}
}

func isFrontendPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = lbAPI.GetFrontend(&lbSDK.ZonedAPIGetFrontendRequest{
			FrontendID: ID,
			Zone:       zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isFrontendDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_frontend" {
				continue
			}

			lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = lbAPI.GetFrontend(&lbSDK.ZonedAPIGetFrontendRequest{
				Zone:       zone,
				FrontendID: ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("LB Frontend (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func TestAccFrontend_ACLBasic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isFrontendDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-acl"
						type = "lb-s"
					}
					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name = "tf-test"
						inbound_port = 80
						timeout_client = "30s"
						acl {
							name = "test-acl"
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["192.168.0.1", "192.168.0.2", "192.168.10.0/24"]
								http_filter = "acl_http_filter_none"
								http_filter_value = []
								invert = "true"
							}
						}
						acl {
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["0.0.0.0/0"]
								http_filter = "path_begin"
								http_filter_value = ["criteria1","criteria2"]
								invert = "true"
							}
						}
						acl {
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["0.0.0.0/0"]
								http_filter = "path_begin"
								http_filter_value = ["criteria1","criteria2"]
							}
						}
						acl {
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["0.0.0.0/0"]
								http_filter = "acl_http_filter_none"
								http_filter_value = []
							}
						}
						acl {
							match {
								http_filter_value = []
								ip_subnet = ["0.0.0.0/0"]
							}
							action {
								type = "deny"
							}
						}

						acl {
							match {
								ip_subnet = ["0.0.0.0/0"]
								http_filter = "http_header_match"
								http_filter_value = ["example.com"]
								http_filter_option = "host"
							}

							action {
								type = "allow"
							}
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isACLCorrect(tt, "scaleway_lb_frontend.frt01", []*lbSDK.ACL{
						{
							Name: "test-acl",
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterACLHTTPFilterNone,
								HTTPFilterValue: []*string{},
								Invert:          true,
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
						},
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"0.0.0.0/0"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterPathBegin,
								HTTPFilterValue: scw.StringSlicePtr([]string{"criteria1", "criteria2"}),
								Invert:          true,
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
						},
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"0.0.0.0/0"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterPathBegin,
								HTTPFilterValue: scw.StringSlicePtr([]string{"criteria1", "criteria2"}),
								Invert:          false,
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
						},
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"0.0.0.0/0"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterACLHTTPFilterNone,
								HTTPFilterValue: []*string{},
								Invert:          false,
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
						},
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"0.0.0.0/0"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterACLHTTPFilterNone,
								HTTPFilterValue: []*string{},
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeDeny},
						},
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:         scw.StringSlicePtr([]string{"0.0.0.0/0"}),
								HTTPFilter:       lbSDK.ACLHTTPFilterHTTPHeaderMatch,
								HTTPFilterValue:  scw.StringSlicePtr([]string{"example.com"}),
								HTTPFilterOption: scw.StringPtr("host"),
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
						},
					}),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-acl"
						type = "lb-s"
					}
					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name = "tf-test"
						inbound_port = 80
						timeout_client = "30s"
						acl {
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["10.0.0.10"]
								http_filter = "path_begin"
								http_filter_value = ["foo","bar"]
							}
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isACLCorrect(tt, "scaleway_lb_frontend.frt01", []*lbSDK.ACL{
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"10.0.0.10"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterPathBegin,
								HTTPFilterValue: scw.StringSlicePtr([]string{"foo", "bar"}),
								Invert:          false,
							},
							Action: &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
						},
					}),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaleway_lb_ip.ip01", "id"),
				),
			},
		},
	})
}

func TestAccFrontend_ACLRedirectAction(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isFrontendDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-acl"
						type = "lb-s"
					}
					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name = "tf-test"
						inbound_port = 80
						timeout_client = "30s"
						acl {
							action {
								type = "redirect"
								redirect {
									type = "location"
									target = "https://example.com"
									code = 307
								}	
							}
							match {
								ip_subnet = ["10.0.0.10"]
								http_filter = "path_begin"
								http_filter_value = ["foo","bar"]
							}
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isACLCorrect(tt, "scaleway_lb_frontend.frt01", []*lbSDK.ACL{
						{
							Match: &lbSDK.ACLMatch{
								IPSubnet:        scw.StringSlicePtr([]string{"10.0.0.10"}),
								HTTPFilter:      lbSDK.ACLHTTPFilterPathBegin,
								HTTPFilterValue: scw.StringSlicePtr([]string{"foo", "bar"}),
								Invert:          false,
							},
							Action: &lbSDK.ACLAction{
								Type: lbSDK.ACLActionTypeRedirect,
								Redirect: &lbSDK.ACLActionRedirect{
									Type:   lbSDK.ACLActionRedirectRedirectTypeLocation,
									Target: "https://example.com",
									Code:   types.ExpandInt32Ptr(307),
								},
							},
						},
					}),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaleway_lb_ip.ip01", "id"),
				),
			},
		},
	})
}

func isACLCorrect(tt *acctest.TestTools, frontendName string, expectedAcls []*lbSDK.ACL) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// define a wrapper for acl comparison
		testCompareAcls := func(testAcl, apiAcl lbSDK.ACL) bool {
			// drop some values which are not part of the testing acl structure
			apiAcl.ID = ""
			apiAcl.Frontend = nil
			// if we do not pass any name, then drop it from comparison
			if testAcl.Name == "" {
				testAcl.Name = apiAcl.Name
			}
			return lb.ACLEquals(&testAcl, &apiAcl)
		}

		rs, ok := s.RootModule().Resources[frontendName]
		if !ok {
			return fmt.Errorf("resource not found: %s", frontendName)
		}

		if rs.Primary.ID == "" {
			return errors.New("resource id is not set")
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		// fetch our acls from the scaleway
		resACL, err := lbAPI.ListACLs(&lbSDK.ZonedAPIListACLsRequest{
			Zone:       zone,
			FrontendID: ID,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error on getting acl list [%s]", err)
		}

		// verify that the count of api acl is the same as we are expecting it to be
		if len(expectedAcls) != len(resACL.ACLs) {
			return errors.New("acl count is wrong")
		}
		// convert them to map indexed by the acl index
		aclMap := make(map[int32]*lbSDK.ACL)
		for _, acl := range resACL.ACLs {
			aclMap[acl.Index] = acl
		}

		// check that every index is set up correctly
		for i := 1; i <= len(expectedAcls); i++ {
			if _, found := aclMap[int32(i)]; !found {
				return fmt.Errorf("cannot find an index set [%d]", i)
			}
			if !testCompareAcls(*expectedAcls[i-1], *aclMap[int32(i)]) {
				return fmt.Errorf("two acls are not equal on stage %d", i)
			}
		}
		// check the actual data

		return nil
	}
}

func TestAclEqual(t *testing.T) {
	aclA := &lbSDK.ACL{
		Name: "test-acl",
		Match: &lbSDK.ACLMatch{
			IPSubnet:        scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24"}),
			HTTPFilter:      lbSDK.ACLHTTPFilterACLHTTPFilterNone,
			HTTPFilterValue: nil,
			Invert:          true,
		},
		Action:   &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
		Frontend: nil,
		Index:    1,
	}
	aclB := &lbSDK.ACL{
		Name: "test-acl",
		Match: &lbSDK.ACLMatch{
			IPSubnet:        scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24"}),
			HTTPFilter:      lbSDK.ACLHTTPFilterACLHTTPFilterNone,
			HTTPFilterValue: nil,
			Invert:          true,
		},
		Action:   &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow},
		Frontend: nil,
		Index:    1,
	}
	assert.True(t, lb.ACLEquals(aclA, aclB))

	// change name
	aclA.Name = "nope"
	assert.False(t, lb.ACLEquals(aclA, aclB))
	aclA.Name = aclB.Name

	// check action
	aclA.Action = nil
	assert.False(t, lb.ACLEquals(aclA, aclB))
	aclA.Action = &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow}
	assert.True(t, lb.ACLEquals(aclA, aclB))
	aclA.Action = &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeDeny}
	assert.False(t, lb.ACLEquals(aclA, aclB))
	aclA.Action = &lbSDK.ACLAction{Type: lbSDK.ACLActionTypeAllow}
	assert.True(t, lb.ACLEquals(aclA, aclB))

	// check match
	aclA.Match.IPSubnet = scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24", "0.0.0.0"})
	assert.False(t, lb.ACLEquals(aclA, aclB))
}
