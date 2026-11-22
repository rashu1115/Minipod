package main

import (
	"context"
	"fmt"
	"os"

	namespace "github.com/containerd/containerd/namespaces"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/oci"
	"go.yaml.in/yaml/v3"
)

var containerdSocket = "/run/containerd/containerd.sock"

type namespaceList []string

// container native info to pass
// on to containerd
type container struct {
	name      string
	image     string
	rootFile  []byte
	namespace []*namespaceList
}

type Pod struct {
	name           string
	containersList []*container
}

func main() {
	args := os.Args
	filePath := args[1]
	output := readYAMLManifest(filePath)
	createContainer(output)
}

func readYAMLManifest(manifest string) *config {
	var cfg config
	data, err := os.ReadFile(manifest)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func createContainer(input *config) {
	cnt := &container{}
	spec := input.Spec
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

	ctx := namespace.WithNamespace(context.Background(), "minipod")
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
