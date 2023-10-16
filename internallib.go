package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

func removeSubstring(originalString, substringToRemove string) string {
	// Check if the substring exists in the original string
	if strings.Contains(originalString, substringToRemove) {
		// Replace the substring with an empty string to remove it
		return strings.Replace(originalString, substringToRemove, "", -1)
	}
	// If the substring is not found, return the original string
	return originalString
}

func startBinary(binaryPath string, args ...string) error {
	cmd := exec.Command(binaryPath, args...)

	// Set the output to the current process's output (your terminal)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	return err
}

func OpenWebsite(url string) {

	open.Run(url)
}
