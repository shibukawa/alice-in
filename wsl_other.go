// +build !linux,!windows

package alicein

// IsInWSL returns current environment is in WSL or not
func IsInWSL() bool {
	return false
}

// IsWSLInstalled returns current environment has WSL guest/host
func IsWSLInstalled() bool {
	return false
}

// WSLGuest returns environment to access Host(Windows) environment
func WSLGuest() Environment {
	return &nonVirtualEnvironment{}
}

// WSLHost returns environment to access Host(Windows) environment
func WSLHost() Environment {
	return &nonVirtualEnvironment{}
}

// ConvertToHostPath returns host style path if current env is in WSL.
// Otherwise, return input path as is.
// Docker doesn't support this.
func ConvertToHostPath(path string) string {
	return path
}

// ConvertToGuestPath returns guest style path.
// Otherwise, return input path as is.
// Docker doesn't support this.
func ConvertToGuestPath(path string) string {
	return path
}
