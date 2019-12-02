package scaleway

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLocalizedID(t *testing.T) {

	testCases := []struct {
		name       string
		localityId string
		id         string
		locality   string
		err        string
	}{
		{
			name:       "simple",
			localityId: "fr-par-1/my-id",
			id:         "my-id",
			locality:   "fr-par-1",
		},
		{
			name:       "id with slashed",
			localityId: "fr-par-1/my/id",
			id:         "my/id",
			locality:   "fr-par-1",
		},
		{
			name:       "id with a region",
			localityId: "fr-par/my-id",
			id:         "my-id",
			locality:   "fr-par",
		},
		{
			name:       "empty",
			localityId: "",
			err:        "cant parse localized id: ",
		},
		{
			name:       "without locality",
			localityId: "my-id",
			err:        "cant parse localized id: my-id",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			locality, id, err := parseLocalizedID(testCase.localityId)
			if testCase.err != "" {
				require.EqualError(t, err, testCase.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.locality, locality)
				assert.Equal(t, testCase.id, id)
			}
		})
	}

}

func TestParseZonedID(t *testing.T) {

	testCases := []struct {
		name       string
		localityId string
		id         string
		zone       scw.Zone
		err        string
	}{
		{
			name:       "simple",
			localityId: "fr-par-1/my-id",
			id:         "my-id",
			zone:       scw.ZoneFrPar1,
		},
		{
			name:       "id with slashed",
			localityId: "fr-par-1/my/id",
			id:         "my/id",
			zone:       scw.ZoneFrPar1,
		},
		{
			name:       "empty",
			localityId: "",
			err:        "cant parse localized id: ",
		},
		{
			name:       "without locality",
			localityId: "my-id",
			err:        "cant parse localized id: my-id",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			zone, id, err := parseZonedID(testCase.localityId)
			if testCase.err != "" {
				require.EqualError(t, err, testCase.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.zone, zone)
				assert.Equal(t, testCase.id, id)
			}
		})
	}

}

func TestParseRegionID(t *testing.T) {

	testCases := []struct {
		name       string
		localityId string
		id         string
		region     scw.Region
		err        string
	}{
		{
			name:       "simple",
			localityId: "fr-par/my-id",
			id:         "my-id",
			region:     scw.RegionFrPar,
		},
		{
			name:       "id with slashed",
			localityId: "fr-par/my/id",
			id:         "my/id",
			region:     scw.RegionFrPar,
		},
		{
			name:       "empty",
			localityId: "",
			err:        "cant parse localized id: ",
		},
		{
			name:       "without locality",
			localityId: "my-id",
			err:        "cant parse localized id: my-id",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			region, id, err := parseRegionalID(testCase.localityId)
			if testCase.err != "" {
				require.EqualError(t, err, testCase.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.region, region)
				assert.Equal(t, testCase.id, id)
			}
		})
	}

}

func TestNewZonedId(t *testing.T) {
	assert.Equal(t, "fr-par-1/my-id", newZonedId(scw.ZoneFrPar1, "my-id"))
}

func TestNewRegionalId(t *testing.T) {
	assert.Equal(t, "fr-par/my-id", newRegionalId(scw.RegionFrPar, "my-id"))
}

func TestIsHTTPCodeError(t *testing.T) {
	assert.True(t, isHTTPCodeError(&scw.ResponseError{StatusCode: http.StatusBadRequest}, http.StatusBadRequest))
	assert.False(t, isHTTPCodeError(nil, http.StatusBadRequest))
	assert.False(t, isHTTPCodeError(&scw.ResponseError{StatusCode: http.StatusBadRequest}, http.StatusNotFound))
	assert.False(t, isHTTPCodeError(fmt.Errorf("not an http error"), http.StatusNotFound))
}

func TestIs404Error(t *testing.T) {
	assert.True(t, is404Error(&scw.ResponseError{StatusCode: http.StatusNotFound}))
	assert.False(t, is404Error(nil))
	assert.False(t, is404Error(&scw.ResponseError{StatusCode: http.StatusBadRequest}))
}

func TestIs403Error(t *testing.T) {
	assert.True(t, is403Error(&scw.ResponseError{StatusCode: http.StatusForbidden}))
	assert.False(t, is403Error(nil))
	assert.False(t, is403Error(&scw.ResponseError{StatusCode: http.StatusBadRequest}))
}

func TestGetRandomName(t *testing.T) {
	name := getRandomName("test")
	assert.True(t, strings.HasPrefix(name, "tf-test-"))
}

func TestDiffSuppressFuncLabelUUID(t *testing.T) {
	testCases := []struct {
		name    string
		old     string
		new     string
		isEqual bool
	}{
		{
			name:    "no label changes",
			old:     "foo/11111111-1111-1111-1111-111111111111",
			new:     "foo",
			isEqual: true,
		},
		{
			name:    "no UUID changes",
			old:     "foo/11111111-1111-1111-1111-111111111111",
			new:     "11111111-1111-1111-1111-111111111111",
			isEqual: true,
		},
		{
			name:    "UUID changes",
			old:     "foo/11111111-1111-1111-1111-111111111111",
			new:     "22222222-2222-2222-2222-222222222222",
			isEqual: false,
		},
		{
			name:    "To label",
			old:     "foo/11111111-1111-1111-1111-111111111111",
			new:     "foo",
			isEqual: true,
		},
		{
			name:    "To label change",
			old:     "foo/11111111-1111-1111-1111-111111111111",
			new:     "bar",
			isEqual: false,
		},
	}

	fakeResourceData := &schema.ResourceData{}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.isEqual, diffSuppressFuncLabelUUID("", c.old, c.new, fakeResourceData))
		})
	}
}

func testCheckResourceAttrFunc(name string, key string, test func(string) error) resource.TestCheckFunc {
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

func testCheckResourceAttrUUID(name string, key string) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(name, key, UUIDRegex)
}

func testCheckResourceAttrIPv4(name string, key string) resource.TestCheckFunc {
	return testCheckResourceAttrFunc(name, key, func(value string) error {
		ip := net.ParseIP(value)
		if ip.To4() == nil {
			return fmt.Errorf("%s is not a valid IPv4", value)
		}
		return nil
	})
}

func testCheckResourceAttrIPv6(name string, key string) resource.TestCheckFunc {
	return testCheckResourceAttrFunc(name, key, func(value string) error {
		ip := net.ParseIP(value)
		if ip.To16() == nil {
			return fmt.Errorf("%s is not a valid IPv6", value)
		}
		return nil
	})
}
