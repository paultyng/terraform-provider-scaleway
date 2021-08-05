package scaleway

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/namegenerator"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/xerrors"
)

var (
	// DefaultWaitRetryInterval is used to set the retry interval to 0 during acceptance tests
	DefaultWaitRetryInterval *time.Duration
)

// RegionalID represents an ID that is linked with a region, eg fr-par/11111111-1111-1111-1111-111111111111
type RegionalID struct {
	ID     string
	Region scw.Region
}

func (z RegionalID) String() string {
	return fmt.Sprintf("%s/%s", z.Region, z.ID)
}

func newRegionalID(region scw.Region, id string) RegionalID {
	return RegionalID{
		ID:     id,
		Region: region,
	}
}

func expandRegionalID(id interface{}) RegionalID {
	regionalID := RegionalID{}
	tab := strings.SplitN(id.(string), "/", -1)
	if len(tab) != 2 {
		regionalID.ID = id.(string)
	} else {
		region, _ := scw.ParseRegion(tab[0])
		regionalID.ID = tab[1]
		regionalID.Region = region
	}

	return regionalID
}

// ZonedID represents an ID that is linked with a zone, eg fr-par-1/11111111-1111-1111-1111-111111111111
type ZonedID struct {
	ID   string
	Zone scw.Zone
}

func (z ZonedID) String() string {
	return fmt.Sprintf("%s/%s", z.Zone, z.ID)
}

func newZonedID(zone scw.Zone, id string) ZonedID {
	return ZonedID{
		ID:   id,
		Zone: zone,
	}
}

func expandZonedID(id interface{}) ZonedID {
	zonedID := ZonedID{}
	tab := strings.SplitN(id.(string), "/", -1)
	if len(tab) != 2 {
		zonedID.ID = id.(string)
	} else {
		zone, _ := scw.ParseZone(tab[0])
		zonedID.ID = tab[1]
		zonedID.Zone = zone
	}

	return zonedID
}

// parseLocalizedID parses a localizedID and extracts the resource locality and id.
func parseLocalizedID(localizedID string) (locality string, ID string, err error) {
	tab := strings.SplitN(localizedID, "/", -1)
	if len(tab) != 2 {
		return "", localizedID, fmt.Errorf("cant parse localized id: %s", localizedID)
	}
	return tab[0], tab[1], nil
}

// parseLocalizedNestedID parses a localizedNestedID and extracts the resource locality, the inner and outer id.
func parseLocalizedNestedID(localizedID string) (locality string, innerID, outerID string, err error) {
	tab := strings.SplitN(localizedID, "/", -1)
	if len(tab) != 3 {
		return "", "", localizedID, fmt.Errorf("cant parse localized id: %s", localizedID)
	}
	return tab[0], tab[1], tab[2], nil
}

// parseZonedID parses a zonedID and extracts the resource zone and id.
func parseZonedID(zonedID string) (zone scw.Zone, id string, err error) {
	locality, id, err := parseLocalizedID(zonedID)
	if err != nil {
		return
	}

	zone, err = scw.ParseZone(locality)
	return
}

// parseZonedNestedID parses a zonedNestedID and extracts the resource zone ,inner and outer ID.
func parseZonedNestedID(zonedNestedID string) (zone scw.Zone, outerID, innerID string, err error) {
	locality, innerID, outerID, err := parseLocalizedNestedID(zonedNestedID)
	if err != nil {
		return
	}

	zone, err = scw.ParseZone(locality)
	return
}

// expandID returns the id whether it is a localizedID or a raw ID.
func expandID(id interface{}) string {
	_, ID, err := parseLocalizedID(id.(string))
	if err != nil {
		return id.(string)
	}
	return ID
}

// parseRegionalID parses a regionalID and extracts the resource region and id.
func parseRegionalID(regionalID string) (region scw.Region, id string, err error) {
	locality, id, err := parseLocalizedID(regionalID)
	if err != nil {
		return
	}

	region, err = scw.ParseRegion(locality)
	return
}

// newZonedIDString constructs a unique identifier based on resource zone and id
func newZonedIDString(zone scw.Zone, id string) string {
	return fmt.Sprintf("%s/%s", zone, id)
}

// newZonedNestedIDString constructs a unique identifier based on resource zone, inner and outer IDs
func newZonedNestedIDString(zone scw.Zone, outerID, innerID string) string {
	return fmt.Sprintf("%s/%s/%s", zone, outerID, innerID)
}

// newRegionalIDString constructs a unique identifier based on resource region and id
func newRegionalIDString(region scw.Region, id string) string {
	return fmt.Sprintf("%s/%s", region, id)
}

// terraformResourceData is an interface for *schema.ResourceData. (used for mock)
type terraformResourceData interface {
	HasChange(string) bool
	GetOkExists(string) (interface{}, bool)
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
	Set(string, interface{}) error
	SetId(string)
	Id() string
}

