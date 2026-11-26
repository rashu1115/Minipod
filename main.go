package main

import (
	"context"
	"fmt"
	"os"

	namespaces "github.com/containerd/containerd/namespaces"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/oci"
	"go.yaml.in/yaml/v3"
)

var namespace string = "minipod.io"
var containerdSocket string = "/run/containerd/containerd.sock"

type namespaceList []string

// container native info to pass
// on to containerd
type container struct {
	name      string
	image     string
	rootFile  []byte
	namespace []*namespaceList
}

type pod struct {
	name           string
	containersList []*container
}

func main() {
	args := os.Args
	filePath := args[1]
	output := readYAMLManifest(filePath)
	createContainer(output)
}

func readYAMLManifest(manifest string) *podSpec {
	var pd podSpec
	data, err := os.ReadFile(manifest)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &pd); err != nil {
		panic(err)
	}

	return &pd
}

func createContainer(input *podSpec) {
	cnt := &container{}
	spec := input.ContainerSpec
	if len(spec) == 0 {
		os.Exit(1)
	}
	for i := range spec {
		cnt.name = spec[i].Name
		cnt.image = spec[i].Image
	}

	fmt.Println("Creating container now...")
	client, err := containerd.New(containerdSocket)

	if err != nil {
		panic(err)
	}

	ctx := namespaces.WithNamespace(context.Background(), namespace)
	fmt.Printf("Pulling %s....", cnt.image)
	image, err := client.Pull(ctx, fmt.Sprintf("docker.io/library/%s", cnt.image), containerd.WithPullUnpack)
	if err != nil {
		panic(err)
	}
	_, err = client.NewContainer(ctx, cnt.name,
		containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", cnt.image), image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println("Container created")
}
