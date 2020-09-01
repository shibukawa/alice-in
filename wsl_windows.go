// +build windows

package alicein

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// IsInWSL returns current environment is in WSL or not
func IsInWSL() bool {
	d, _ := os.Getwd()
	return strings.HasPrefix(d, `\\wsl$`)
}

var isInWSLInstalled bool
var checkEnv sync.Once

// IsWSLInstalled returns current environment has WSL guest/host
func IsWSLInstalled() bool {
	checkEnv.Do(func() {
		_, err := exec.LookPath("wslconfig")
		isInWSLInstalled = err == nil
	})
	return isInWSLInstalled
}

// WSLGuest returns environment to access Guest(Linux) environment
// If WSL is not installed, it returns host environment
func WSLGuest() Environment {
	if IsWSLInstalled() {
		return &wslGuestEnvironment{}
	}
	return &nonVirtualEnvironment{}
}

// WSLHost returns environment to access Host(Windows) environment
func WSLHost() Environment {
	return &nonVirtualEnvironment{}
}

type wslGuestEnvironment struct {
}

var _ Environment = &nonVirtualEnvironment{}

func (e wslGuestEnvironment) Exec(ctx context.Context, cmd string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "wsl", append([]string{cmd}, args...)...)
}

func (e wslGuestEnvironment) Open(input string) {
	convertedPath := ConvertToGuestPath(input)
	cmd := exec.Command("wsl", "xdg-open", convertedPath)
	cmd.Run()
}

func (e wslGuestEnvironment) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (e wslGuestEnvironment) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (e wslGuestEnvironment) UserCacheDir() (string, error) {
	return os.UserCacheDir()
}

var cachedGuestEnv map[string]string
var checkEnvOnce sync.Once

func (e wslGuestEnvironment) Environ() map[string]string {
	checkEnvOnce.Do(func() {
		// wslvar is better because cmd.exe /C set specify
		// static location of windows. but wsvar reset terminal
		// unexpectedly
		//cmd := exec.Command("wslvar", "--getsys")
		cmd := exec.Command("wsl", "/mnt/c/Windows/System32/cmd.exe")
		envs, err := cmd.Output()
		if err != nil {
			return
		}
		envStrs := strings.Split(string(envs), "\n")
		cachedGuestEnv = make(map[string]string, len(envStrs))
		for _, l := range envStrs {
			f := strings.SplitN(l, "=", 2)
			cachedGuestEnv[f[0]] = f[1]
		}
	})
	return cachedGuestEnv
}

// ConvertToHostPath returns host style path if current env is in WSL.
// Otherwise, return input path as is.
// Docker doesn't support this.
func ConvertToHostPath(path string) string {
	cmd := exec.Command("wsl", "wslpath", "-w", path)
	result, err := cmd.Output()
	if err != nil {
		return path
	}
	return strings.TrimSpace(string(result))
}

// ConvertToGuestPath returns guest style path.
// Otherwise, return input path as is.
// Docker doesn't support this.
func ConvertToGuestPath(path string) string {
	cmd := exec.Command("wsl", "wslpath", path)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return path
	}
	return strings.TrimSpace(string(result))
}
