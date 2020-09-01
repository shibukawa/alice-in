package alicein

import (
	"context"
	"os/exec"
)

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
}
