# alice-in -- Alice in VIRTUAL Land

This library is for treating virtual environment. This library can handle the following environments:

* WSL2
* Docker

For example,

* From host Windows, you can get environment variables in WSL and exec command in WSL
* From guest WSL Linux, you can get environment variables in Windows and exec command in Host Windows.

## Basic Usage

### Environment Detection

* `alicein.IsInWSL()`: Detect current Linux environment is on WSL
* `alicein.IsWSLInstalled()`: Detect current Windows environment has WSL
* `alicein.IsInDocker()`: Detect current environment is in Docker

## Path conversion

* `ConvertToHostPath(path string) string`

It returns host style path if current env is in WSL.
Otherwise, return input path as is.
Docker doesn't support this.

* `ConvertToGuestPath(path string) string`

It returns guest style path.
Otherwise, return input path as is.
Docker doesn't support this.


## Get Environment

* `alicein.WSLGuest() Environment`: On Windows and WSL is installed, it returns interface to handle WSL. Otherwise, it returns native environment.
* `alicein.WSLHost() Environment`: On Linux in WSL, it returns interface to handle Host Windows. Otherwise, it returns native environment.
* `alicein.Docker(target) Environment`: Detect current environment is in Docker. If target prefixed "image:", it means image name. Otherwise, it means running container name.

## Environment

`Environment` interface wraps `os` / `os/exec` packages. Docker environment only support `Exec()`.

### `Exec(ctx context.Context, command string, args ...string) *exec.Cmd`

It execs commands in current environment.
In Windows's WSL host environment, exec via "wsl" command wrapper.
In Docker, exec via "docker exec".

### `Open(file string)`

It opens file by associated program.
On Windows, "start" is used. And "open" is used on mac,
and "xdg-open" is used on other environment.

### `UserHomeDir() (string, error)`

It returns in virtual environment.
Otherwise it returns standard os.UserHomeDir()

### `UserConfigDir() (string, error)`

It returns in virtual environment.
Otherwise it returns standard os.UserConfigDir()

### `UserCacheDir() (string, error)`

It returns in virtual environment.
Otherwise it returns standard os.UserCacheDir()

### `Environ() map[string]string`

It returns environment variables in virtual environment.
Otherwise it returns os.Environment

## Sample execution (Linux application inside WSL)

```sh
$ go run testdata/detect/main.go 
🐋 alicein.IsInDocker(): false
田  aliciein.IsInWSL(): true
🐧  WSL Guest: alicein.WSLGuest()
   🐧 UserHomeDir(): /home/shibu
   🐧 UserConfigDir(): /home/shibu/.config
   🐧 UserCacheDir(): /home/shibu/.cache
   田 Exec(context.Background(), "calc")

田  WSL Host: alicein.WSLHost()
   田 UserHomeDir(): C:\Users\yoshi
   田 UserConfigDir(): C:\Users\yoshi\AppData\Roaming
   田 UserCacheDir(): C:\Users\yoshi\AppData\Local

Path conversion between 🐧 ⇔ 田
🐧 alicein.ConvertToGuestPath(`/usr/bin/yes`): /usr/bin/yes
🐧 alicein.ConvertToGuestPath(`C:\windows\system32`): /mnt/c/windows/system32
田 alicein.ConvertToHostPath(`/usr/bin/echo`): \\wsl$\Ubuntu-20.04\usr\bin\echo
田 alicein.ConvertToHostPath(`C:\windows\system32`): C:\\windows\\system32
```

## License

Apache 2

## Contributors

* [@aodag](https://github.com/aodag): Cool name
* [@mopemope](https://github.com/mopemope)
* [@shirou](https://github.com/shirou)