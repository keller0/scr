package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDocker(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		t.Log(container.ID[:10], container.Image)
	}
	assert.Equal(t, 200, 200)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	id, err := CreateContainer("yximages/gcc:8.3")
	if err != nil {
		panic(err)
	}
	t.Log("created container:", id)
	defer func() {
		err = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
		if err != nil {
			panic(err)
		}
		t.Log("container", id, "removed")
	}()

	status, err := cli.ContainerStats(ctx, id, false)
	if err != nil {
		panic(err)
	}

	t.Log(status)
	assert.Equal(t, 200, 200)
}

func TestStartManagers(t *testing.T) {
	StartManagers()
	time.Sleep(13 * time.Second)
	t.Log(len(GccWorker))
	t.Log(len(GoWorker))

	assert.Equal(t, true, len(GccWorker) < 13)

}
