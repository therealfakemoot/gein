package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Execute takes a series of strings as arguments, executes a command, and returns the string output and an error
func Execute(env []string, args ...string) (string, string, error) {
	cmdName := args[0]
	cmdArgs := args[:1]

	cmd := exec.Command(cmdName, cmdArgs...)

	newEnv := append(os.Environ(), env...)
	cmd.Env = newEnv

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

// ExecuteP attempts to spawn an external process and panics with exit code 0.
func ExecuteP(env []string, msg string, args ...string) {
	var err error

	_, _, err = Execute(env, args...)
	if err != nil {
		fmt.Println(err)
		if len(msg) > 0 {
			fmt.Println(msg)
		}
		os.Exit(1)
	}

}

// ExecuteWithEnv returns a closure. This is a convenience function for setting up repeated calls using no or identical environments.
func ExecuteWithEnv(e []string) func(args ...string) (string, string, error) {
	return func(args ...string) (string, string, error) {
		return Execute(e, args...)
	}
}

// ExecutePWithEnv behaves as ExecuteWithEnv for the ExecuteP variant.
func ExecutePWithEnv(e []string) func(msg string, args ...string) {
	return func(msg string, args ...string) {
		Execute(e, args...)
	}

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
