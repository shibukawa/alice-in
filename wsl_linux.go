// +build linux

package isisolated

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
		_, err := os.Stat("/usr/bin/wslpath")
		isInWSL = err == nil
	})
	return isInWSL
}

// ConvertToHostPath returns host style path if current env is in WSL.
// Otherwise, return input path as is.
func ConvertToHostPath(path string) string {
	if IsInWSL() {
		cmd := exec.Command("wslpath", "-w", path)
		result, err := cmd.Output()
		if err != nil {
			return path
		}
		return strings.TrimSpace(string(result))
	}
	return path
}

// ConvertToHostPath returns guest style path.
// Otherwise, return input path as is.
func ConvertToGuestPath(path string) string {
	if IsInWSL() {
		cmd := exec.Command("wslpath", "-u", path)
		result, err := cmd.Output()
		if err != nil {
			return path
		}
		return strings.TrimSpace(string(result))
	}
	return path
}

func ExecInGuestEnv(ctx context.Context, command string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, command, args...)
}

func ExecInHostEnv(ctx context.Context, command string, args ...string) *exec.Cmd {
	if IsInWSL() {
		exts := []string{"", ".com", ".exe", ".bat", ".cmd", ".lnk"}
		for _, ext := range exts {
			_, err := exec.LookPath(command + ext)
			if err == nil {
				command = command + ext
				break
			}
		}
	}
	return exec.CommandContext(ctx, command, args...)
}

func UserCacheDirInGuest() (string, error) {
	return os.UserCacheDir()
}

func UserCacheDirInHost() (string, error) {
	if IsInWSL() {
		envs := EnvironInHost()
		cache, ok := envs["LOCALAPPDATA"]
		if ok {
			return cache, nil
		}
		return "", errors.New("not found")
	}
	return os.UserCacheDir()
}

func UserConfigDirInGuest() (string, error) {
	if IsInWSL() {
		envs := EnvironInHost()
		cache, ok := envs["APPDATA"]
		if ok {
			return cache, nil
		}
		return "", errors.New("not found")
	}
	return os.UserCacheDir()
}

func UserConfigDirInHost() (string, error) {
	if IsInWSL() {
		envs := EnvironInHost()
		cache, ok := envs["APPDATA"]
		if ok {
			return cache, nil
		}
		return "", errors.New("not found")
	}
	return os.UserConfigDir()
}

func UserHomeDirInGuest() (string, error) {
	return os.UserHomeDir()
}

func UserHomeDirInHost() (string, error) {
	if IsInWSL() {
		envs := EnvironInHost()
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
	return os.UserHomeDir()
}

func EnvironInGuest() map[string]string {
	envs := os.Environ()
	result := make(map[string]string, len(envs))
	for _, l := range envs {
		f := strings.SplitN(l, "=", 2)
		result[f[0]] = f[1]
	}
	return result
}

var cachedHostEnv map[string]string
var checkEnvOnce sync.Once

func EnvironInHost() map[string]string {
	if IsInWSL() {
		checkEnvOnce.Do(func() {
			cachedHostEnv = make(map[string]string)
			cmd := exec.Command("wslvar", "--getsys")
			vars, err := cmd.Output()
			if err != nil {
				return
			}
			for _, vr := range strings.Split(string(vars), "\n")[2:] {
				f := strings.SplitN(vr, " ", 2)
				if len(f) == 2 {
					cachedHostEnv[strings.TrimSpace(f[0])] = strings.TrimSpace(f[1])
				}
			}
		})
		return cachedHostEnv
	}
	return EnvironInGuest()
}