// ErrZoneNotFound is returned when no zone can be detected
var ErrZoneNotFound = fmt.Errorf("could not detect zone. Scaleway uses regions and zones. For more information, refer to https://www.terraform.io/docs/providers/scaleway/guides/regions_and_zones.html")

// extractZone will try to guess the zone from the following:
//  - zone field of the resource data
//  - default zone from config
func extractZone(d terraformResourceData, meta *Meta) (scw.Zone, error) {
	rawZone, exist := d.GetOkExists("zone")
	if exist {
		return scw.ParseZone(rawZone.(string))
	}

	zone, exist := meta.scwClient.GetDefaultZone()
	if exist {
		return zone, nil
	}

	return "", ErrZoneNotFound
}

// ErrRegionNotFound is returned when no region can be detected
var ErrRegionNotFound = fmt.Errorf("could not detect region")

// extractRegion will try to guess the region from the following:
//  - region field of the resource data
//  - default region from config
func extractRegion(d terraformResourceData, meta *Meta) (scw.Region, error) {
	rawRegion, exist := d.GetOkExists("region")
	if exist {
		return scw.ParseRegion(rawRegion.(string))
	}

	region, exist := meta.scwClient.GetDefaultRegion()
	if exist {
		return region, nil
	}

	return "", ErrRegionNotFound
}

// isHTTPCodeError returns true if err is an http error with code statusCode
func isHTTPCodeError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	responseError := &scw.ResponseError{}
	if xerrors.As(err, &responseError) && responseError.StatusCode == statusCode {
		return true
	}
	return false
}

// is404Error returns true if err is an HTTP 404 error
func is404Error(err error) bool {
	notFoundError := &scw.ResourceNotFoundError{}
	return isHTTPCodeError(err, http.StatusNotFound) || xerrors.As(err, &notFoundError)
}

// is403Error returns true if err is an HTTP 403 error
func is403Error(err error) bool {
	permissionsDeniedError := &scw.PermissionsDeniedError{}
	return isHTTPCodeError(err, http.StatusForbidden) || xerrors.As(err, &permissionsDeniedError)
}

// organizationIDSchema returns a standard schema for a organization_id
func organizationIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The organization_id you want to attach the resource to",
		Computed:    true,
	}
}

// projectIDSchema returns a standard schema for a project_id
func projectIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The project_id you want to attach the resource to",
		Optional:     true,
		ForceNew:     true,
		Computed:     true,
		ValidateFunc: validationUUID(),
	}
}

// zoneSchema returns a standard schema for a zone
func zoneSchema() *schema.Schema {
	var allZones []string
	for _, z := range scw.AllZones {
		allZones = append(allZones, z.String())
	}
	return &schema.Schema{
		Type:             schema.TypeString,
		Description:      "The zone you want to attach the resource to",
		Optional:         true,
		ForceNew:         true,
		Computed:         true,
		ValidateDiagFunc: validateStringInSliceWithWarning(allZones, "zone"),
	}
}

// regionSchema returns a standard schema for a zone
func regionSchema() *schema.Schema {
	var allRegions []string
	for _, z := range scw.AllRegions {
		allRegions = append(allRegions, z.String())
	}
	return &schema.Schema{
		Type:             schema.TypeString,
		Description:      "The region you want to attach the resource to",
		Optional:         true,
		ForceNew:         true,
		Computed:         true,
		ValidateDiagFunc: validateStringInSliceWithWarning(allRegions, "region"),
	}
}

// regionComputedSchema returns a standard schema for a zone
func regionComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The region of the resource",
		Computed:    true,
	}
}

// validateStringInSliceWithWarning helps to only returns warnings in case we got a non public locality passed
func validateStringInSliceWithWarning(correctValues []string, field string) func(i interface{}, path cty.Path) diag.Diagnostics {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		_, rawErr := validation.StringInSlice(correctValues, true)(i, field)
		var res diag.Diagnostics
		for _, e := range rawErr {
			res = append(res, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       e.Error(),
				AttributePath: path,
			})
		}
		return res
	}
}

// newRandomName returns a random name prefixed for terraform.
func newRandomName(prefix string) string {
	return namegenerator.GetRandomName("tf", prefix)
}

const gb uint64 = 1000 * 1000 * 1000

func flattenTime(date *time.Time) interface{} {
	if date != nil {
		return date.Format(time.RFC3339)
	}
	return ""
}

func flattenDuration(duration *time.Duration) interface{} {
	if duration != nil {
		return duration.String()
	}
	return ""
}

