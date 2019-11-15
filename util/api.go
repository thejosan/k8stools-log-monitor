package util

import (
	"net/url"
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
)

func Api_Post(url string, postValue url.Values) (string) {
	postString := postValue.Encode()
	req, err := http.NewRequest("POST",url, strings.NewReader(postString))
	if err != nil {
		fmt.Println("ERROR")
	}
	defer req.Body.Close()
	// 表单方式(必须)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//AJAX 方式请求
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "{\"Status\": \"error\",\"Message\": \"请求链接失败\"}"
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return string(body)
}
