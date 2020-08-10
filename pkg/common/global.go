package common

import (
	"go/build"
	"os"
	"runtime"
	"time"
	"strings"
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

// GetFirstDateOfMonth 月初
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day() + 1)
	return GetZeroTime(d)
}
// GetLastDateOfMonth 月末
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetMondayDate - 周一
func GetMondayDate(d time.Time) time.Time {
	offset := int(time.Monday - d.Weekday())
    if offset > 0 {
        offset = -6
	}
	return GetZeroTime(d).AddDate(0, 0, offset)
}
 
// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// FormatTime - 格式化日期
func FormatTime(d time.Time, format string) string {
	m := map[string]string{
		"Y": "2000",
		"y": "00",
		"m": "02",
		"d": "02",
		"H": "20",
		"i": "00",
		"s": "00", 
	}
	for k, v := range m {
		format = strings.Replace(format, k, v, -1)
    }
	return d.Format(format)
}