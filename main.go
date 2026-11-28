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

type (
	upStatus   int
	downStatus int
)

type ports struct {
	ContainerPorts int `yaml:"containerPorts"`
}

// get container spec from the manifest
type containerSpec struct {
	up       upStatus
	down     downStatus
	Name     string   `yaml:"name"`
	Image    string   `yaml:"image"`
	Commands []string `yaml:"commands"`
	Ports    []*ports `yaml:"ports"`
}

type podSpec struct {
	Name          string           `yaml:"name"`
	ContainerSpec []*containerSpec `yaml:"containerSpec"`
}

func main() {
	args := os.Args
	filePath := args[1]

	output, err := readYAMLManifest(filePath)
	if err != 0 {
		os.Exit(err)
	}
	output.runPauseContainer()
	output.createContainer()

	if os.Args[1] == "list" {
		fmt.Fprint(os.Stdout, "Name\tRunning")
		// listPods()
	}
}

// func listPods() int {
// 	contregistered := len(podList)
// 	return contregistered
// }

func (pd *podSpec) runPauseContainer() {
	pd.createContainer()
}

func readYAMLManifest(manifest string) (*podSpec, int) {
	var pd podSpec
	data, err := os.ReadFile(manifest)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, 1
	}

	if err := yaml.Unmarshal(data, &pd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, 1
	}

	return &pd, 0
}

func (pd *podSpec) createContainer() {
	cnt := &containerSpec{}
	spec := pd.ContainerSpec
	if len(spec) == 0 {
		os.Exit(1)
	}
	for i := range spec {
		cnt.Name = spec[i].Name
		cnt.Image = spec[i].Image
	}

	fmt.Println("Creating container now...")
	client, err := containerd.New(containerdSocket)

	if err != nil {
		panic(err)
	}

	ctx := namespaces.WithNamespace(context.Background(), namespace)
	fmt.Printf("Pulling %s....", cnt.Image)
	image, err := client.Pull(ctx, fmt.Sprintf("docker.io/library/%s", cnt.Image), containerd.WithPullUnpack)
	if err != nil {
		panic(err)
	}
	_, err = client.NewContainer(ctx, cnt.Name,
		containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", cnt.Image), image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println("Container created")
}
