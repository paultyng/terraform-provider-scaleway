package iot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iotSDK "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iot"
)

const customDevCert = `-----BEGIN CERTIFICATE-----
MIIDkjCCAnqgAwIBAgIUJ/Xxw1ucfzPkzOTQFeojiTd9WDMwDQYJKoZIhvcNAQEL
BQAwRzELMAkGA1UEBhMCRlIxDjAMBgNVBAcTBVBhcmlzMREwDwYDVQQKEwhTY2Fs
ZXdheTEVMBMGA1UECxMMc2NhbGV3YXkuY29tMB4XDTIxMDkwOTEzMDYwMFoXDTMx
MDkwNzEzMDYwMFowRzELMAkGA1UEBhMCRlIxDjAMBgNVBAcTBVBhcmlzMREwDwYD
VQQKEwhTY2FsZXdheTEVMBMGA1UECxMMc2NhbGV3YXkuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJ7LJtcEM9FCv3+4ZCIbB9p+0Nx0YbOZhGh1
X3cdALXoZ803uOvvO86rkYB73flHvRsSCgKyhD3uw2c4r3QJslRvljFYg2jAb5c0
A2rldJGsVxsv9FogY0i4sP8UA6ixCham9Tq5s0CD1VXJ+EaD92jp6FILhzJ7UBGD
PDLDPF73LcFTVjgNM8EQuQrkah38Et83j1Cqy/MLfrMWo6SY/oUyHTa1N9BZQHif
t6wxYCdV/i9JIRUmPL4w8TQRAURiMRjmAejnUyIekLhrrFm2W4R+p0WnVARMDrx3
THuxk5L3984l8n0ewwHRF+NGwxA09xqGrT+0ELcCmkhhOniS5wIDAQABo3YwdDAO
BgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIw
ADAdBgNVHQ4EFgQUr/KKPNxVmB/IZ/gGgRmpYcOKz+cwIAYDVR0RBBkwF4IVdGVz
dC10Zi5pb3Quc2N3LmNsb3VkMA0GCSqGSIb3DQEBCwUAA4IBAQAMxzVBFBw3U9fj
fUabLpk9+O/9iPlDJfW09c21P3iuI76CnxCLaOCAMctNtdMQSodaYnnA1w1A9+Oq
QH3B/ydZCHfVL/5FkUayYHG1uzsJir9IOtJ8QWEpeDaprO5XBqMtRfGzrz7fDB6x
uNolPpQkOwHtgALqnHJHpxnI49NoUDXM9ZHvk7YY4WE8gEfDwMN185k78qElf2d5
WuePu1khrEuTaXKyaiLD3pmxM86F/6Ho6V86mJpKXr/wmMU56TcKk9UURucVQZ1o
2m0aSh7KsPOWXiIKFpGQMXZZkdpHDOPFAeX2Z7PzLD7AwYSfgnhTmgxL1XqkoHG8
3Vu6zPcD
-----END CERTIFICATE-----
`

func TestAccDevice_Minimal(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_device" "default-1" {
							name = "default-1"
							hub_id = scaleway_iot_hub.minimal-1.id
						}
						resource "scaleway_iot_hub" "minimal-1" {
							name = "minimal-1"
							product_plan = "plan_shared"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal-1"),
					isDevicePresent(tt, "scaleway_iot_device.default-1"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-1", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-1", "hub_id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-1", "name"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-1", "allow_insecure", "false"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-1", "allow_multiple_connections", "false"),
				),
			},
		},
	})
}

func TestAccDevice_MessageFilters(t *testing.T) {
	t.Skip("Some checks seem to be flaky.")
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_device" "default-4" {
							name = "default-4"
							hub_id = scaleway_iot_hub.minimal-4.id
							message_filters {
								publish {
									policy = "reject"
									topics = ["1", "2", "3"]
								}
								subscribe {
									policy = "accept"
									topics = ["4", "5", "6"]
								}
							}
						}
						resource "scaleway_iot_device" "empty" {
							name = "empty"
							hub_id = scaleway_iot_hub.minimal-4.id
							message_filters {
								publish { }
								subscribe { }
							}
						}
						resource "scaleway_iot_hub" "minimal-4" {
							name = "minimal-4"
							product_plan = "plan_shared"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal-4"),
					isDevicePresent(tt, "scaleway_iot_device.default-4"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-4", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-4", "hub_id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-4", "name"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-4", "allow_insecure", "false"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-4", "allow_multiple_connections", "false"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-4", "message_filters.0.publish.0.policy", "reject"),
					// TODO: The following checks seem to be flaky.
					// resource.TestCheckResourceAttr("scaleway_iot_device.default-4", "message_filters.0.publish.0.topics.0", "1"),
					// resource.TestCheckResourceAttr("scaleway_iot_device.default-4", "message_filters.0.subscribe.0.policy", "accept"),
					// resource.TestCheckResourceAttr("scaleway_iot_device.default-4", "message_filters.0.subscribe.0.topics.0", "4"),
				),
			},
		},
	})
}

func TestAccDevice_AllowInsecure(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_device" "default-2" {
							name = "default-2"
							hub_id = scaleway_iot_hub.minimal-2.id
							allow_insecure = true
						}
						resource "scaleway_iot_hub" "minimal-2" {
							name = "minimal-2"
							product_plan = "plan_shared"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal-2"),
					isDevicePresent(tt, "scaleway_iot_device.default-2"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-2", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-2", "hub_id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-2", "name"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-2", "allow_insecure", "true"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-2", "allow_multiple_connections", "false"),
				),
			},
			{
				Config: `
						resource "scaleway_iot_device" "default-3" {
							name = "default-3"
							hub_id = scaleway_iot_hub.minimal-3.id
							allow_insecure = false
						}
						resource "scaleway_iot_hub" "minimal-3" {
							name = "minimal-3"
							product_plan = "plan_shared"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal-3"),
					isDevicePresent(tt, "scaleway_iot_device.default-3"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-3", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-3", "hub_id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-3", "name"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-3", "allow_insecure", "false"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-3", "allow_multiple_connections", "false"),
				),
			},
		},
	})
}

func TestAccDevice_Certificate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_device" "default-5" {
							name = "default-5"
							hub_id = scaleway_iot_hub.minimal-5.id
							allow_insecure = true
							certificate {
								crt = <<EOF
` + customDevCert + `EOF
							}
						}
						resource "scaleway_iot_hub" "minimal-5" {
							name = "minimal-5"
							product_plan = "plan_shared"
						}`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal-5"),
					isDevicePresent(tt, "scaleway_iot_device.default-5"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-5", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-5", "hub_id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_device.default-5", "name"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-5", "allow_insecure", "true"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-5", "allow_multiple_connections", "false"),
					resource.TestCheckResourceAttr("scaleway_iot_device.default-5", "certificate.0.crt", customDevCert),
				),
			},
		},
	})
}

func isDevicePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		iotAPI, region, deviceID, err := iot.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = iotAPI.GetDevice(&iotSDK.GetDeviceRequest{
			Region:   region,
			DeviceID: deviceID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
