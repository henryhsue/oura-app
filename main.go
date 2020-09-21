package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// create request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "/v1/userinfo", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.URL.Scheme = "https"
	req.URL.Host = "api.ouraring.com"
	token, err := ioutil.ReadFile("bearer.token")
	check(err)
	req.Header.Add("Authorization", string(token))

	// make request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
