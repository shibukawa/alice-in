package alicein

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

type nonVirtualEnvironment struct {
}

var _ Environment = &nonVirtualEnvironment{}

func (e nonVirtualEnvironment) Exec(ctx context.Context, cmd string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, cmd, args...)
}

func (e nonVirtualEnvironment) Open(file string) {
	open.Start(file)
}

func (e nonVirtualEnvironment) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (e nonVirtualEnvironment) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (e nonVirtualEnvironment) UserCacheDir() (string, error) {
	return os.UserCacheDir()
}

func (e nonVirtualEnvironment) Environ() map[string]string {
	return stdEnviron()
}

func (e nonVirtualEnvironment) Type() EnvType {
	return NativeEnv
}

func stdEnviron() map[string]string {
	envs := os.Environ()
	result := make(map[string]string, len(envs))
	for _, l := range envs {
		f := strings.SplitN(l, "=", 2)
		result[f[0]] = f[1]
	}
	return result
}

// NativeEnvironment returns environment to access current environment
func NativeEnvironment() Environment {
	return &nonVirtualEnvironment{}
}
