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
	const S3Arch = "amd64"
	const S3Date = "20170727"

	_exec := ExecuteWithEnv([]string{})
	_execP := ExecutePWithEnv([]string{})

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

	// TODO Is a correct time sync strictly necessary for install to succeed?
	_, _, err := _exec("ntpd", "-q", "-q")
	if err != nil {
		fmt.Println(err)
		fmt.Println("WARNING: ntpd failed to synchronize.")
	}

	archEnv := []string{"S3_ARCH=amd64", "S3_DATE=20170727"}

	wgetURL := fmt.Sprintf("http://distfiles.gentoo.org/releases/%s/autobuilds/%s/stage3-%s-nomultilib-%s.tar.bz2", S3Arch, S3Date)

	wgetCmd := []string{"wget", wgetURL}

	ExecuteP(archEnv, "WARNING: Install archive download failed. Aborting.", wgetCmd...)

	tarCmd := []string{"tar", "xvjpf", "stage3-*.tar.bz2", "--xattrs", "--numeric-owner"}

	// For now, I'm using ExecuteP here. During the second pass, this should be replaced with Execute(); we can inspect the error and perhaps re-download, re-extract, etc.
	_execP("WARNING: Archive extraction failed. Aborting.", tarCmd...)

	procMountCmd := []string{"mount", "-t", "proc", "/proc", "/mnt/gentoo/proc"}
	_execP("WARNING: Mounting /proc to chroot failed. Aborting.", procMountCmd...)

	procRbindCmd := []string{"mount", "--rbind", "/sys", "/mnt/gentoo/sys"}
	_execP("WARNING: Mounting /proc to chroot failed. Aborting.", procRbindCmd...)

	devMountCmd := []string{"mount", "--rbind", "/sys", "/mnt/gentoo/sys"}
	_execP("WARNING: Rslave sys mounting failed.", devMountCmd...)

	devRbindCmd := []string{"mount", "--make-rslave", "/mnt/gentoo/dev"}
	_execP("WARNING: Mounting /proc to chroot failed. Aborting.", devRbindCmd...)

	swapFile := "/mnt/gentoo/swapfile"

	swapCmds := [][]string{
		[]string{"mkswap", swapFile},
		[]string{"chmod", "0600", swapFile},
		[]string{"swapon", swapFile},
	}

	if _, err := os.Stat(swapFile); os.IsNotExist(err) {
		makeSwapCmd := []string{"fallocate", "-l", "2G", swapFile}
		_execP("Creating swapfile failed.", makeSwapCmd...)
	}

	for _, cmd := range swapCmds {
		_execP("Activating swap failed.", cmd...)
	}

	fetchMakeConfig := []string{"wget", "-O", "/mnt/gentoo/etc/portage/make.conf", "https://raw.githubusercontent.com/jcmdln/gein/master/etc/portage/make.conf"}
	_execP("Failed to download make config file.", fetchMakeConfig...)

	fetchUseConfig := []string{"wget", "-O", "/mnt/gentoo/etc/portage/package.use", "https://raw.githubusercontent.com/jcmdln/gein/master/etc/portage/package.use"}
	_execP("Failed to download package.use config file.", fetchUseConfig...)

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
