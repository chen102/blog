package service

import (
	"blog/gintest"
	"log"
	"testing"
)

func TestAddComment(t *testing.T) {
	r := gintest.NewTest()
	Cookies := make(map[int]string, 3)
	//登录获取Cookie
	login := gintest.TestCases{
		{CaseName: "用户2", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{

			"account":  "66666666",
			"password": "12341234",
		}, Exp: false},
		{CaseName: "用户6", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{
			"account":  "user00004",
			"password": "12341234",
		}, Exp: false},
		{CaseName: "用户3", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{
			"account":  "66666666",
			"password": "12341234",
		}, Exp: false},
	}
	for k, v := range login {
		req, err := gintest.NewRequest(v.Method, v.Url, v.BodyType, "", v.Param.(map[string]interface{}))
		if err != nil {
			t.Error(err.Error())
		}
		_, headers, err := gintest.StartHandler(r, req)
		Cookies[k] = headers["Set-Cookie"][0]

	}
	testCases := gintest.TestCases{
		{CaseName: "功能测试1", Method: "POST", Url: "/api/v0/article/comment", Cookie: Cookies[0], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleId": 2,
			"Content":   "可以可以",
		}, Exp: false},
		{CaseName: "功能测试2", Method: "POST", Url: "/api/v0/article/comment", Cookie: Cookies[1], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleId": 2,
			"Content":   "写的好啊",
		}, Exp: false},
		{CaseName: "功能测试3", Method: "POST", Url: "/api/v0/article/comment", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleId":  2,
			"Rev":        1,
			"LandlordId": 1,
			"Content":    "GOOD",
		}, Exp: false},
		{CaseName: "功能测试4", Method: "POST", Url: "/api/v0/article/comment", Cookie: Cookies[1], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleId":  2,
			"Rev":        3,
			"LandlordId": 1,
			"Content":    "GOOD什么",
		}, Exp: false},
		{CaseName: "功能测试5", Method: "POST", Url: "/api/v0/article/comment", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleId":  2,
			"Rev":        4,
			"LandlordId": 1,
			"Content":    "没什么",
		}, Exp: false},
		{CaseName: "功能测试1", Method: "POST", Url: "/api/v0/article/comment", Cookie: Cookies[0], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleId":  2,
			"Rev":        5,
			"LandlordId": 1,
			"Content":    "你们真有意思",
		}, Exp: false},
	}

	for _, v := range testCases {

		log.Println(v.CaseName, ": Testing")
		req, err := gintest.NewRequest(v.Method, v.Url, v.BodyType, v.Cookie, v.Param.(map[string]interface{}))
		if err != nil {
			t.Error(err.Error())
		}
		body, headers, err := gintest.StartHandler(r, req)
		if err != nil {
			t.Error(err.Error())
		}

		log.Println(body, headers)

	}
}
