package main

import (
	"os"

	"go.yaml.in/yaml/v3"
)

// var containerdSocket = "/run/containerd/containerd.sock"

type namespaceList []string

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
	createContainer(output)
}

func readYAMLManifest(manifest string) *config {
	var cfg *config
	data, err := os.ReadFile(manifest)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	return cfg
}

func createContainer(input *config) *container {
	cnt := &container{}
	spec := input.spec
	cnt.name = spec[0].name
	return cnt
}
