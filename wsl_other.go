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

// WSLGuestEnvironemnt returns environment to access Host(Windows) environment
func WSLGuestEnvironemnt() Environment {
	return &nonVirtualEnvironment{}
}

// WSLHostEnvironemnt returns environment to access Host(Windows) environment
func WSLHostEnvironemnt() Environment {
	return &nonVirtualEnvironment{}
}

// DetectEnvType returns detected parent environment type
//
// This function check's parent process and return result
func DetectEnvType() EnvType {
	return NativeEnv
}

// DetectedEnvironment returns environment of parent
func DetectedEnvironment() Environment {
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
