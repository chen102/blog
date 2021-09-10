package tool

import (
	//"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
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
func IntSliceToString(i []int) string {
	var temp string
	for _, v := range i {
		temp += strconv.Itoa(v) + ","
	}
	return temp[:len(temp)-1]

}
func StringSliceTOIntSlice(str []string) []int64 {
	res := make([]int64, len(str))
	for k, v := range str {
		res[k], _ = strconv.ParseInt(v, 10, 64)
	}
	return res
}
func RandomTime() int64 {
	return rand.Int63n(time.Now().Unix()) + 7200 //随机生成现在到两个小时以内的时间戳

}

//分页工具
func Pagination(item []string, offset, count int) []string {
	if offset*count < 0 { //不能是负数
		return nil
	}
	leftover := offset + count
	if int(offset) >= len(item) { //若偏移量超了，直接返回
		return nil
	} else if leftover >= len(item) { //若输出的超出了已有的输出剩下的全部
		leftover = len(item) - offset
	}
	return item[offset:leftover]
}

//输出[]int
func PaginationINT(item []string, offset, count int) []int {
	toint := make([]int, len(item))
	for k, v := range item {
		toint[k], _ = strconv.Atoi(v)
	}
	if offset*count < 0 { //不能是负数
		return nil
	}
	leftover := offset + count
	if int(offset) >= len(item) { //若偏移量超了，直接返回
		return nil
	} else if leftover >= len(item) { //若输出的超出了已有的输出剩下的全部
		leftover = len(item) - offset
	}
	return toint[offset:leftover]
}
