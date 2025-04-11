package kind

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nutanix-cloud-native/nkp-nutanix-product-catalog/apptests/docker"
)

var ErrMisconfiguredNetwork = errors.New("misconfigured kind network")

// GetDockerNetworkName returns docker network name for kind cluster.
func GetDockerNetworkName() string {
	// default network name
	kindNetwork := "kind"

	// env var for override network name
	if network, ok := os.LookupEnv("KIND_EXPERIMENTAL_DOCKER_NETWORK"); ok {
		kindNetwork = network
	}

	return kindNetwork
}

// EnsureDockerNetworkExist ensures docker network exist with given configurations.
func EnsureDockerNetworkExist(ctx context.Context, subnet string, internal bool) (*docker.NetworkResource, error) {
	dapi, err := docker.NewAPI()
	if err != nil {
		return nil, fmt.Errorf("could not create docker api:%w", err)
	}

	name := GetDockerNetworkName()

	ok, networkResource, err := dapi.GetNetwork(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker network %s: %w", name, err)
	}

	if ok && networkResource.Internal != internal {
		return nil, fmt.Errorf("internal flag does not match for network %s: %w", name, ErrMisconfiguredNetwork)
	}

	if ok && subnet != "" {
		ipamConfigs := networkResource.IPAM.Config

		if len(ipamConfigs) == 0 {
			return nil, fmt.Errorf("subnet configuration is missing for network %s: %w", name, ErrMisconfiguredNetwork)
		}

		// we take only the first
		actual := ipamConfigs[0].Subnet

		if actual != subnet {
			return nil, fmt.Errorf("subnet expected %s actual %s for network %s: %w", actual, subnet, name, ErrMisconfiguredNetwork)
		}
	}

	if !ok {
		networkResource, err = dapi.CreateNetwork(context.Background(), name, internal, subnet)

		if err != nil {
			return nil, fmt.Errorf("failed to create network %s: %w", name, err)
		}
	}

	return networkResource, nil
}
