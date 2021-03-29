package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestAccScalewayLbFrontend_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayLbFrontendDestroy(tt),
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
					testAccCheckScalewayLbFrontendExists(tt, "scaleway_lb_frontend.frt01"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "inbound_port", "80"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "timeout_client", ""),
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
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayLbFrontendExists(tt, "scaleway_lb_frontend.frt01"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "name", "tf-test"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "inbound_port", "443"),
					resource.TestCheckResourceAttr("scaleway_lb_frontend.frt01", "timeout_client", "30s"),
				),
			},
		},
	})
}

func testAccCheckScalewayLbFrontendExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, region, ID, err := lbAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = lbAPI.GetFrontend(&lb.GetFrontendRequest{
			FrontendID: ID,
			Region:     region,
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayLbFrontendDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_frontend" {
				continue
			}

			lbAPI, region, ID, err := lbAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = lbAPI.GetFrontend(&lb.GetFrontendRequest{
				Region:     region,
				FrontendID: ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("LB Frontend (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}

func TestAclEqual(t *testing.T) {
	aclA := &lb.ACL{
		Name: "test-acl",
		Match: &lb.ACLMatch{
			IPSubnet:        scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24"}),
			HTTPFilter:      lb.ACLHTTPFilterACLHTTPFilterNone,
			HTTPFilterValue: nil,
			Invert:          true,
		},
		Action:   &lb.ACLAction{Type: lb.ACLActionTypeAllow},
		Frontend: nil,
		Index:    1,
	}
	aclB := &lb.ACL{
		Name: "test-acl",
		Match: &lb.ACLMatch{
			IPSubnet:        scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24"}),
			HTTPFilter:      lb.ACLHTTPFilterACLHTTPFilterNone,
			HTTPFilterValue: nil,
			Invert:          true,
		},
		Action:   &lb.ACLAction{Type: lb.ACLActionTypeAllow},
		Frontend: nil,
		Index:    1,
	}
	assert.True(t, aclEquals(aclA, aclB))

	//change name
	aclA.Name = "nope"
	assert.False(t, aclEquals(aclA, aclB))
	aclA.Name = aclB.Name

	//check action
	aclA.Action = nil
	assert.False(t, aclEquals(aclA, aclB))
	aclA.Action = &lb.ACLAction{Type: lb.ACLActionTypeAllow}
	assert.True(t, aclEquals(aclA, aclB))
	aclA.Action = &lb.ACLAction{Type: lb.ACLActionTypeDeny}
	assert.False(t, aclEquals(aclA, aclB))
	aclA.Action = &lb.ACLAction{Type: lb.ACLActionTypeAllow}
	assert.True(t, aclEquals(aclA, aclB))

	//check match
	aclA.Match.IPSubnet = scw.StringSlicePtr([]string{"192.168.0.1", "192.168.0.2", "192.168.10.0/24", "0.0.0.0"})
	assert.False(t, aclEquals(aclA, aclB))
}
