# is-isolated

## Sample execution (Linux application inside WSL)

```sh
$ go run testdata/detect/main.go 
🐋 isisolated.IsInDocker(): false
🐧 isisolated.IsInWSL(): true
🐧 isisolated.ConvertToGuestPath(`/usr/bin/yes`): /usr/bin/yes
🐧 isisolated.ConvertToGuestPath(`C:\windows\system32`): /mnt/c/windows/system32
田 isisolated.ConvertToHostPath(`/usr/bin/yes`): \\wsl$\Ubuntu-20.04\usr\bin\yes
田 isisolated.ConvertToHostPath(`C:\windows\system32`): C:\windows\system32
田 isisolated.ExecInHostEnv(context.Background(), "calc")
    result = '', err = <nil>
🐧 isisolated.EnvironInGuest():
  🔑 VSCODE_GIT_ASKPASS_MAIN = 🗨 /home/shibu/.vscode-server/bin/3dd905126b34dcd4de81fa624eb3a8cbe7485f13/extensions/git/dist/askpass-main.js
  🔑 VSCODE_IPC_HOOK_CLI = 🗨 /tmp/vscode-ipc-2e573d81-3bb8-4908-b35e-3c522b07227a.sock
  🔑 APPLICATION_INSIGHTS_NO_DIAGNOSTIC_CHANNEL = 🗨 true
  🔑 MOTD_SHOWN = 🗨 update-motd
  🔑 LANG = 🗨 C.UTF-8
  :
  🔑 VERBOSE_LOGGING = 🗨 true
  🔑 TERM_PROGRAM = 🗨 vscode
  🔑 OLDPWD = 🗨 /home/shibu/develop/go-isisolated
田 isisolated.EnvironInHost():
  🔑 VSCODE_WSL_EXT_LOCATION = 🗨 c:\Users\yoshi\.vscode\extensions\ms-vscode-remote.remote-wsl-0.44.4
  🔑 DriverData = 🗨 C:\Windows\System32\Drivers\DriverData
  🔑 ORIGINAL_XDG_CURRENT_DESKTOP = 🗨 undefined
  🔑 PROCESSOR_ARCHITECTURE = 🗨 AMD64
  :
  🔑 windir = 🗨 C:\WINDOWS
  🔑 ComSpec = 🗨 C:\WINDOWS\system32\cmd.exe
  🔑 HOMEPATH = 🗨 \Users\yoshi
  🔑 VSCODE_WSL_DISTRO = 🗨 Ubuntu-20.04
🐧 isisolated.UserHomeDirInGuest(): /home/shibu
田 isisolated.UserHomeDirInHost(): C:\Users\yoshi
🐧 isisolated.UserConfigDirInGuest(): C:\Users\yoshi\AppData\Roaming
田 isisolated.UserConfigDirInHost(): C:\Users\yoshi\AppData\Roaming
🐧 isisolated.UserCacheDirInGuest(): /home/shibu/.cache
田 isisolated.UserCacheDirInHost(): C:\Users\yoshi\AppData\Local
```