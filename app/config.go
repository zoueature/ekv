package main

type Conf struct {
	Host    string   `yaml:"host"`
	Cluster []string `yaml:"cluster""`
}
