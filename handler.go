package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"

	"scm.eadn.dz/DevOps/ssh_tunneling/config"
)

const (
	VPNSubnet = "10.10.10.0/24"
	ServerIP  = "10.1.0.100"
)

type Tunnel struct {
	Name      string
	Category  string
	LocalPort string
	DistPort  string
	IP        string
	SSHServer string
	db        *sql.DB
}

func New(ctx context.Context, name string, app config.App, db *sql.DB) *Tunnel {
	s := Tunnel{}
	s.Name = name
	s.SSHServer = app.Shh
	s.db = db

	if len(app.Bdd) > 0 {
		s.Category = "bdd"
		s.DistPort = strings.Split(app.Bdd, ":")[1]
		s.IP = strings.Split(app.Bdd, ":")[0]
		s.LocalPort = GetRandomNumber()
	}
	if len(app.Web) > 0 {
		s.Category = "app"
		s.DistPort = strings.Split(app.Web, ":")[1]
		s.IP = strings.Split(app.Web, ":")[0]
		s.LocalPort = GetRandomNumber()
	}

	return &s
}

func (sshTunnel *Tunnel) GenerateSystemdService(ctx context.Context) error {
	Filename := "/etc/systemd/system/ssh-tunnel-" + sshTunnel.Name + "-" + sshTunnel.Category + ".service"
	f, err := os.Create(Filename)
	if err != nil {
		return err
	}

	temp := template.Must(template.ParseFiles("templates/ssh-tunnel.service"))
	err = temp.Execute(f, sshTunnel)
	if err != nil {
		return err
	}
	return nil
}

func (sshTunnel *Tunnel) GenerateHAProxyBackend(ctx context.Context) error {
	Filename := "/etc/haproxy/conf.d/" + sshTunnel.Name + ".cfg"
	f, err := os.Create(Filename)
	if err != nil {
		return err
	}

	temp := template.Must(template.ParseFiles("templates/haproxy-backend.cfg"))
	err = temp.Execute(f, sshTunnel)
	if err != nil {
		return err
	}
	return nil
}

func (sshTunnel *Tunnel) StartSSHTunnel(ctx context.Context) error {
	err := sshTunnel.GenerateSystemdService(ctx)
	if err != nil {
		return err
	}

	systemdService := "ssh-tunnel-" + sshTunnel.Name + "-" + sshTunnel.Category + ".service"
	cmd := exec.Command("systemctl", "start", systemdService)
	err = cmd.Run()
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	if CheckPortStatus("127.0.0.1", sshTunnel.LocalPort) {
		cmd = exec.Command("systemctl", "enable", systemdService)
		err = cmd.Run()
		if err != nil {
			return err
		}

		if sshTunnel.Category == "app" {
			err := sshTunnel.GenerateHAProxyBackend(ctx)
			if err != nil {
				return err
			}
		}

		FirewallRule := "allow from " + VPNSubnet + " to " + ServerIP + " proto tcp port" + " " + sshTunnel.LocalPort
		err = ExecFirewallRule(FirewallRule)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Something went wrong, cannot check port")
	}

	err = sshTunnel.Insert(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (sshTunnel *Tunnel) StopSSHTunnel(ctx context.Context, Name, Type string) error {

	tunnel, err := sshTunnel.RetrieveByID(ctx, Name, Type)
	if err != nil {
		return err
	}
	if len(tunnel) == 0 {
		log.Println("SSL Tunnel does not Exist!")
		return nil
	}

	FirewallRule := "delete allow from " + VPNSubnet + " to " + ServerIP + " proto tcp port" + " " + strconv.Itoa(tunnel[0].Port)
	err = ExecFirewallRule(FirewallRule)
	if err != nil {
		return err
	}

	//This Needs Improvements...
	systemdService := "ssh-tunnel-" + Name + "-" + Type + ".service"
	systemdServicePath := "/etc/systemd/system/" + systemdService
	cmd := exec.Command("systemctl", "start", systemdService)
	err = cmd.Run()
	cmd = exec.Command("systemctl", "disable", systemdService)
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("rm", "-f", systemdServicePath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = sshTunnel.Delete(ctx, Name, Type)
	if err != nil {
		return err
	}

	return nil
}
