package main

import (
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

type Service struct {
	Name      string
	Category  string
	LocalPort string
	DistPort  string
	IP        string
	SSHServer string
}

func New(name string, app config.App) *Service {
	s := Service{}
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

func (service *Service) CreateSystemdService() error {
	Filename := "/etc/systemd/system/ssh-tunnel-" + service.Name + "-" + service.Category + ".service"
	f, err := os.Create(Filename)
	if err != nil {
		return err
	}

	temp := template.Must(template.ParseFiles("templates/ssh-tunnel.service"))
	err = temp.Execute(f, service)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) StartSSHTunnel() error {
	err := service.CreateSystemdService()
	if err != nil {
		return err
	}

	systemdService := "ssh-tunnel-" + service.Name + "-" + service.Category + ".service"
	cmd := exec.Command("systemctl", "start", systemdService)
	err = cmd.Run()
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	if CheckPortStatus("127.0.0.1", service.LocalPort) {
		cmd = exec.Command("systemctl", "enable", systemdService)
		err = cmd.Run()
		if err != nil {
			return err
		}

		FirewallRule := "allow from " + VPNSubnet + " to " + ServerIP + " proto tcp port" + " " + service.LocalPort
		err = AddFirewallRule(FirewallRule)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Something went wrong, cannot check port")
	}
	return nil
}
