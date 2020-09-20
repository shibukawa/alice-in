// +build windows

package alicein

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shirou/gopsutil/process"
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

// WSLGuestEnvironment returns environment to access Guest(Linux) environment
// If WSL is not installed, it returns host environment
func WSLGuestEnvironment() Environment {
	if IsWSLInstalled() {
		return &wslGuestEnvironment{}
	}
	return &nonVirtualEnvironment{}
}

// WSLHostEnvironment returns environment to access Host(Windows) environment
func WSLHostEnvironment() Environment {
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
	envs := e.Environ()
	if v, ok := envs["HOME"]; ok {
		return v, nil
	}
	return "", errors.New("$HOME is not defined")
}

func (e wslGuestEnvironment) UserConfigDir() (string, error) {
	envs := e.Environ()
	dir, ok := envs["XDG_CONFIG_HOME"]
	if ok {
		return dir, nil
	}
	dir, ok = envs["HOME"]
	if !ok {
		return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
	}
	return dir + "/.config", nil
}

func (e wslGuestEnvironment) UserCacheDir() (string, error) {
	envs := e.Environ()
	dir, ok := envs["XDG_CACHE_HOME"]
	if ok {
		return dir, nil
	}
	dir, ok = envs["HOME"]
	if !ok {
		return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
	}
	return dir + "/.cache", nil
}

var cachedGuestEnv map[string]string
var checkEnvOnce sync.Once

func (e wslGuestEnvironment) Environ() map[string]string {
	checkEnvOnce.Do(func() {
		// wslvar is better because cmd.exe /C set specify
		// static location of windows. but wsvar reset terminal
		// unexpectedly
		//cmd := exec.Command("wslvar", "--getsys")
		cmd := exec.Command("wsl", "env")
		envs, err := cmd.Output()
		if err != nil {
			return
		}
		envStrs := strings.Split(string(envs), "\n")
		cachedGuestEnv = make(map[string]string, len(envStrs))
		for _, l := range envStrs {
			f := strings.SplitN(l, "=", 2)
			if len(f) == 2 {
				cachedGuestEnv[f[0]] = strings.TrimSpace(f[1])
			}
		}
	})
	return cachedGuestEnv
}

func (e wslGuestEnvironment) Type() EnvType {
	return WSLEnv
}

// DetectEnvType returns detected parent environment type
//
// This function check's parent process and return result
func DetectEnvType() (EnvType, error) {
	p, err := process.NewProcess(int32(os.Getppid()))
	if err != nil {
		return 0, err
	}
	c, err := p.CmdlineSlice()
	if err != nil {
		return 0, err
	}
	if filepath.Base(c[0]) == "wsl.exe" {
		return WSLEnv, nil
	}
	return NativeEnv, nil
}

// DetectedEnvironment returns environment of parent
//
// On Windows, it returns native or WSL environment based on parent's process
func DetectedEnvironment() Environment {
	env, _ := DetectEnvType()
	if env == WSLEnv {
		return WSLGuestEnvironment()
	}
	return NativeEnvironment()
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
