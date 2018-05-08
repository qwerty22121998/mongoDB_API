package main

import (
	"testing"
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/json"
)

const host = "http://localhost:8080"

func readResp(r *http.Response) string {
	defer r.Body.Close()
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err.Error()
	}
	return string(str)

}

type Test struct {
	Type string
	Link string
	Req  interface{} // io.reader
	Resp Request
}

func TestApi(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodDelete, host+"/student", nil)
	client.Do(req)
	// ////
	students := []Student{
		{"1", "A", 15},
		{"2", "B", 16},
		{"3", "C", 17},
		{"1", "A", 10},
		{},
		{"1", "Change", 10},
		{},
		{},
		{"1", "Change", 11},
	}
	methods := []string{
		"POST",
		"POST",
		"POST",
		"POST",
		"GET",
		"PATCH",
		"DELETE",
		"DELETE",
		"PATCH",
	}
	links := []string{
		host + "/insert",
		host + "/insert",
		host + "/insert",
		host + "/insert",
		host + "/student",
		host + "/update",
		host + "/delete/4",
		host + "/delete/1",
		host + "/update",
	}
	rightResp := []Request{
		{200, "Success", students[0]},
		{200, "Success", students[1]},
		{200, "Success", students[2]},
		{400, "Duplicated", nil},
		{200, "Success", []Student{students[0], students[1], students[2]}},
		{200, "Success", students[5]},
		{404, "Not found", nil},
		{200, "Success", students[5]},
		{404, "Not found", nil},
	}



	test := make([]Test, 0)
	for i := range students {
		test = append(test, Test{methods[i], links[i], students[i], rightResp[i]})
	}

	// ////

	for _, v := range test {
		js, _ := json.Marshal(v.Req.(Student))
		req, _ := http.NewRequest(v.Type, v.Link, strings.NewReader(string(js)))
		fmt.Println(string(js))
		res := ""
		resp, err := client.Do(req)
		if err != nil {
			res = err.Error()
		} else {
			res = readResp(resp)
		}

		expect, _ := json.Marshal(v.Resp)

		fmt.Printf("Result\n%s\nExpected:\n%s\n\n", res, fmt.Sprintf("%v", string(expect)))
		if res != fmt.Sprintf("%v", string(expect)) {
			t.Errorf("Wrong case")
			return
		}
	}
}
