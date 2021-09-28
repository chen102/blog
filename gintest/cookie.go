package gintest

func GetCookie() map[int]string {
	r := NewTest()
	Cookies := make(map[int]string, 3)
	//登录获取Cookie
	login := TestCases{
		{CaseType: Functional, Description: "获取cookie", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{

			"account":  "66666666",
			"password": "12341234",
		}, Exp: false},
		{CaseType: Functional, Description: "获取cookie", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{
			"account":  "user00004",
			"password": "12341234",
		}, Exp: false},
		{CaseType: Functional, Description: "获取cookie", Method: "POST", Url: "/api/v0/user/login", BodyType: "JSON", Param: map[string]interface{}{
			"account":  "66666666",
			"password": "12341234",
		}, Exp: false},
	}
	for k, v := range login {
		req, err := NewRequest(v.Method, v.Url, v.BodyType, "", v.Param.(map[string]interface{}))
		if err != nil {
			panic(err)
		}
		_, headers, err := StartHandler(r, req)
		Cookies[k] = headers["Set-Cookie"][0]

	}
	return Cookies
}
