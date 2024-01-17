package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	log "github.com/charmbracelet/log"
)

func Shutdown() error {

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("shutdown", "/s")
	case "linux":
		cmd = exec.Command("shutdown", "-h", "now")
	case "darwin":
		cmd = exec.Command("shutdown", "-h", "now")
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Run()
	if err != nil {
		log.Error("Error shutting down:", err)
		return err
	}

	log.Warn("System shutdown triggered")
	os.Exit(0)
	return nil
}

func Reboot() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("shutdown", "/r")
	case "linux":
		cmd = exec.Command("reboot")
	case "darwin":
		cmd = exec.Command("shutdown", "-r", "now")
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Run()
	if err != nil {
		log.Error("Error rebooting: ", err)
		return err
	}

	log.Warn("System reboot triggered")
	os.Exit(0)
	return nil

}

func Sleep() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
	case "linux":
		cmd = exec.Command("systemctl", "suspend")
	case "darwin":
		cmd = exec.Command("pmset", "sleepnow")
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Run()
	if err != nil {
		log.Error("Error sleeping:", err)
		return err
	}
	log.Warn("System suspend triggered")
	return nil
}
