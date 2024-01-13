package main

import (
	"os"
	"os/exec"

	log "github.com/charmbracelet/log"
)

//Functions

func shutdown() error {
	cmd := exec.Command("shutdown", "-h", "now")
	err := cmd.Run()
	if err != nil {
		log.Error("Error shutting down:", err)
		return err
	}

	log.Error("System shutdown triggered")
	os.Exit(0)
	return nil
}

func reboot() error {
	// Execute the reboot command
	cmd := exec.Command("reboot")
	err := cmd.Run()
	if err != nil {
		log.Error("Error rebooting: ", err)
		return err
	}

	log.Info("System reboot triggered")
	os.Exit(0)
	return nil

}

func sleep() error {
	cmd := exec.Command("systemctl", "suspend")
	err := cmd.Run()
	if err != nil {
		log.Error("Error sleeping:", err)
		return err
	}
	log.Error("System suspend triggered")
	return nil
}
