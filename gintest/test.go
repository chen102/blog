package gintest

import (
	"blog/model"
	"blog/router"
	"blog/tool"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	//"gopkg.in/yaml.v2"
	"io"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

const (
	Functional      = iota //功能测试
	LocalData              //局部数据结构
	IndependentPath        //独立路径
	Error                  //各种错误处理
	Boundary               //临界值
)

type TestCase struct {
	CaseType                      uint        //测试用例类型
	Description                   string      //描述
	Method, Url, BodyType, Cookie string      //用例名称,方法，地址，body类型
	Param, Exp                    interface{} //参数，期望输出
}
type TestCases []TestCase

func NewTest() *gin.Engine {
	config := tool.NewConfig()
	config.TestReadConfig()
	model.DelMysql(*config)
	model.DelRedis(*config)
	return router.New()
}

//JSON转map
func JsonToMap(str []byte) (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	err := json.Unmarshal(str, &tempMap)
	if err != nil {
		return nil, err
	}
	return tempMap, nil
}
func JsonMapParam(param map[string]interface{}) (io.Reader, error) {
	j, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(j), nil
}
func NewRequest(method, url, bodytype, cookie string, body map[string]interface{}) (*http.Request, error) {
	if method == "GET" {
		return http.NewRequest(method, url, nil)
	}
	switch bodytype {
	case "JSON":
		jsonBody, err := JsonMapParam(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(method, url, jsonBody)
		if err != nil {
			return nil, err
		}
		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}

		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		return req, nil
	case "FORM":
		//form表单的参数可以通过querystring的形式符在URI地址后面进行传递
		null := strings.NewReader("")
		req, err := http.NewRequest(method, url+ParamToStr(body), null) //body不能输入nil,gin报空body
		if err != nil {
			return nil, err
		}
		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
		return req, nil
	}
	return nil, fmt.Errorf("bodytype:JSON OR FORM")
}
func StartHandler(router *gin.Engine, req *http.Request) (string, http.Header, error) {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	//result := w.Result()
	//defer result.Body.Close()
	//bodyByte, err := ioutil.ReadAll(result.Body) //返回[]byte
	//if err != nil {
	//return nil, err
	//}

	return w.Body.String(), w.HeaderMap, nil
}

//输出成querystring形式
func ParamToStr(mp map[string]interface{}) string {
	if len(mp) == 0 {
		return ""
	}
	values := ""
	for k, v := range mp {
		switch v.(type) {
		case int:
			values += "&" + k + "=" + strconv.Itoa(v.(int))
		case uint:
			values += "&" + k + "=" + strconv.Itoa(int(v.(uint)))
		case string:
			values += "&" + k + "=" + v.(string)
		case []string: //若是切片，则key名相同
			for _, v1 := range v.([]string) {
				values += "&" + k + "=" + v1
			}
		}
	}
	temp := values[1:]
	values = "?" + temp
	return values
}
