package container

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

const (
	defaultContainerNamespaceTimeout = 5 * time.Minute
	defaultContainerCronTimeout      = 5 * time.Minute
	defaultContainerTimeout          = 12*time.Minute + 30*time.Second
	defaultContainerDomainTimeout    = 10 * time.Minute
	DefaultContainerRetryInterval    = 5 * time.Second
	defaultTriggerRetryInterval      = 5 * time.Second
)

// newAPIWithRegion returns a new container API and the region.
func newAPIWithRegion(d *schema.ResourceData, m interface{}) (*container.API, scw.Region, error) {
	api := container.NewAPI(meta.ExtractScwClient(m))

	region, err := meta.ExtractRegion(d, m)
	if err != nil {
		return nil, "", err
	}
	return api, region, nil
}

// NewAPIWithRegionAndID returns a new container API, region and ID.
func NewAPIWithRegionAndID(m interface{}, id string) (*container.API, scw.Region, string, error) {
	api := container.NewAPI(meta.ExtractScwClient(m))

	region, id, err := regional.ParseID(id)
	if err != nil {
		return nil, "", "", err
	}
	return api, region, id, nil
}

func setCreateContainerRequest(d *schema.ResourceData, region scw.Region) (*container.CreateContainerRequest, error) {
	// required
	nameRaw := d.Get("name")
	namespaceID := d.Get("namespace_id")

	name := types.ExpandOrGenerateString(nameRaw.(string), "co")
	privacyType := d.Get("privacy")
	protocol := d.Get("protocol")
	httpOption := d.Get("http_option")

	req := &container.CreateContainerRequest{
		Region:      region,
		NamespaceID: locality.ExpandID(namespaceID),
		Name:        name,
		Privacy:     container.ContainerPrivacy(privacyType.(string)),
		Protocol:    container.ContainerProtocol(*types.ExpandStringPtr(protocol)),
		HTTPOption:  container.ContainerHTTPOption(httpOption.(string)),
	}

	// optional
	if envVariablesRaw, ok := d.GetOk("environment_variables"); ok {
		req.EnvironmentVariables = types.ExpandMapPtrStringString(envVariablesRaw)
	}

	if secretEnvVariablesRaw, ok := d.GetOk("secret_environment_variables"); ok {
		req.SecretEnvironmentVariables = expandContainerSecrets(secretEnvVariablesRaw)
	}

	if minScale, ok := d.GetOk("min_scale"); ok {
		req.MinScale = scw.Uint32Ptr(uint32(minScale.(int)))
	}

	if maxScale, ok := d.GetOk("max_scale"); ok {
		req.MaxScale = scw.Uint32Ptr(uint32(maxScale.(int)))
	}

	if memoryLimit, ok := d.GetOk("memory_limit"); ok {
		req.MemoryLimit = scw.Uint32Ptr(uint32(memoryLimit.(int)))
	}

	if cpuLimit, ok := d.GetOk("cpu_limit"); ok {
		req.CPULimit = scw.Uint32Ptr(uint32(cpuLimit.(int)))
	}

	if timeout, ok := d.GetOk("timeout"); ok {
		timeInt := timeout.(int)
		req.Timeout = &scw.Duration{Seconds: int64(timeInt)}
	}

	if port, ok := d.GetOk("port"); ok {
		req.Port = scw.Uint32Ptr(uint32(port.(int)))
	}

	if description, ok := d.GetOk("description"); ok {
		req.Description = types.ExpandStringPtr(description)
	}

	if registryImage, ok := d.GetOk("registry_image"); ok {
		req.RegistryImage = types.ExpandStringPtr(registryImage)
	}

	if maxConcurrency, ok := d.GetOk("max_concurrency"); ok {
		req.MaxConcurrency = scw.Uint32Ptr(uint32(maxConcurrency.(int))) //nolint:staticcheck
	}

	if sandbox, ok := d.GetOk("sandbox"); ok {
		req.Sandbox = container.ContainerSandbox(sandbox.(string))
	}

	if scalingOption, ok := d.GetOk("scaling_option"); ok {
		scalingOptionReq, err := expandScalingOptions(scalingOption)
		if err != nil {
			return nil, err
		}
		req.ScalingOption = scalingOptionReq
	}

	return req, nil
}

func expandScalingOptions(scalingOptionSchema interface{}) (*container.ContainerScalingOption, error) {
	scalingOption, ok := scalingOptionSchema.(*schema.Set)
	if !ok {
		return &container.ContainerScalingOption{}, nil
	}

	for _, option := range scalingOption.List() {
		rawOption, isRawOption := option.(map[string]interface{})
		if !isRawOption {
			continue
		}

		setFields := 0
		cso := &container.ContainerScalingOption{}
		if concurrentRequestThresold, ok := rawOption["concurrent_requests_threshold"].(int); ok && concurrentRequestThresold != 0 {
			cso.ConcurrentRequestsThreshold = scw.Uint32Ptr(uint32(concurrentRequestThresold))
			setFields++
		}
		if cpuUsageThreshold, ok := rawOption["cpu_usage_threshold"].(int); ok && cpuUsageThreshold != 0 {
			cso.CPUUsageThreshold = scw.Uint32Ptr(uint32(cpuUsageThreshold))
			setFields++
		}
		if memoryUsageThreshold, ok := rawOption["memory_usage_threshold"].(int); ok && memoryUsageThreshold != 0 {
			cso.MemoryUsageThreshold = scw.Uint32Ptr(uint32(memoryUsageThreshold))
			setFields++
		}

		if setFields > 1 {
			return &container.ContainerScalingOption{}, errors.New("a maximum of one scaling option can be set")
		}
		return cso, nil
	}

	return &container.ContainerScalingOption{}, nil
}

func flattenScalingOption(scalingOption *container.ContainerScalingOption) interface{} {
	if scalingOption == nil {
		return nil
	}

	flattenedScalingOption := []map[string]interface{}(nil)
	flattenedScalingOption = append(flattenedScalingOption, map[string]interface{}{
		"concurrent_requests_threshold": types.FlattenUint32Ptr(scalingOption.ConcurrentRequestsThreshold),
		"cpu_usage_threshold":           types.FlattenUint32Ptr(scalingOption.CPUUsageThreshold),
		"memory_usage_threshold":        types.FlattenUint32Ptr(scalingOption.MemoryUsageThreshold),
	})

	return flattenedScalingOption
}

func expandContainerSecrets(secretsRawMap interface{}) []*container.Secret {
	secretsMap := secretsRawMap.(map[string]interface{})
	secrets := make([]*container.Secret, 0, len(secretsMap))

	for k, v := range secretsMap {
		secrets = append(secrets, &container.Secret{
			Key:   k,
			Value: types.ExpandStringPtr(v),
		})
	}

	return secrets
}

func isContainerDNSResolveError(err error) bool {
	responseError := &scw.ResponseError{}

	if !errors.As(err, &responseError) {
		return false
	}

	if strings.HasPrefix(responseError.Message, "could not validate domain") {
		return true
	}

	return false
}

func retryCreateContainerDomain(ctx context.Context, containerAPI *container.API, req *container.CreateDomainRequest, timeout time.Duration) (*container.Domain, error) {
	timeoutChannel := time.After(timeout)

	for {
		select {
		case <-time.After(DefaultContainerRetryInterval):
			domain, err := containerAPI.CreateDomain(req, scw.WithContext(ctx))
			if err != nil && isContainerDNSResolveError(err) {
				continue
			}
			return domain, err
		case <-timeoutChannel:
			return containerAPI.CreateDomain(req, scw.WithContext(ctx))
		}
	}
}
