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

	defer db.Close()

	var sshTunnel *Tunnel
	for key, value := range config.Apps {
		sshTunnel = New(key, value, db)

		err := sshTunnel.GenerateHAProxyBackend()
		if err != nil {
			panic(err)
		}
		//err := sshTunnel.Insert()
		//err = service.StartSSHTunnel()
		//if err != nil {
		//	panic(err)
		//}

	}

	List, err := sshTunnel.List()
	if err != nil {
		panic(err)
	}

	Display(List)
}
