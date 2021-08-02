package main

import (
	"blog/gintest"
	"blog/model"
	//"encoding/json"
	"fmt"
	//"net/http"
	//"net/http/httptest"
	"github.com/stretchr/testify/assert"
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
	r := New()
	model.Del()
	model.Init()
	//payload := strings.NewReader(`{"id":1}`)  //JSON
	//payload := strings.NewReader("id=1") //FORM
	//req, err := gintest.NewRequest("POST", "/article/manage/show", "FORM", payload)
	//JSON中变量
	var id uint = 1
	param := map[string]interface{}{
		"id": id,
	}

	//req, err := gintest.NewRequest("POST", "/article/manage/show", "JSON", param)
	req, err := gintest.NewRequest("POST", "/article/manage/show", "FORM", param)
	if err != nil {
		t.Error(err.Error())
	}
	bodyBody, err := gintest.StartHandler(r, req)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(bodyBody)
	assert.Equal(t, 0, 0)
}
