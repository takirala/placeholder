package docker

import (
	"errors"

	"github.com/docker/docker/client"
)

var (
	ErrMissingParameter = errors.New("missing parameter")
)

type docker struct {
	*client.Client
}

type API interface {
	NetworkAPI
}

func NewAPI() (API, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &docker{dc}, nil
}
