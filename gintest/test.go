package gintest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
	"blog/model"
	"blog/router"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type TestCase struct {
	CaseName, Method, Url, BodyType string      //用例名称,方法，地址，body类型
	Param, Exp                      interface{} //参数，期望输出
}
type TestCases []TestCase

func NewTest() *gin.Engine {
	model.Del()
	model.Init()
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
func NewRequest(method, url, bodytype string, body map[string]interface{}) (*http.Request, error) {
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

		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		return req, nil
	case "FORM":
		//form表单的参数可以通过querystring的形式符在URI地址后面进行传递
		null := strings.NewReader("")
		req, err := http.NewRequest(method, url+ParamToStr(body), null) //body不能输入nil,gin报空body
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
		return req, nil
	}
	return nil, fmt.Errorf("bodytype:JSON OR FORM")
}
func StartHandler(router *gin.Engine, req *http.Request) (string, error) {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	//result := w.Result()
	//defer result.Body.Close()
	//bodyByte, err := ioutil.ReadAll(result.Body) //返回[]byte
	//if err != nil {
	//return nil, err
	//}

	return w.Body.String(), nil
}

//输出成querystring形式
func ParamToStr(mp map[string]interface{}) string {
	values := ""
	for k, v := range mp {
		switch v.(type) {
		case int: //真的蠢 感觉这样写
			values += "&" + k + "=" + strconv.Itoa(v.(int))
		case uint:
			values += "&" + k + "=" + strconv.Itoa(int(v.(uint)))
		case string:
			values += "&" + k + "=" + v.(string)
		}
	}
	temp := values[1:]
	values = "?" + temp
	return values
}
