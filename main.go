package main

var containerdSocket = "/run/containerd/containerd.sock"

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
	var p Pod
	output := readYAMLManifest()
	p.createContainer(output)
}

func readYAMLManifest() string { return "" }

func (p *Pod) createContainer(input string) *container { return &container{} }
