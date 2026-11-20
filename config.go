package main

type containerSpec struct {
	name     string   `yaml:"name"`
	image    string   `yaml:"image"`
	commands []string `yaml:"commands"`
	ports    []int    `yaml:"ports"`
}

type config struct {
	spec []*containerSpec `yaml:"containerSpec"`
}
