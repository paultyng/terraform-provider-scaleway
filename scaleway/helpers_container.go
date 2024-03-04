package scaleway

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	defaultContainerNamespaceTimeout = 5 * time.Minute
	defaultContainerCronTimeout      = 5 * time.Minute
	defaultContainerTimeout          = 12*time.Minute + 30*time.Second
	defaultContainerDomainTimeout    = 10 * time.Minute
	defaultContainerRetryInterval    = 5 * time.Second
)

// containerAPIWithRegion returns a new container API and the region.
func containerAPIWithRegion(d *schema.ResourceData, m interface{}) (*container.API, scw.Region, error) {
	meta := m.(*Meta)
	api := container.NewAPI(meta.scwClient)

	region, err := extractRegion(d, meta)
	if err != nil {
		return nil, "", err
	}
	return api, region, nil
}

// containerAPIWithRegionAndID returns a new container API, region and ID.
func containerAPIWithRegionAndID(m interface{}, id string) (*container.API, scw.Region, string, error) {
	meta := m.(*Meta)
	api := container.NewAPI(meta.scwClient)

	region, id, err := parseRegionalID(id)
	if err != nil {
		return nil, "", "", err
	}
	return api, region, id, nil
}

func setCreateContainerRequest(d *schema.ResourceData, region scw.Region) (*container.CreateContainerRequest, error) {
	// required
	nameRaw := d.Get("name")
	namespaceID := d.Get("namespace_id")

	name := expandOrGenerateString(nameRaw.(string), "co")
	privacyType := d.Get("privacy")
	protocol := d.Get("protocol")
	httpOption := d.Get("http_option")

	req := &container.CreateContainerRequest{
		Region:      region,
		NamespaceID: expandID(namespaceID),
		Name:        name,
		Privacy:     container.ContainerPrivacy(privacyType.(string)),
		Protocol:    container.ContainerProtocol(*expandStringPtr(protocol)),
		HTTPOption:  container.ContainerHTTPOption(httpOption.(string)),
	}

	// optional
	if envVariablesRaw, ok := d.GetOk("environment_variables"); ok {
		req.EnvironmentVariables = expandMapPtrStringString(envVariablesRaw)
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
		req.Description = expandStringPtr(description)
	}

	if registryImage, ok := d.GetOk("registry_image"); ok {
		req.RegistryImage = expandStringPtr(registryImage)
	}

	if maxConcurrency, ok := d.GetOk("max_concurrency"); ok {
		req.MaxConcurrency = scw.Uint32Ptr(uint32(maxConcurrency.(int)))
	}

	return req, nil
}

func waitForContainerNamespace(ctx context.Context, containerAPI *container.API, region scw.Region, namespaceID string, timeout time.Duration) (*container.Namespace, error) {
	retryInterval := defaultContainerRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	ns, err := containerAPI.WaitForNamespace(&container.WaitForNamespaceRequest{
		Region:        region,
		NamespaceID:   namespaceID,
		RetryInterval: &retryInterval,
		Timeout:       scw.TimeDurationPtr(timeout),
	}, scw.WithContext(ctx))

	return ns, err
}

func waitForContainerCron(ctx context.Context, api *container.API, cronID string, region scw.Region, timeout time.Duration) (*container.Cron, error) {
	retryInterval := defaultContainerRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	request := container.WaitForCronRequest{
		CronID:        cronID,
		Region:        region,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}

	return api.WaitForCron(&request, scw.WithContext(ctx))
}

func waitForContainer(ctx context.Context, api *container.API, containerID string, region scw.Region, timeout time.Duration) (*container.Container, error) {
	retryInterval := defaultContainerRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	request := container.WaitForContainerRequest{
		ContainerID:   containerID,
		Region:        region,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}

	return api.WaitForContainer(&request, scw.WithContext(ctx))
}

func waitForContainerDomain(ctx context.Context, api *container.API, domainID string, region scw.Region, timeout time.Duration) (*container.Domain, error) {
	retryInterval := defaultContainerRetryInterval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	request := container.WaitForDomainRequest{
		DomainID:      domainID,
		Region:        region,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}

	return api.WaitForDomain(&request, scw.WithContext(ctx))
}

func expandContainerSecrets(secretsRawMap interface{}) []*container.Secret {
	secretsMap := secretsRawMap.(map[string]interface{})
	secrets := make([]*container.Secret, 0, len(secretsMap))

	for k, v := range secretsMap {
		secrets = append(secrets, &container.Secret{
			Key:   k,
			Value: expandStringPtr(v),
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
		case <-time.After(defaultContainerRetryInterval):
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

func flattenContainerSecretEnvironmentVariables(secrets []*container.SecretHashedValue) map[string]any {
	m := make(map[string]any, len(secrets))
	for _, s := range secrets {
		m[s.Key] = s.HashedValue
	}

	return m
}
