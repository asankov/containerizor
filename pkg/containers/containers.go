package containers

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// Orchestrator is the struct that
// holds the logic of operating
// the containers
type Orchestrator struct {
	dockerClient client.CommonAPIClient
}

type Container struct {
	ID    string
	Image string
}

// NewOrchestrator creates a new Orchestrator
// and returns a reference to that
func NewOrchestrator(cl client.CommonAPIClient) *Orchestrator {
	return &Orchestrator{
		dockerClient: cl,
	}
}

// Start creates a new container with the image whose name is passed
// starts the container and returns its id or an error if such occurred.
// The format of imageName is <account>/<image>:<label(optional)>.
func (o *Orchestrator) Start(imageName string) (string, error) {
	// TODO: imageName validation

	ctx := context.Background()

	// TODO: see if image is present before pulling
	if _, err := o.dockerClient.ImagePull(ctx, "docker.io/"+imageName, types.ImagePullOptions{}); err != nil {
		return "", fmt.Errorf("error pulling image `%s` from docker.io: %w", imageName, err)
	}

	container, err := o.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, &container.HostConfig{}, &network.NetworkingConfig{}, "")
	if err != nil {
		return "", fmt.Errorf("error creating container: %w", err)
	}

	if err := o.dockerClient.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("error starting container: %w", err)
	}

	return container.ID, nil
}

func (o *Orchestrator) ListContainers() ([]*Container, error) {
	containersResponse, err := o.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	containers := make([]*Container, 0, len(containersResponse))
	for _, c := range containersResponse {
		containers = append(containers, &Container{
			ID:    c.ID,
			Image: c.Image,
		})
	}

	return containers, nil
}
