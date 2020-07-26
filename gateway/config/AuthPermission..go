package config

import (
	"regexp"
	"strings"
)

var (
	AuthPermitConfig AuthPermitAll
)

//Http配置
type AuthPermitAll struct {
	PermitAll []interface{}
}

func Match(str string) bool {
	if count := len(AuthPermitConfig.PermitAll); count > 0 {
		for i := 0; i < count; i++ {
			s := AuthPermitConfig.PermitAll[i].(string)
			res, _ := regexp.MatchString(strings.ReplaceAll(s, "**", "(.*?)"), str)
			if res {
				return true
			}
		}
	}
	
	return false
}