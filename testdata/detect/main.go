package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	alicein "github.com/shibukawa/alice-in"
)

func main() {
	var hi string
	var gi string
	if runtime.GOOS == "windows" {
		hi = "田 "
		gi = "🐧 "
	} else if runtime.GOOS == "linux" {
		if alicein.IsInWSL() {
			hi = "田 "
		} else {
			hi = "🐧 "
		}
		gi = "🐧 "
	} else if strings.Contains(runtime.GOOS, "bsd") {
		hi = "👿 "
		gi = "👿 "
	} else if runtime.GOOS == "darwin" {
		hi = "🍎 "
		gi = "🍎 "
	}

	fmt.Printf("🐋 alicein.IsInDocker(): %v\n", alicein.IsInDocker())
	fmt.Printf("%s aliciein.IsInWSL(): %v\n", hi, alicein.IsInWSL())

	fmt.Printf("%s WSL Guest: alicein.WSLGuest()\n", gi)
	guest := alicein.WSLGuest()

	homePath, _ := guest.UserHomeDir()
	fmt.Printf("   %sUserHomeDir(): %s\n", gi, homePath)
	configPath, _ := guest.UserConfigDir()
	fmt.Printf("   %sUserConfigDir(): %s\n", gi, configPath)
	cachePath, _ := guest.UserCacheDir()
	fmt.Printf("   %sUserCacheDir(): %s\n", gi, cachePath)
	if alicein.IsInWSL() && runtime.GOOS == "linux" {
		fmt.Println("   田 Exec(context.Background(), \"calc\")")
		cmd := guest.Exec(context.Background(), "calc")
		res, err := cmd.Output()
		fmt.Printf("      result = '%s', err = %v\n", string(res), err)
	}

	fmt.Printf("\n%s WSL Host: alicein.WSLHost()\n", hi)
	host := alicein.WSLHost()

	homePath, _ = host.UserHomeDir()
	fmt.Printf("   %sUserHomeDir(): %s\n", hi, homePath)
	configPath, _ = host.UserConfigDir()
	fmt.Printf("   %sUserConfigDir(): %s\n", hi, configPath)
	cachePath, _ = host.UserCacheDir()
	fmt.Printf("   %sUserCacheDir(): %s\n", hi, cachePath)
	if alicein.IsInWSL() && runtime.GOOS == "windows" {
		fmt.Println("   田 Exec(context.Background(), \"date\")")
		cmd := host.Exec(context.Background(), "date")
		res, err := cmd.Output()
		fmt.Printf("      result = '%s', err = %v\n", string(res), err)
	}

	fmt.Println("\nPath conversion between 🐧 ⇔ 田")
	guestPath := alicein.ConvertToGuestPath(`/usr/bin/yes`)
	fmt.Printf("🐧 alicein.ConvertToGuestPath(`/usr/bin/yes`): %s\n", guestPath)
	guestPath = alicein.ConvertToGuestPath(`C:\\windows\\system32`)
	fmt.Printf("🐧 alicein.ConvertToGuestPath(`C:\\windows\\system32`): %s\n", guestPath)
	hostPath := alicein.ConvertToHostPath(`/usr/bin/echo`)
	fmt.Printf("田 alicein.ConvertToHostPath(`/usr/bin/echo`): %s\n", hostPath)
	hostPath = alicein.ConvertToHostPath(`C:\\windows\\system32`)
	fmt.Printf("田 alicein.ConvertToHostPath(`C:\\windows\\system32`): %s\n", hostPath)

}
