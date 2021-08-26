package tool

import (
	//"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const standard = "2006-01-02 15:04:05"

func StringToTime(str string) time.Time {
	fmt.Println(str)
	st, _ := time.Parse(time.RFC3339, str)
	fmt.Println(st)
	return st
}
func TimeToString(t time.Time) string {
	return t.Format(standard)
}
func ShortTime() string {
	now := time.Now()
	tepm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local)
	str := tepm.Format(standard)
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
		temp += v + ","
	}
	return temp[:len(temp)-1] //最后一个,不要
}
func RandomTime() int64 {
	return rand.Int63n(time.Now().Unix()) + 7200 //随机生成现在到两个小时以内的时间戳

}
