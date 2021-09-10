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

func TestUserRegister(t *testing.T) {
	r := gintest.NewTest()
	testCases := gintest.TestCases{
		{CaseName: "两次输入密码不同", Method: "POST", Url: "/api/v0/user/register", BodyType: "JSON", Param: map[string]interface{}{

			"username":    "chenmou",
			"account":     "11111111",
			"password":    "12341234",
			"reppassword": "12341235",
		}, Exp: false},
		{CaseName: "用户名已存在", Method: "POST", Url: "/api/v0/user/register", BodyType: "FORM", Param: map[string]interface{}{

			"username":    "chenmou",
			"account":     "11111111",
			"password":    "12341234",
			"reppassword": "12341234",
		}, Exp: false},
		{CaseName: "账户已存在", Method: "POST", Url: "/api/v0/user/register", BodyType: "JSON", Param: map[string]interface{}{

			"username":    "chenmou1",
			"account":     "11111111",
			"password":    "12341234",
			"reppassword": "12341234",
		}, Exp: false},
		{CaseName: "正常", Method: "POST", Url: "/api/v0/user/register", BodyType: "JSON", Param: map[string]interface{}{

			"username":    "chenmou3",
			"account":     "111111111",
			"password":    "12341234",
			"reppassword": "12341234",
		}, Exp: false},
	}
	for _, v := range testCases {

		fmt.Println(v.CaseName, ": Testing")
		req, err := gintest.NewRequest(v.Method, v.Url, v.BodyType, v.Param.(map[string]interface{}))
		if err != nil {
			t.Error(err.Error())
		}
		bodyBody, _, err := gintest.StartHandler(r, req)
		if err != nil {
			t.Error(err.Error())
		}

		fmt.Println(bodyBody)

	}
}
func TestUserLogin(t *testing.T) {
	r := gintest.NewTest()
	testCases := gintest.TestCases{
		{CaseName: "账户不对", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{

			"account":  "1111111315",
			"password": "12341234",
		}, Exp: false},
		{CaseName: "正常", Method: "POST", Url: "/api/v0/user/login", BodyType: "FORM", Param: map[string]interface{}{

			"account":  "111111111",
			"password": "12341234",
		}, Exp: false},
	}
	for _, v := range testCases {

		fmt.Println(v.CaseName, ": Testing")
		req, err := gintest.NewRequest(v.Method, v.Url, v.BodyType, v.Param.(map[string]interface{}))
		if err != nil {
			t.Error(err.Error())
		}
		bodyBody, _, err := gintest.StartHandler(r, req)
		if err != nil {
			t.Error(err.Error())
		}

		fmt.Println(bodyBody)

	}
}
