package main

import (
	"scm.eadn.dz/DevOps/ssh_tunneling/config"
)

func main() {

	config, err := config.NewConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	for key, value := range config.Apps {
		service := New(key, value)

		err = service.StartSSHTunnel()
		if err != nil {
			panic(err)
		}
	}

}
