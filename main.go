package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ntp = []string{"ntpd", "-q", "-g"}

// Execute takes a series of strings as arguments, executes a command, and returns the string output and an error
func Execute(args ...string) (string, string, error) {
	cmdName := args[0]
	cmdArgs := args[:1]

	cmd := exec.Command(cmdName, cmdArgs...)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return "", "", nil
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return "", "", nil
	}

	outScanner := bufio.NewScanner(outReader)
	go func() {
		for outScanner.Scan() {
		}
	}()

	errScanner := bufio.NewScanner(errReader)
	go func() {
		for errScanner.Scan() {
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return "", "", err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return "", "", err
	}

	return outScanner.Text(), errScanner.Text(), err
}

// Bootstrap performs the initial Gentoo Stage3 install.
func Bootstrap() {
	fmt.Println("Please partition and mount your disks before continuing!")
	fmt.Println("Proceed with installation? [Y/N]: ")
	var proceed string
	fmt.Scan(&proceed)

	switch strings.ToLower(proceed) {
	case "yes":
		fallthrough
	case "y":
		fmt.Println("Yes selected")
	case "no":
		fallthrough
	case "n":
		fmt.Println("No selected")
	default:
		fmt.Println("Invalid choice.")
		os.Exit(1)
	}

	stdOut, stdErr, err := Execute("ntpd", "-q", "-q")
	fmt.Println(stdOut)
	fmt.Println(stdErr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// Minimal performs a basic installation of essentials: vim, ssh, tmux, etc.
func Minimal() {}

// Desktop installs a graphical environment.
func Desktop() {}

// Laptop installs wireless drivers and so on.
func Laptop() {}

// Server installs kubernetes, rkt/docker, etc.
func Server() {}

// Cleanup runs after installation and removes unnecessary source files.
func Cleanup() {}

func main() {
	fmt.Println("vim-go")
}
