package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	cli *client.Client
)

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

type (
	Container struct {
		ID      string
		ShortID string
		Names   []string
		Status  string
		State   string
	}

	Image struct {
		ID       string
		RepoTags []string
		Size     int64
	}

	Volume struct {
		Driver string
		Name   string
	}
)

func ListContainer(ctx context.Context) ([]*Container, error) {
	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	containers := make([]*Container, 0)
	for _, v := range containerList {
		containers = append(containers, &Container{
			ID:      v.ID,
			ShortID: v.ID[:10],
			Names:   v.Names,
			Status:  v.Status,
			State:   v.State,
		})
	}
	return containers, nil
}

func ListImages(ctx context.Context) ([]*Image, error) {
	imageSummary, err := cli.ImageList(ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	images := make([]*Image, 0)
	for _, v := range imageSummary {
		images = append(images, &Image{
			ID:       v.ID,
			RepoTags: v.RepoTags,
			Size:     v.Size,
		})
	}
	return images, nil
}

func ListVolumes(ctx context.Context) ([]*Volume, error) {
	volumeListOkBody, err := cli.VolumeList(ctx, filters.Args{})
	if err != nil {
		return nil, err
	}

	volumes := make([]*Volume, 0)
	for _, v := range volumeListOkBody.Volumes {
		volumes = append(volumes, &Volume{
			Driver: v.Driver,
			Name:   v.Name,
		})
	}
	return volumes, nil
}
