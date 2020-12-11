package scaleway

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Deprecated
	testAccProviders map[string]*schema.Provider
	// Deprecated
	testAccProvider *schema.Provider

	// UpdateCassettes will update all cassettes of a given test
	UpdateCassettes = flag.Bool("cassettes", os.Getenv("TF_UPDATE_CASSETTES") == "true", "Record Cassettes")
)

func init() {
	p := Provider(DefaultProviderConfig())()
	testAccProvider = p
	Version += "-tftest"
	testAccProviders = map[string]*schema.Provider{
		"scaleway": p,
	}
}

func testAccPreCheck(_ *testing.T) {}

// getTestFilePath returns a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
func getTestFilePath(t *testing.T, suffix string) string {
	specialChars := regexp.MustCompile(`[\\?%*:|"<>. ]`)

	// Replace nested tests separators.
	fileName := strings.Replace(t.Name(), "/", "-", -1)

	fileName = strcase.ToBashArg(fileName)

	// Replace special characters.
	fileName = specialChars.ReplaceAllLiteralString(fileName, "") + suffix

	// Remove prefix to simplify
	fileName = strings.TrimPrefix(fileName, "test-acc-scaleway-")

	return filepath.Join(".", "testdata", fileName)
}

// getHTTPRecoder creates a new httpClient that records all HTTP requests in a cassette.
// This cassette is then replayed whenever tests are executed again. This means that once the
// requests are recorded in the cassette, no more real HTTP requests must be made to run the tests.
//
// It is important to add a `defer cleanup()` so the given cassette files are correctly
// closed and saved after the requests.
func getHTTPRecoder(t *testing.T, update bool) (client *http.Client, cleanup func(), err error) {
	recorderMode := recorder.ModeReplaying
	if update {
		recorderMode = recorder.ModeRecording
	}

	// Setup recorder and scw client
	r, err := recorder.NewAsMode(getTestFilePath(t, ".cassette"), recorderMode, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add a filter which removes Authorization headers from all requests:
	r.AddFilter(func(i *cassette.Interaction) error {
		i.Request.Headers = i.Request.Headers.Clone()
		delete(i.Request.Headers, "x-auth-token")
		delete(i.Request.Headers, "X-Auth-Token")
		return nil
	})

	return &http.Client{Transport: newRetryableTransport(r)}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}

type TestTools struct {
	T                 *testing.T
	Meta              *Meta
	ProviderFactories map[string]func() (*schema.Provider, error)
	Cleanup           func()
	ctx               context.Context
}

func NewTestTools(t *testing.T) *TestTools {
	// Create an http client with recording capabilities
	httpClient, cleanup, err := getHTTPRecoder(t, *UpdateCassettes)
	require.NoError(t, err)

	// Create meta that will be passed in the provider config
	meta, err := buildMeta(&MetaConfig{
		providerSchema:   nil,
		terraformVersion: "terraform-tests",
		httpClient:       httpClient,
	})
	require.NoError(t, err)

	return &TestTools{
		T:    t,
		Meta: meta,
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"scaleway": func() (*schema.Provider, error) {
				return Provider(&ProviderConfig{Meta: meta})(), nil
			},
		},
		Cleanup: cleanup,
		ctx:     context.Background(),
	}
}
