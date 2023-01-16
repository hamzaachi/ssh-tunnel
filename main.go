package main

import (
	"context"

	"scm.eadn.dz/DevOps/ssh_tunneling/config"
)

func main() {

	config, err := config.NewConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	db, err := connect(ctx)
	if err != nil {
		panic(err)
	}

	for key, value := range config.Apps {
		service := New(key, value)
		service.db = db
		//err = service.StartSSHTunnel()
		//if err != nil {
		//	panic(err)
		//}

		if service.Name == "invest" {
			err := service.Insert()
			if err != nil {
				panic(err)
			}
		}
	}

}
