package lb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccCertificate_Basic(t *testing.T) {
	/**
	* See the discussion on https://github.com/scaleway/terraform-provider-scaleway/pull/396
	* Long story short, scaleway API will not permit you to request a certificate in case common name is not pointed
	* to the load balancer IP (which is unknown before creating it). In production, this can be overcome by introducing
	* an additional step which creates a DNS record and depending on it
	* We use a DNS name like: 195-154-70-235.lb.fr-par.scw.cloud
	 */
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isLbDestroyed(tt),
			lbchecks.IsIPDestroyed(tt),
			isCertificateDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {
					}
			
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}
			
					resource scaleway_lb_certificate cert01 {
						lb_id = scaleway_lb.lb01.id
						name = "test-cert"
					  	letsencrypt {
							common_name = "${replace(scaleway_lb_ip.ip01.ip_address,".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud"
					  	}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "name", "test-cert"),
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "letsencrypt.#", "1"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {
					}
			
					resource scaleway_lb lb01 {
						ip_id  = scaleway_lb_ip.ip01.id
						name   = "test-lb"
						type   = "lb-s"
					}
			
					resource scaleway_lb_certificate cert01 {
						lb_id = scaleway_lb.lb01.id
						name = "test-cert-new"
						letsencrypt {
							common_name = "${replace(scaleway_lb.lb01.ip_address, ".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "name", "test-cert-new"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {
					}

					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_lb_certificate cert01 {
						lb_id = scaleway_lb.lb01.id
						name = "test-cert"
						letsencrypt {
							common_name = "${replace(scaleway_lb.lb01.ip_address, ".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud"
							subject_alternative_name = [
							  "sub1.${replace(scaleway_lb.lb01.ip_address, ".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud",
							  "sub2.${replace(scaleway_lb.lb01.ip_address, ".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud"
							]
						}
					}
							`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "name", "test-cert"),
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "letsencrypt.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "letsencrypt.0.subject_alternative_name.#", "2"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {
					}

					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}

					resource scaleway_lb_certificate cert01 {
						lb_id = scaleway_lb.lb01.id
						name = "test-custom-cert"
					  	custom_certificate {
							certificate_chain = <<EOF
-----BEGIN CERTIFICATE-----
MIIEeDCCA2CgAwIBAgISA+78x4/3radnfUGMWDp4jLA+MA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yMTEyMjAxMzU1MDlaFw0yMjAzMjAxMzU1MDhaMCkxJzAlBgNVBAMT
HnN0YWdpbmcuc2NhbGV3YXktdGVycmFmb3JtLmNvbTBZMBMGByqGSM49AgEGCCqG
SM49AwEHA0IABPqtv8kzJW2QI3LG93ScS13SSm5SIJI2uAoIhny4Lj06g05ff7mM
gmd+XgID6OmWOB+Paso3Udq6hccEUN/URlKjggJaMIICVjAOBgNVHQ8BAf8EBAMC
B4AwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAw
HQYDVR0OBBYEFJvpLZyegSLA1uGb1o2QvthXFb+OMB8GA1UdIwQYMBaAFBQusxe3
WFbLrlAJQOYfr52LFMLGMFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0
cDovL3IzLm8ubGVuY3Iub3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5j
ci5vcmcvMCkGA1UdEQQiMCCCHnN0YWdpbmcuc2NhbGV3YXktdGVycmFmb3JtLmNv
bTBMBgNVHSAERTBDMAgGBmeBDAECATA3BgsrBgEEAYLfEwEBATAoMCYGCCsGAQUF
BwIBFhpodHRwOi8vY3BzLmxldHNlbmNyeXB0Lm9yZzCCAQUGCisGAQQB1nkCBAIE
gfYEgfMA8QB2AN+lXqtogk8fbK3uuF9OPlrqzaISpGpejjsSwCBEXCpzAAABfdhW
1YQAAAQDAEcwRQIgIFwFx6ihGAbrAGEtXsncJB2JGg899FAuQIS4tar08UgCIQCY
xmuSCme644NT4cq1HQYaO0EX0FvIgXRGULciNzsoGwB3AEalVet1+pEgMLWiiWn0
830RLEF0vv1JuIWr8vxw/m1HAAABfdhW1awAAAQDAEgwRgIhAKkjKvUYsQqqn1rI
z/I9ryz3hwUBIrnh33kCAthLXXwNAiEAhM7UwVoljG5UK9g+QxckiL9sWn+W+z+8
YLoJJQhrKbwwDQYJKoZIhvcNAQELBQADggEBAKbLFTFyWyXkfPqI0x4Z2I7PHXH4
w7J5E1fqVLyg3rIojCSzwsYzbOhVEjZ5CK8W3LKbG9s/kOBqmoNgoGgT4thZg4Ks
uZz9am4AeKWz0z6SjUgqLt53UCIj5VbefccOkuuqPa0l5sCAo3gCeV9BfwHmAg3m
zPkjEfqNkEzot3tei/M4RRDn7caURghxbOG9tm8c1aCe/gsq0G0j1KRM8zwXCDQ2
nKzB5TIJsAHKlIkUfyheK+xG+B4F2f5mfdY4l4a/mQH4+nLcWdMly1DzHMrdyC8c
6m2o9DBaRsXON54j8Fdgct/C5P3r0IjK6ne2PEiLhi/WYFiw60/sNjQHKsQ=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFFjCCAv6gAwIBAgIRAJErCErPDBinU/bWLiWnX1owDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjAwOTA0MDAwMDAw
WhcNMjUwOTE1MTYwMDAwWjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
RW5jcnlwdDELMAkGA1UEAxMCUjMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQC7AhUozPaglNMPEuyNVZLD+ILxmaZ6QoinXSaqtSu5xUyxr45r+XXIo9cP
R5QUVTVXjJ6oojkZ9YI8QqlObvU7wy7bjcCwXPNZOOftz2nwWgsbvsCUJCWH+jdx
sxPnHKzhm+/b5DtFUkWWqcFTzjTIUu61ru2P3mBw4qVUq7ZtDpelQDRrK9O8Zutm
NHz6a4uPVymZ+DAXXbpyb/uBxa3Shlg9F8fnCbvxK/eG3MHacV3URuPMrSXBiLxg
Z3Vms/EY96Jc5lP/Ooi2R6X/ExjqmAl3P51T+c8B5fWmcBcUr2Ok/5mzk53cU6cG
/kiFHaFpriV1uxPMUgP17VGhi9sVAgMBAAGjggEIMIIBBDAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMBIGA1UdEwEB/wQIMAYB
Af8CAQAwHQYDVR0OBBYEFBQusxe3WFbLrlAJQOYfr52LFMLGMB8GA1UdIwQYMBaA
FHm0WeZ7tuXkAXOACIjIGlj26ZtuMDIGCCsGAQUFBwEBBCYwJDAiBggrBgEFBQcw
AoYWaHR0cDovL3gxLmkubGVuY3Iub3JnLzAnBgNVHR8EIDAeMBygGqAYhhZodHRw
Oi8veDEuYy5sZW5jci5vcmcvMCIGA1UdIAQbMBkwCAYGZ4EMAQIBMA0GCysGAQQB
gt8TAQEBMA0GCSqGSIb3DQEBCwUAA4ICAQCFyk5HPqP3hUSFvNVneLKYY611TR6W
PTNlclQtgaDqw+34IL9fzLdwALduO/ZelN7kIJ+m74uyA+eitRY8kc607TkC53wl
ikfmZW4/RvTZ8M6UK+5UzhK8jCdLuMGYL6KvzXGRSgi3yLgjewQtCPkIVz6D2QQz
CkcheAmCJ8MqyJu5zlzyZMjAvnnAT45tRAxekrsu94sQ4egdRCnbWSDtY7kh+BIm
lJNXoB1lBMEKIq4QDUOXoRgffuDghje1WrG9ML+Hbisq/yFOGwXD9RiX8F6sw6W4
avAuvDszue5L3sz85K+EC4Y/wFVDNvZo4TYXao6Z0f+lQKc0t8DQYzk1OXVu8rp2
yJMC6alLbBfODALZvYH7n7do1AZls4I9d1P4jnkDrQoxB3UqQ9hVl3LEKQ73xF1O
yK5GhDDX8oVfGKF5u+decIsH4YaTw7mP3GFxJSqv3+0lUFJoi5Lc5da149p90Ids
hCExroL1+7mryIkXPeFM5TgO9r0rvZaBFOvV2z0gp35Z0+L4WPlbuEjN/lxPFin+
HlUjr8gRsI3qfJOQFy/9rKIJR0Y/8Omwt/8oTWgy1mdeHmmjk7j1nYsvC9JSQ6Zv
MldlTTKB3zhThV1+XWYp6rjd5JW1zbVWEkLNxE7GJThEUG3szgBVGP7pSWTUTsqX
nLRbwHOoq7hHwg==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFYDCCBEigAwIBAgIQQAF3ITfU6UK47naqPGQKtzANBgkqhkiG9w0BAQsFADA/
MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMT
DkRTVCBSb290IENBIFgzMB4XDTIxMDEyMDE5MTQwM1oXDTI0MDkzMDE4MTQwM1ow
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQCt6CRz9BQ385ueK1coHIe+3LffOJCMbjzmV6B493XC
ov71am72AE8o295ohmxEk7axY/0UEmu/H9LqMZshftEzPLpI9d1537O4/xLxIZpL
wYqGcWlKZmZsj348cL+tKSIG8+TA5oCu4kuPt5l+lAOf00eXfJlII1PoOK5PCm+D
LtFJV4yAdLbaL9A4jXsDcCEbdfIwPPqPrt3aY6vrFk/CjhFLfs8L6P+1dy70sntK
4EwSJQxwjQMpoOFTJOwT2e4ZvxCzSow/iaNhUd6shweU9GNx7C7ib1uYgeGJXDR5
bHbvO5BieebbpJovJsXQEOEO3tkQjhb7t/eo98flAgeYjzYIlefiN5YNNnWe+w5y
sR2bvAP5SQXYgd0FtCrWQemsAXaVCg/Y39W9Eh81LygXbNKYwagJZHduRze6zqxZ
Xmidf3LWicUGQSk+WT7dJvUkyRGnWqNMQB9GoZm1pzpRboY7nn1ypxIFeFntPlF4
FQsDj43QLwWyPntKHEtzBRL8xurgUBN8Q5N0s8p0544fAQjQMNRbcTa0B7rBMDBc
SLeCO5imfWCKoqMpgsy6vYMEG6KDA0Gh1gXxG8K28Kh8hjtGqEgqiNx2mna/H2ql
PRmP6zjzZN7IKw0KKP/32+IVQtQi0Cdd4Xn+GOdwiK1O5tmLOsbdJ1Fu/7xk9TND
TwIDAQABo4IBRjCCAUIwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYw
SwYIKwYBBQUHAQEEPzA9MDsGCCsGAQUFBzAChi9odHRwOi8vYXBwcy5pZGVudHJ1
c3QuY29tL3Jvb3RzL2RzdHJvb3RjYXgzLnA3YzAfBgNVHSMEGDAWgBTEp7Gkeyxx
+tvhS5B1/8QVYIWJEDBUBgNVHSAETTBLMAgGBmeBDAECATA/BgsrBgEEAYLfEwEB
ATAwMC4GCCsGAQUFBwIBFiJodHRwOi8vY3BzLnJvb3QteDEubGV0c2VuY3J5cHQu
b3JnMDwGA1UdHwQ1MDMwMaAvoC2GK2h0dHA6Ly9jcmwuaWRlbnRydXN0LmNvbS9E
U1RST09UQ0FYM0NSTC5jcmwwHQYDVR0OBBYEFHm0WeZ7tuXkAXOACIjIGlj26Ztu
MA0GCSqGSIb3DQEBCwUAA4IBAQAKcwBslm7/DlLQrt2M51oGrS+o44+/yQoDFVDC
5WxCu2+b9LRPwkSICHXM6webFGJueN7sJ7o5XPWioW5WlHAQU7G75K/QosMrAdSW
9MUgNTP52GE24HGNtLi1qoJFlcDyqSMo59ahy2cI2qBDLKobkx/J3vWraV0T9VuG
WCLKTVXkcGdtwlfFRjlBz4pYg1htmf5X6DYO8A4jqv2Il9DjXA6USbW1FzXSLr9O
he8Y4IWS6wY7bCkjCWDcRQJMEhg76fsO3txE+FiYruq9RUWhiF1myv4Q6W+CyBFC
Dfvp7OOGAN6dEOM4+qR9sdjoSYKEBpsr6GtPAQw4dy753ec5
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEII45PLJDlsxyXYk55ladXnPUQOwYphEOy3Z3qlt5EoRBoAoGCCqGSM49
AwEHoUQDQgAE+q2/yTMlbZAjcsb3dJxLXdJKblIgkja4CgiGfLguPTqDTl9/uYyC
Z35eAgPo6ZY4H49qyjdR2rqFxwRQ39RGUg==
-----END EC PRIVATE KEY-----
EOF
					  	}
					} 
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "name", "test-custom-cert"),
					resource.TestCheckResourceAttr("scaleway_lb_certificate.cert01", "custom_certificate.#", "1"),
				),
			},
		},
	})
}

func isCertificateDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_certificate" {
				continue
			}

			lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = lbAPI.GetCertificate(&lbSDK.ZonedAPIGetCertificateRequest{
				CertificateID: ID,
				Zone:          zone,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("LB Certificate (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
