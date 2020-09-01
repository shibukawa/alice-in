// +build linux

package alicein

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var checkWSL sync.Once
var isInWSL bool

// IsInWSL returns current environment is in WSL or not
func IsInWSL() bool {
	checkWSL.Do(func() {
		_, err := os.Stat("/run/WSL")
		isInWSL = err == nil
	})
	return isInWSL
}

// IsWSLInstalled returns current environment has WSL guest/host
func IsWSLInstalled() bool {
	return IsInWSL()
}

// WSLGuest returns environment to access Guest(Linux) environment
func WSLGuest() Environment {
	return &nonVirtualEnvironment{}
}

// WSLHost returns environment to access Host(Windows) environment
func WSLHost() Environment {
	if IsInWSL() {
		return &wslHostEnvironment{}
	}
	return &nonVirtualEnvironment{}
}

type wslHostEnvironment struct {
}

var _ Environment = &nonVirtualEnvironment{}

func (e wslHostEnvironment) Exec(ctx context.Context, cmd string, args ...string) *exec.Cmd {
	pathExt, ok := e.Environ()["PATHEXT"]
	exts := []string{"", ".com", ".exe", ".bat", ".cmd", ".lnk"}
	if ok {
		exts = []string{}
		for _, ext := range strings.Split(pathExt, ";") {
			exts = append(exts, strings.TrimSpace(ext))
		}
	}
	for _, ext := range exts {
		// todo: reimplement exec.LookPath to overwrite
		// PATH environment variable by using host's PATH
		_, err := exec.LookPath(cmd + ext)
		if err == nil {
			cmd = cmd + ext
			break
		}
	}
	return exec.CommandContext(ctx, cmd, args...)
}

func (e wslHostEnvironment) Open(input string) {
	// wslview is better because cmd.exe /C set specify
	// static location of windows. but wsvar reset terminal
	// unexpectedly
	convertedPath := ConvertToHostPath(input)
	cmd := exec.Command("/mnt/C/Windows/System32/cmd.exe", "/C", "start", convertedPath)
	cmd.Run()
}

func (e wslHostEnvironment) UserHomeDir() (string, error) {
	envs := e.Environ()
	home, ok := envs["USERPROFILE"]
	if ok {
		return home, nil
	}
	home, ok = envs["HOMEPATH"]
	if ok {
		return home, nil
	}
	return "", errors.New("not found")
}

func (e wslHostEnvironment) UserConfigDir() (string, error) {
	envs := e.Environ()
	cache, ok := envs["APPDATA"]
	if ok {
		return cache, nil
	}
	return "", errors.New("not found")
}

func (e wslHostEnvironment) UserCacheDir() (string, error) {
	envs := e.Environ()
	cache, ok := envs["LOCALAPPDATA"]
	if ok {
		return cache, nil
	}
	return "", errors.New("not found")
}

var cachedHostEnv map[string]string
var checkEnvOnce sync.Once

func (e wslHostEnvironment) Environ() map[string]string {
	checkEnvOnce.Do(func() {
		cachedHostEnv = make(map[string]string)
		// wslvar is better because cmd.exe /C set specify
		// static location of windows. but wsvar reset terminal
		// unexpectedly
		//cmd := exec.Command("wslvar", "--getsys")
		cmd := exec.Command("/mnt/c/Windows/System32/cmd.exe", "/C", "set")
		vars, err := cmd.Output()
		if err != nil {
			return
		}
		for _, vr := range strings.Split(string(vars), "\n") {
			f := strings.SplitN(vr, "=", 2)
			if len(f) == 2 {
				cachedHostEnv[f[0]] = strings.TrimSpace(f[1])
			}
		}
	})
	return cachedHostEnv
}

// ConvertToHostPath returns host style path if current env is in WSL.
// Otherwise, return input path as is.
// Docker doesn't support this.
func ConvertToHostPath(path string) string {
	cmd := exec.Command("wslpath", "-w", path)
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
	cmd := exec.Command("wslpath", path)
	result, err := cmd.Output()
	if err != nil {
		return path
	}
	return strings.TrimSpace(string(result))
}
