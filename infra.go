package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var instance *client.Client
var networkName = "metricbeat-devnet"

func getDockerClient() *client.Client {
	if instance != nil {
		return instance
	}

	instance, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	return instance
}

func getDevNetwork() (types.NetworkResource, error) {
	dockerClient := getDockerClient()

	ctx := context.Background()

	networkResource, err := dockerClient.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{
		Verbose: true,
	})
	if err != nil {
		fmt.Println("Not found! Creating " + networkName)
		initDevNetwork()
	} else {
		fmt.Println("Reusing existing " + networkName)
	}

	return networkResource, err
}

func removeDevNetwork() error {
	dockerClient := getDockerClient()

	ctx := context.Background()

	fmt.Println("Removing existing " + networkName)
	if err := dockerClient.NetworkRemove(ctx, networkName); err != nil {
		return err
	}

	fmt.Println("Network " + networkName + " removed!")

	return nil
}

func initDevNetwork() types.NetworkCreateResponse {
	dockerClient := getDockerClient()

	ctx := context.Background()

	// Construct network create request body
	nc := types.NetworkCreate{
		Driver:         "bridge",
		CheckDuplicate: true,
		Internal:       true,
		EnableIPv6:     false,
		Attachable:     true,
		Labels: map[string]string{
			"project": "metricbeat",
			"runtime": "test",
		},
	}

	response, err := dockerClient.NetworkCreate(ctx, networkName, nc)
	if err != nil {
		panic("Cannot create Docker network. Aborting")
	}

	fmt.Println("network created: " + response.ID)

	return response
}
