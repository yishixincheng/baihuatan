package common

import (
	"go/build"
	"os"
	"runtime"
)

// IsWindows 判断是否是windows系统
func IsWindows() bool {

	sysType := runtime.GOOS
	if sysType == "windows" {

		// windows系统
		return true
	}

	return false
}

// GetGOPATH 获取gopath路径
func GetGOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}
