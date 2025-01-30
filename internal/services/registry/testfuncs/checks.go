package registrytestfuncs

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	dockerRegistrySDK "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry"
)

var testDockerIMG = "docker.io/library/alpine:latest"

func PushImageToRegistry(tt *acctest.TestTools, registryEndpoint string, tagName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if !*acctest.UpdateCassettes {
			return nil
		}

		meta := tt.Meta
		var errorMessage registry.ErrorRegistryMessage

		accessKey, _ := meta.ScwClient().GetAccessKey()
		secretKey, _ := meta.ScwClient().GetSecretKey()

		authConfig := dockerRegistrySDK.AuthConfig{
			ServerAddress: registryEndpoint,
			Username:      accessKey,
			Password:      secretKey,
		}

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return fmt.Errorf("could not connect to Docker: %v", err)
		}

		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return fmt.Errorf("could not marshal auth config: %v", err)
		}

		ctx := context.Background()
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		out, err := cli.ImagePull(ctx, testDockerIMG, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("could not pull image: %v", err)
		}
		defer out.Close()

		buffIOReader := bufio.NewReader(out)
		for {
			streamBytes, errPull := buffIOReader.ReadBytes('\n')
			if errPull == io.EOF {
				break
			}
			err = json.Unmarshal(streamBytes, &errorMessage)
			if err != nil {
				return fmt.Errorf("could not unmarshal: %v", err)
			}

			if errorMessage.Error != "" {
				return errors.New(errorMessage.Error)
			}
		}

		scwTag := registryEndpoint + "/alpine:" + tagName
		err = cli.ImageTag(ctx, testDockerIMG, scwTag)
		if err != nil {
			return fmt.Errorf("could not tag image: %v", err)
		}

		pusher, err := cli.ImagePush(ctx, scwTag, image.PushOptions{RegistryAuth: authStr})
		if err != nil {
			return fmt.Errorf("could not push image: %v", err)
		}
		defer pusher.Close()

		buffIOReader = bufio.NewReader(pusher)
		for {
			streamBytes, errPush := buffIOReader.ReadBytes('\n')
			if errPush == io.EOF {
				break
			}
			err = json.Unmarshal(streamBytes, &errorMessage)
			if err != nil {
				return fmt.Errorf("could not unmarshal: %v", err)
			}

			if errorMessage.Error != "" {
				return errors.New(errorMessage.Error)
			}
		}

		return nil
	}
}
