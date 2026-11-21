package main

// get container spec from the manifest
type containerSpec struct {
	Name     string   `yaml:"name"`
	Image    string   `yaml:"image"`
	Commands []string `yaml:"commands"`
	Ports    []int    `yaml:"ports"`
}

type config struct {
	Spec []*containerSpec `yaml:"containerSpec"`
}
