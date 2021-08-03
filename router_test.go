package main

import (
	"blog/gintest"
	//"encoding/json"
	"fmt"
	//"net/http"
	//"net/http/httptest"
	//"github.com/stretchr/testify/assert"
	//"strings"
	"testing"
)

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

func TestShowRoute(t *testing.T) {
	r := gintest.NewTest()
	param := map[string]interface{}{
		"id": 1,
	}
	testCases := []struct {
		Casename              string
		Method, Url, Bodytype string
		Param                 interface{}
		Exp                   interface{}
	}{
		{"C01", "POST", "/article/manage/show", "JSON", param, false},
		{"C02", "POST", "/article/manage/show", "FORM", param, false},
	}
	for _, v := range testCases {

		fmt.Println(v.Casename, ": Testing")
		req, err := gintest.NewRequest(v.Method, v.Url, v.Bodytype, v.Param.(map[string]interface{}))
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
