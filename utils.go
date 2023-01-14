package main

import (
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MinPort = 1000
const MaxPort = 5000

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