func expandDuration(data interface{}) *time.Duration {
	if data == nil || data == "" {
		return nil
	}
	d, err := time.ParseDuration(data.(string))
	if err != nil {
		// We panic as this should never happened. Data from state should be validate using a validate func
		panic(err) // lintignore:R009
	}
	return &d
}

func expandOrGenerateString(data interface{}, prefix string) string {
	if data == nil || data == "" {
		return newRandomName(prefix)
	}
	return data.(string)
}

func expandStringWithDefault(data interface{}, defaultValue string) string {
	if data == nil || data.(string) == "" {
		return defaultValue
	}
	return data.(string)
}

func expandStrings(data interface{}) []string {
	stringSlice := []string{}
	for _, s := range data.([]interface{}) {
		stringSlice = append(stringSlice, s.(string))
	}
	return stringSlice
}

func expandStringsPtr(data interface{}) *[]string {
	var stringSlice []string
	if data == nil {
		return &stringSlice
	}
	for _, s := range data.([]interface{}) {
		stringSlice = append(stringSlice, s.(string))
	}
	return &stringSlice
}

func expandStringsOrEmpty(data interface{}) []string {
	if data == nil {
		return []string{}
	}
	stringSlice := []string{}
	for _, s := range data.([]interface{}) {
		stringSlice = append(stringSlice, s.(string))
	}
	return stringSlice
}

func expandSliceStringPtr(data interface{}) []*string {
	if data == nil {
		return nil
	}
	stringSlice := []*string(nil)
	for _, s := range data.([]interface{}) {
		stringSlice = append(stringSlice, expandStringPtr(s))
	}
	return stringSlice
}

func flattenIPPtr(ip *net.IP) interface{} {
	if ip == nil {
		return ""
	}
	return ip.String()
}

func flattenStringPtr(s *string) interface{} {
	if s == nil {
		return ""
	}
	return *s
}

func flattenSliceStringPtr(s []*string) interface{} {
	res := make([]interface{}, 0, len(s))
	for _, strPtr := range s {
		res = append(res, flattenStringPtr(strPtr))
	}
	return res
}

func expandStringPtr(data interface{}) *string {
	if data == nil || data == "" {
		return nil
	}
	return scw.StringPtr(data.(string))
}

func expandBoolPtr(data interface{}) *bool {
	if data == nil {
		return nil
	}
	return scw.BoolPtr(data.(bool))
}

func flattenInt32Ptr(i *int32) interface{} {
	if i == nil {
		return 0
	}
	return *i
}

func expandInt32Ptr(data interface{}) *int32 {
	if data == nil || data == "" {
		return nil
	}
	return scw.Int32Ptr(int32(data.(int)))
}

func expandIPNet(raw string) scw.IPNet {
	if raw == "" {
		return scw.IPNet{}
	}
	var ipNet scw.IPNet
	raw = `"` + raw + `"`
	err := json.Unmarshal([]byte(raw), &ipNet)
	if err != nil {
		// We panic as this should never happen. Data from state should be validate using a validate func
		panic(fmt.Errorf("%s could not be marshaled: %v", raw, err)) // lintignore:R009
	}

	return ipNet
}

func flattenIPNet(ipNet scw.IPNet) string {
	raw, err := json.Marshal(ipNet)
	if err != nil {
		// We panic as this should never happen.
		panic(err) // lintignore:R009
	}
	return string(raw[1 : len(raw)-1]) // remove quotes
}

func validateDuration() schema.SchemaValidateFunc {
	return func(i interface{}, s string) (strings []string, errors []error) {
		str, isStr := i.(string)
		if !isStr {
			return nil, []error{fmt.Errorf("%v is not a string", i)}
		}
		_, err := time.ParseDuration(str)
		if err != nil {
			return nil, []error{fmt.Errorf("cannot parse duration for value %s", str)}
		}
		return nil, nil
	}
}

func diffSuppressFuncDuration(k, old, new string, d *schema.ResourceData) bool {
	if old == new {
		return true
	}
	d1, err1 := time.ParseDuration(old)
	d2, err2 := time.ParseDuration(new)
	if err1 != nil || err2 != nil {
		return false
	}
	return d1 == d2
}

func diffSuppressFuncIgnoreCase(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func diffSuppressFuncIgnoreCaseAndHyphen(k, old, new string, d *schema.ResourceData) bool {
	return strings.Replace(strings.ToLower(old), "-", "_", -1) == strings.Replace(strings.ToLower(new), "-", "_", -1)
}

// diffSuppressFuncLocality is a SuppressDiffFunc to remove the locality from an ID when checking diff.
// e.g. 2c1a1716-5570-4668-a50a-860c90beabf6 == fr-par-1/2c1a1716-5570-4668-a50a-860c90beabf6
func diffSuppressFuncLocality(k, old, new string, d *schema.ResourceData) bool {
	return expandID(old) == expandID(new)
}
