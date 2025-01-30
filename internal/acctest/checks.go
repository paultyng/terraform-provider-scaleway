package acctest

import (
	"errors"
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
)

// CheckResourceIDChanged checks that the ID of the resource has indeed changed, in case of ForceNew for example.
// It will fail if resourceID is empty so be sure to use acctest.CheckResourceIDPersisted first in a test suite.
func CheckResourceIDChanged(resourceName string, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if resourceID == nil || *resourceID == "" {
			return errors.New("resourceID was not set")
		}
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource was not found: %s", resourceName)
		}
		if *resourceID == rs.Primary.ID {
			return errors.New("resource ID persisted when it should have changed")
		}
		*resourceID = rs.Primary.ID
		return nil
	}
}

// CheckResourceIDPersisted checks that the ID of the resource is the same throughout tests of migration or mutation
// It can be used to check that no ForceNew has been done
func CheckResourceIDPersisted(resourceName string, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource was not found: %s", resourceName)
		}
		if *resourceID != "" && *resourceID != rs.Primary.ID {
			return errors.New("resource ID changed when it should have persisted")
		}
		*resourceID = rs.Primary.ID
		return nil
	}
}

// CheckResourceRawIDMatches asserts the equality of IDs from two specified attributes of two Scaleway resources.
func CheckResourceRawIDMatches(res1, attr1, res2, attr2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, ok1 := s.RootModule().Resources[res1]
		if !ok1 {
			return fmt.Errorf("not found: %s", res1)
		}

		rs2, ok2 := s.RootModule().Resources[res2]
		if !ok2 {
			return fmt.Errorf("not found: %s", res2)
		}

		id1 := locality.ExpandID(rs1.Primary.Attributes[attr1])
		id2 := locality.ExpandID(rs2.Primary.Attributes[attr2])

		if id1 != id2 {
			return fmt.Errorf("ID mismatch: %s from resource %s does not match ID %s from resource %s", id1, res1, id2, res2)
		}

		return nil
	}
}

func CheckResourceAttrUUID(name string, key string) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(name, key, UUIDRegex)
}

var UUIDRegex = regexp.MustCompile(`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

func CheckResourceAttrFunc(name string, key string, test func(string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}
		value, ok := rs.Primary.Attributes[key]
		if !ok {
			return fmt.Errorf("key not found: %s", key)
		}
		err := test(value)
		if err != nil {
			return fmt.Errorf("test for %s %s did not pass test: %s", name, key, err)
		}
		return nil
	}
}

func CheckResourceAttrIPv4(name string, key string) resource.TestCheckFunc {
	return CheckResourceAttrFunc(name, key, func(value string) error {
		ip := net.ParseIP(value)
		if ip.To4() == nil {
			return fmt.Errorf("%s is not a valid IPv4", value)
		}
		return nil
	})
}

func CheckResourceAttrIPv6(name string, key string) resource.TestCheckFunc {
	return CheckResourceAttrFunc(name, key, func(value string) error {
		ip := net.ParseIP(value)
		if ip.To16() == nil {
			return fmt.Errorf("%s is not a valid IPv6", value)
		}
		return nil
	})
}

func CheckResourceAttrIP(name string, key string) resource.TestCheckFunc {
	return CheckResourceAttrFunc(name, key, func(value string) error {
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("%s is not a valid IP", value)
		}
		return nil
	})
}
