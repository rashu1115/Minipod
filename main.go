package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// var containerdSocket = "/run/containerd/containerd.sock"

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
	containersList []*container
}

func main() {
	args := os.Args
	filePath := args[1]
	output := readYAMLManifest(filePath)
	cnt := createContainer(output)
	fmt.Println(cnt)
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

func createContainer(input *config) *container {
	cnt := &container{}
	spec := input.Spec
	if len(spec) == 0 {
		os.Exit(1)
	}
	for i := range spec {
		cnt.name = spec[i].Name
		cnt.image = spec[i].Image
	}
	return cnt
}
