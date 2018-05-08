package main

import (
	"testing"
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/gin-gonic/gin/json"
	"strings"
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
	Resp string
}

func TestApi(t *testing.T) {
	client := &http.Client{}

	// // clear db
	req, err := http.NewRequest(http.MethodDelete, host+"/student", nil)

	resp, err := client.Do(req)

	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(readResp(resp))

	// insert
	students := []Student{
		Student{"1", "A", 15},
		Student{"2", "B", 16},
		Student{"3", "C", 17},
		Student{"1", "A", 10},
	}

	rightResp := []string{
		"{\"Status\":200,\"Msg\":\"Success\",\"Data\":{\"ID\":\"1\",\"Name\":\"A\",\"Age\":15}}",
		"{\"Status\":200,\"Msg\":\"Success\",\"Data\":{\"ID\":\"2\",\"Name\":\"B\",\"Age\":16}}",
		"{\"Status\":200,\"Msg\":\"Success\",\"Data\":{\"ID\":\"3\",\"Name\":\"C\",\"Age\":17}}",
		"{\"Status\":400,\"Msg\":\"Duplicated\",\"Data\":null}",
	}

	test := []Test{
		Test{"POST", host + "/insert", students[0], rightResp[0]},
		Test{"POST", host + "/insert", students[1], rightResp[1]},
		Test{"POST", host + "/insert", students[2], rightResp[2]},
		Test{"POST", host + "/insert", students[3], rightResp[3]},
	}

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

		fmt.Printf("Result\n%s\nExpected:\n%s\n", res, v.Resp)
		if res != v.Resp {
			t.Errorf("Wrong case")
			return
		}

	}

}
