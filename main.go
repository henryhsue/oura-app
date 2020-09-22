package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func main() {
	spew.Config.DisablePointerAddresses = true

	// process user info
	body, err := makeRequest("/v1/userinfo")
	check(err)
	var userInfo *UserInfo
	err = json.Unmarshal(body, &userInfo)
	check(err)
	spew.Dump("User Info: %+v\n", userInfo)

	// process sleep data
	body, err = makeRequest("/v1/sleep?start=2020-09-20&end=2020-09-21") // TODO: parameterize dates
	check(err)
	var sleepSummary interface{}
	err = json.Unmarshal(body, &sleepSummary)
	check(err)
	spew.Dump("Sleep Summary: %+v\n", sleepSummary)
}

type UserInfo struct {
	Age    int32
	Weight float32
	Height float32
	Gender string
	Email  string
}

// TODO: add only necessary fields in sleep summary, with struct

func makeRequest(path string) ([]byte, error) {
	// create request
	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.URL.Scheme = "https"
	req.URL.Host = "api.ouraring.com"
	token, err := ioutil.ReadFile("bearer.token")
	check(err)
	req.Header.Add("Authorization", string(token))

	// process request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}
