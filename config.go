package main

type ports struct {
	ContainerPorts int `yaml:"containerPorts"`
}

// get container spec from the manifest
type containerSpec struct {
	Name     string   `yaml:"name"`
	Image    string   `yaml:"image"`
	Commands []string `yaml:"commands"`
	Ports    []*ports `yaml:"ports"`
}

type config struct {
	Spec []*containerSpec `yaml:"containerSpec"`
}
