package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/shibukawa/go-isisolated"
)

func main() {
	fmt.Printf("🐋 isisolated.IsInDocker(): %v\n", isisolated.IsInDocker())
	fmt.Printf("🐧 isisolated.IsInWSL(): %v\n", isisolated.IsInWSL())

	guestPath := isisolated.ConvertToGuestPath(`/usr/bin/yes`)
	fmt.Printf("🐧 isisolated.ConvertToGuestPath(`/usr/bin/yes`): %s\n", guestPath)
	guestPath = isisolated.ConvertToGuestPath(`C:\windows\system32`)
	fmt.Printf("🐧 isisolated.ConvertToGuestPath(`C:\\windows\\system32`): %s\n", guestPath)
	hostPath := isisolated.ConvertToHostPath(`/usr/bin/yes`)
	fmt.Printf("田 isisolated.ConvertToHostPath(`/usr/bin/yes`): %s\n", hostPath)
	hostPath = isisolated.ConvertToHostPath(`C:\windows\system32`)
	fmt.Printf("田 isisolated.ConvertToHostPath(`C:\\windows\\system32`): %s\n", hostPath)
	if isisolated.IsInWSL() && runtime.GOOS == "linux" {
		fmt.Println("田 isisolated.ExecInHostEnv(context.Background(), \"calc\")")
		cmd := isisolated.ExecInHostEnv(context.Background(), "calc")
		res, err := cmd.Output()
		fmt.Printf("    result = '%s', err = %v\n", string(res), err)
	}
	envs := isisolated.EnvironInGuest()
	fmt.Println("🐧 isisolated.EnvironInGuest():")
	for k, v := range envs {
		fmt.Printf("  🔑 %s = 🗨 %s\n", k, v)
	}
	envs = isisolated.EnvironInHost()
	fmt.Println("田 isisolated.EnvironInHost():")
	for k, v := range envs {
		fmt.Printf("  🔑 %s = 🗨 %s\n", k, v)
	}
	homePath, _ := isisolated.UserHomeDirInGuest()
	fmt.Printf("🐧 isisolated.UserHomeDirInGuest(): %s\n", homePath)
	homePath, _ = isisolated.UserHomeDirInHost()
	fmt.Printf("田 isisolated.UserHomeDirInHost(): %s\n", homePath)
	configPath, _ := isisolated.UserConfigDirInGuest()
	fmt.Printf("🐧 isisolated.UserConfigDirInGuest(): %s\n", configPath)
	configPath, _ = isisolated.UserConfigDirInHost()
	fmt.Printf("田 isisolated.UserConfigDirInHost(): %s\n", configPath)
	cachePath, _ := isisolated.UserCacheDirInGuest()
	fmt.Printf("🐧 isisolated.UserCacheDirInGuest(): %s\n", cachePath)
	cachePath, _ = isisolated.UserCacheDirInHost()
	fmt.Printf("田 isisolated.UserCacheDirInHost(): %s\n", cachePath)
}
