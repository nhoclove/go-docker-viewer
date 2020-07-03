package main

import (
	"context"
	"fmt"
	"go-docker-viewer/pkg/docker"
	"go-docker-viewer/pkg/gui"
)

func main() {
	gui.ShowMenu()
}

func init() {
	containers, err := docker.ListVolumes(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", containers)
}
