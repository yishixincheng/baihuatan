package common

import (
	"go/build"
	"os"
	"runtime"
	"time"
	"strings"
	"math/rand"
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

// UpsetSlice 打乱切片
func UpsetSlice(sli []interface{}) (aft []interface{}) {
	num := len(sli)
	if num == 0 {
		return
	}
	rand.Seed(time.Now().UnixNano()) // 随机种子
	sdIdxList := []int{}
	C: 
	for i := 0; i < num; i++ {
		rd := rand.Intn(int(num)) // 产生随机数
		for _, av := range sdIdxList {
			if rd == av {
				continue C
			}
		}
		sdIdxList = append(sdIdxList, rd)
	}

	for _, idx := range sdIdxList {
		aft = append(aft, sli[idx])
	}
	diffIdxList := []int{}
	// 求差集
	for i := 0; i < num; i++ {
		isExist := false
		for _, av := range sdIdxList {
			if i == av {
				isExist = true
				break
			}
		}
		if !isExist {
			diffIdxList = append(diffIdxList, i)
		}
	}
	
	if len(diffIdxList) != 0 {
		tmpSli := []interface{}{}
		for _, i := range diffIdxList {
			tmpSli = append(tmpSli, sli[i])
		}
		aft = append(aft,UpsetSlice(tmpSli)...)
	}

	return
}