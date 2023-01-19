package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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

	Add(ctx, *config, db)

	sshTunnel = &Tunnel{Name: "", Category: "", LocalPort: "", DistPort: "", IP: "", SSHServer: "", db: db}
	List, err := sshTunnel.List(ctx)
	if err != nil {
		panic(err)
	}

	Display(List)
}

func Add(ctx context.Context, conf config.Applications, db *sql.DB) {
	for key, value := range conf.Apps {
		sshTunnel := New(ctx, key, value, db)

		tunnel, err := sshTunnel.RetrieveByID(ctx, sshTunnel.Name, sshTunnel.Category)
		if err != nil {
			panic(err)
		}
		fmt.Println(tunnel)
		if len(tunnel) > 0 {
			log.Println("SSH Tunnel Already Exist!, Skipping...")
			continue
		}
		err = sshTunnel.StartSSHTunnel(ctx)
		if err != nil {
			panic(err)
		}
		log.Println("SSH Tunnel: ", sshTunnel.Name, "Type: ", sshTunnel.Category, ",Added")

	}
}
