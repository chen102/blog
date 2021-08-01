package tool

import (
	"strings"
	"time"
)

func ShortTime() string {
	const shortForm = "2006-01-02 15:04:05"
	now := time.Now()
	tepm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local)
	str := tepm.Format(shortForm)
	return str

}

//字符串拼接
func StrSplicing(str1, str2 string) string {
	var build strings.Builder
	build.WriteString(str1)
	build.WriteString(str2)
	return build.String()

}
