package alicein

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// ErrNotImplemented indicates this method is not implemented yet
	ErrNotImplemented = errors.New("not implemented")
)

// IsInDocker returns current process works
func IsInDocker() bool {
	_, err := os.Stat(filepath.Join("/", ".dockerenv"))
	return err == nil
}

type dockerEnvironment struct {
	isImage    bool
	identifier string
}

var _ Environment = &dockerEnvironment{}

func (e dockerEnvironment) ConvertToHostPath(path string) string {
	return ""
}

func (e dockerEnvironment) ConvertToGuestPath(path string) string {
	return ""
}

func (e dockerEnvironment) Exec(ctx context.Context, cmd string, args ...string) *exec.Cmd {
	if e.isImage {
		cmdArgs := append([]string{}, "run", "-it", "--rm", e.identifier, cmd)
		cmdArgs = append(cmdArgs, args...)
		return exec.CommandContext(ctx, "docker", cmdArgs...)
	}
	cmdArgs := append([]string{}, "exec", e.identifier, cmd)
	cmdArgs = append(cmdArgs, args...)
	return exec.CommandContext(ctx, "docker", cmdArgs...)
}

func (e dockerEnvironment) Open(file string) {
}

func (e dockerEnvironment) UserHomeDir() (string, error) {
	return "", ErrNotImplemented
}

func (e dockerEnvironment) UserConfigDir() (string, error) {
	return "", ErrNotImplemented
}

func (e dockerEnvironment) UserCacheDir() (string, error) {
	return "", ErrNotImplemented
}

func (e dockerEnvironment) Environ() map[string]string {
	return nil
}

func (e dockerEnvironment) Type() EnvType {
	return DockerEnv
}

// Docker returns Docker environment.
//
// This environment only support Exec
func Docker(target string) Environment {
	return &dockerEnvironment{
		isImage:    strings.HasPrefix(target, "image:"),
		identifier: strings.TrimPrefix(target, "image:"),
	}
}
