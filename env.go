package alicein

import (
	"context"
	"os/exec"
)

// EnvType is type of environment. DetectEnv() returns this
type EnvType int

const (
	// NativeEnv means current env equals to GOOS
	NativeEnv EnvType = iota + 1
	// WSLEnv means WSL guest environment
	WSLEnv
	// HostEnv means Host Windows environemnt
	HostEnv
	// DockerEnv means Docker environemnt
	DockerEnv
)

func (e EnvType) String() string {
	switch e {
	case NativeEnv:
		return "native"
	case WSLEnv:
		return "wsl"
	case HostEnv:
		return "host"
	}
	return "(unknown)"
}

// Environment is an interface that returns information or exec commands
type Environment interface {
	// Exec execs commands in current environment.
	// In Windows's WSL host environment, exec via "wsl" command wrapper.
	// In Docker, exec via "docker exec".
	Exec(ctx context.Context, command string, args ...string) *exec.Cmd

	// Open opens file by associated program.
	// On Windows, "start" is used. And "open" is used on mac,
	// and "xdg-open" is used on other environment.
	Open(file string)

	// UserHomeDir returns in virtual environment.
	// Otherwise it returns standard os.UserHomeDir()
	UserHomeDir() (string, error)

	// UserConfigDir returns in virtual environment.
	// Otherwise it returns standard os.UserConfigDir()
	UserConfigDir() (string, error)

	// UserCacheDir returns in virtual environment.
	// Otherwise it returns standard os.UserCacheDir()
	UserCacheDir() (string, error)

	// Environ returns environment variables in virtual environment.
	// Otherwise it returns os.Environment
	Environ() map[string]string

	// Type returns current environment type
	Type() EnvType
}
