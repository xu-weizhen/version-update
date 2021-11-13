package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var updateVersionCodes []string

func main() {
	var buffer []byte
	// 自定义一个访问连接
	urlString := "http://127.0.0.1:8080/newversion?aid=10000&channel=huawei&cpu_arch=32&device_id=9V4ZMM7I8DKW6DOZ&device_platform=iOS&os_api=20&update_version_code=1.0.0.0&version=-v-&version_code=0.0.0.0"
	resp, err := http.Get(urlString)
	if err != nil {
		log.Fatal(err)
	}
	if buffer, err = ioutil.ReadAll(resp.Body); err != nil && err != io.EOF {
		log.Fatal(err)
	}
	str := string(buffer)
	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
	if str != "null" {
		fmt.Println(urlString)
		fmt.Println(str)
	}
}
