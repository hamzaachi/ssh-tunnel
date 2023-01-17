package main

import (
	"database/sql"
	"errors"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"scm.eadn.dz/DevOps/ssh_tunneling/config"
)

const VPNSubnet = "10.10.10.0/24"
const ServerIP = "10.1.0.100"

type Tunnel struct {
	Name      string
	Category  string
	LocalPort string
	DistPort  string
	IP        string
	SSHServer string
	db        *sql.DB
}

func New(name string, app config.App) *Tunnel {
	s := Tunnel{}
	s.Name = name
	s.SSHServer = app.Shh

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

func (sshTunnel *Tunnel) CreateSystemdService() error {
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

func (sshTunnel *Tunnel) StartSSHTunnel() error {
	err := sshTunnel.CreateSystemdService()
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

		FirewallRule := "allow from " + VPNSubnet + " to " + ServerIP + " proto tcp port" + " " + sshTunnel.LocalPort
		err = AddFirewallRule(FirewallRule)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Something went wrong, cannot check port")
	}
	return nil
}
