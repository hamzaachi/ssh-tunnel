package main

import (
	"context"
	"fmt"

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

	defer db.Close()

	var service *Service
	for key, value := range config.Apps {
		service = New(key, value)
		service.db = db
		//err = service.StartSSHTunnel()
		//if err != nil {
		//	panic(err)
		//}

	}

	List, err := service.RetrieveByName("invpest")
	if err != nil {
		panic(err)
	}

	fmt.Println(List)
}
