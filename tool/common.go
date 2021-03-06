package tool

import (
	//"encoding/json"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const standard = "2006-01-02 15:04:05"

func StringToInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}
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

//压缩
func Compress(str string) string {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write([]byte(str))
	w.Close()
	return in.String()
}

//解压在
func UnCompress(str string) (string, error) {
	var out bytes.Buffer
	b := bytes.NewReader([]byte(str))
	res, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}
	io.Copy(&out, res)
	return out.String(), nil
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

//模拟队列操作
func Push(queue []int64, e int64) []int64 {
	queue = append(queue, e)
	return queue
}
func Pop(queue []int64) (int64, []int64) {
	temp := queue[0]
	queue = queue[1:]
	return temp, queue
}
