package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func main() {
	// process user info
	body, err := makeRequest("/v1/userinfo")
	check(err)
	var userInfo *UserInfo
	err = json.Unmarshal(body, &userInfo)
	check(err)
	fmt.Printf("UserInfo: %+v\n", userInfo)
}

type UserInfo struct {
	Age    int32
	Weight float32
	Height float32
	Gender string
	Email  string
}

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
