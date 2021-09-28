package service

import (
	"blog/gintest"
	"blog/tool"
	"log"
	"testing"
)

func TestAddArticl(t *testing.T) {

	r := gintest.NewTest()
	//登录获取Cookie
	Cookies := gintest.GetCookie()
	tags := []string{"one", "two", "three"}

	testCases := gintest.TestCases{
		{CaseType: gintest.Functional, Description: "Case1:功能测试1", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[0], BodyType: "FORM", Param: map[string]interface{}{
			"ArticleTitle":   "BOLG",
			"ArticleContent": "this is a blog",
			"Tags":           tags,
		}, Exp: false},
		{CaseType: gintest.Functional, Description: "Case2:功能测试2", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[1], BodyType: "FORM", Param: map[string]interface{}{
			"ArticleTitle":   "BOLG1",
			"ArticleContent": "this is a blog1",
			"Tags":           tags,
		}, Exp: false},
		{CaseType: gintest.Error, Description: "Case3:错误处理测试1", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[1], BodyType: "FORM", Param: map[string]interface{}{
			"ArticleTitle":   "",
			"ArticleContent": "",
			"Tags":           new([]string),
		}, Exp: false},

		{CaseType: gintest.Error, Description: "Case4:错误处理测试2", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[1], BodyType: "FORM", Param: map[string]interface{}{}, Exp: false},
		{CaseType: gintest.Boundary, Description: "Case5:临界值测试1", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleTitle":   TitleBoundart(),
			"ArticleContent": "test",
		}, Exp: false},
		{CaseType: gintest.Boundary, Description: "Case5:临界值测试2", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleTitle":   TitleBoundart() + "哈",
			"ArticleContent": "test",
		}, Exp: false},
		//{CaseType: gintest.Boundary, Description: "Case6:临界值测试3", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
		//"ArticleTitle":   "test",
		//"ArticleContent": ContentBoundart(),
		//}, Exp: false},
		//{CaseType: gintest.Boundary, Description: "Case7:临界值测试4", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
		//"ArticleTitle":   "test",
		//"ArticleContent": ContentBoundart()[3:], //一个中文占3个字节
		//}, Exp: false},
		{CaseType: gintest.Boundary, Description: "Case8:临界值测试5", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleTitle":   "test",
			"ArticleContent": "test",
			"Tags":           TagsBoundart(),
		}, Exp: false},
		{CaseType: gintest.Boundary, Description: "Case9:临界值测试6", Method: "POST", Url: "/api/v0/article/add", Cookie: Cookies[2], BodyType: "JSON", Param: map[string]interface{}{
			"ArticleTitle":   "test",
			"ArticleContent": "test",
			"Tags":           TagsBoundart()[1:],
		}, Exp: false},
	}
	for _, v := range testCases {

		log.Println(v.Description, ": Testing")
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
func TagsBoundart() []string {
	return []string{"test", "test", "test", "test", "test", "test"}
}
func TitleBoundart() string {
	res := ""
	for i := 0; i < 20; i++ {
		res = tool.StrSplicing(res, "我")
	}
	return res
}
func ContentBoundart() string {
	res := "哈"
	for i := 0; i < 16; i++ {
		res = tool.StrSplicing(res, res)
	}
	return res
}
