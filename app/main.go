package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func main() {
	config := conf()
	fmt.Printf("%+v", config)
}

func conf() *Conf {
	conf := new(Conf)
	f, err := os.Open("./conf/conf.yaml")
	if err != nil {
		log.Printf("open config file error: %s", err.Error())
		os.Exit(-1)
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(conf)
	if err != nil {
		log.Printf("config syntax error: %s", err.Error())
		os.Exit(-1)
	}
	return conf
}