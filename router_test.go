package main

import (
	"blog/gintest"
	//"encoding/json"
	"fmt"
	//"github.com/stretchr/testify/assert"
	//"net/http"
	//"net/http/httptest"
	//"strings"
	"testing"
)

//官方gin单元测试Demo
//func TestPingRoute(t *testing.T) {
//r := New()
//w := httptest.NewRecorder()
//req, _ := http.NewRequest("GET", "/ping", nil)
//r.ServeHTTP(w, req)
//assert.Equal(t, 200, w.Code)
//body := w.Body.String()
//bodymap := make(map[string]interface{})
//err := json.Unmarshal([]byte(body), &bodymap)
//if err != nil {
//t.Errorf("string to map error")
//}
//fmt.Println(bodymap["code"])
//assert.Equal(t, "pong", w.Body.String())
//}

func TestArticleAddRoute(t *testing.T) {
	r := gintest.NewTest()
	articleParam := map[string]interface{}{
		"AuthorId":       1,
		"ArticleTitle":   "redis教程",
		"ArticleContent": "先这样，再这样，再那样",
		"Tags":           []string{"redis", "nosql"},
	}
	testCases := gintest.TestCases{
		{CaseName: "C01", Method: "POST", Url: "/article/manage/add", BodyType: "JSON", Param: articleParam, Exp: false},
		{CaseName: "C02", Method: "POST", Url: "/article/manage/add", BodyType: "FORM", Param: articleParam, Exp: false},
	}
	for _, v := range testCases {

		fmt.Println(v.CaseName, ": Testing")
		req, err := gintest.NewRequest(v.Method, v.Url, v.BodyType, v.Param.(map[string]interface{}))
		if err != nil {
			t.Error(err.Error())
		}
		bodyBody, err := gintest.StartHandler(r, req)
		if err != nil {
			t.Error(err.Error())
		}

		fmt.Println(bodyBody)
	}
	//assert.Equal(t, 0, 0)
}
func TestArticleShowRoute(t *testing.T) {
	r := gintest.NewTest()
	param := map[string]interface{}{
		"AuthorId":  1,
		"ArticleId": 1,
	}
	testCases := gintest.TestCases{
		{CaseName: "C03", Method: "POST", Url: "/article/manage/show", BodyType: "JSON", Param: param, Exp: false},
		{CaseName: "C04", Method: "POST", Url: "/article/manage/show", BodyType: "FORM", Param: param, Exp: false},
	}
	for _, v := range testCases {

		fmt.Println(v.CaseName, ": Testing")
		req, err := gintest.NewRequest(v.Method, v.Url, v.BodyType, v.Param.(map[string]interface{}))
		if err != nil {
			t.Error(err.Error())
		}
		bodyBody, err := gintest.StartHandler(r, req)
		if err != nil {
			t.Error(err.Error())
		}

		fmt.Println(bodyBody)
	}
	//assert.Equal(t, 0, 0)
}
