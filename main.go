package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"

	"scm.eadn.dz/DevOps/ssh_tunneling/config"
)

func main() {

	var list, create bool
	var cmd, name, Type string

	flag.BoolVar(&list, "list", false, "List All Created SSH Tunnels")
	flag.BoolVar(&create, "create", false, "Create a new SSH Tunnel")
	deleteCMD := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCMD.StringVar(&name, "name", "", "Name of the SSH channel to be deleted")
	deleteCMD.StringVar(&Type, "type", "", "type of the SSH channel to be deleted")
	flag.Parse()
	args := flag.Args()
	if len(args) != 0 {
		cmd = args[0]

	}

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

	sshTunnel = &Tunnel{Name: "", Category: "", LocalPort: "", DistPort: "", IP: "", SSHServer: "", db: db}

	switch {
	case list:
		List, err := sshTunnel.List(ctx)
		if err != nil {
			panic(err)
		}
		Display(List)
		os.Exit(0)
	case create:
		Add(ctx, *config, db)
		os.Exit(0)
	}
	switch cmd {
	case "delete":
		deleteCMD.Parse(os.Args[2:])
		if name == "" || Type == "" {
			flag.Usage()
			deleteCMD.Usage()
			os.Exit(1)
		}
		err := sshTunnel.StopSSHTunnel(ctx, name, Type)
		if err != nil {
			panic(err)
		}
	default:
		flag.Usage()
		deleteCMD.Usage()
		os.Exit(1)
	}
}

func Add(ctx context.Context, conf config.Applications, db *sql.DB) {
	for key, value := range conf.Apps {
		sshTunnel := New(ctx, key, value, db)

		tunnel, err := sshTunnel.RetrieveByID(ctx, sshTunnel.Name, sshTunnel.Category)
		if err != nil {
			panic(err)
		}

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
