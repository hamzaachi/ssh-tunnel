package main

import (
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	MinPort = 1000
	MaxPort = 9000
)

func CheckPortStatus(IP, Port string) bool {
	out, err := exec.Command("nc", "-nzv", IP, Port).CombinedOutput()

	if err != nil {
		return false
	}
	if strings.Contains(string(out), "succeeded!") {
		return true
	} else {
		return false
	}
}

func AddFirewallRule(Rule string) error {
	args := strings.Split(Rule, " ")
	cmd := exec.Command("ufw", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func GetRandomNumber() string {
	rand.Seed(time.Now().UnixNano())

	return strconv.Itoa(rand.Intn(MaxPort-MinPort+1) + MinPort)
}

func Display(output []Recod) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Name", "Type", "Port"})
	t.Style().Color.Header = text.Colors{text.Bold, text.FgBlue}

	for _, item := range output {
		t.AppendRow([]interface{}{item.ID, item.Name, item.Type, item.Port})
		t.AppendSeparator()
	}
	t.Render()
}
