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

//多字符串拼接
func StrSplicing(str ...string) string {
	var build strings.Builder
	for _, v := range str {

		build.WriteString(v)
	}
	return build.String()

}
func SliceToString(str []string) string {
	var temp string
	for _, v := range str {
		if v != "" {

			temp += v + ","
		}
	}
	return temp
}
