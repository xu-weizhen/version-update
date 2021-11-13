package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var updateVersionCodes []string

func main() {
	rand.Seed(time.Now().Unix())
	platformFunc := platformFuncCreate()
	osApiFunc := osApiFuncCreate()
	channelFunc := channelFuncCreate()
	deviceIDFunc := deviceIDFuncCreate()
	versionCodeFunc := versionCodeFuncCreate()
	appIDFunc := appIDFuncCreate()
	cpuArchFunc := cpuArchFuncCreate()
	//for {
	//log.Println(osApiFunc())
	//log.Println(download)
	var buffer []byte
	//keys := []string{"version", "device_platform", "device_id", "os_api", "channel", "version_code", "update_version_code", "aid", "cpu_arch"}
	for {
		params := url.Values{
			"version":             {"-v-"},
			"device_platform":     {platformFunc()},
			"device_id":           {deviceIDFunc()},
			"os_api":              {osApiFunc()},
			"channel":             {channelFunc()},
			"version_code":        {versionCodeFunc()},
			"update_version_code": {versionCodeFunc()},
			"aid":                 {appIDFunc()},
			"cpu_arch":            {cpuArchFunc()},
		}
		//fmt.Println("\n\n")
		urlString := "http://127.0.0.1:8080/newversion?" + params.Encode()
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
			//for _, key := range keys {
			//	fmt.Println(params.Get(key))
			//}
			fmt.Println(str)
		}
	}

}

// platformFuncCreate 随机平台版本
func platformFuncCreate() func() string {
	platform := []string{"Android", "iOS"}
	return func() string {
		return platform[rand.Intn(2)]
	}
}

// osApiFuncCreate 随机 os_api
func osApiFuncCreate() func() string {
	osApiRange := []int{1, 20}
	return func() string {
		result := rand.Intn(osApiRange[1]-osApiRange[0]+1) + osApiRange[0]
		return strconv.Itoa(result)
	}
}

// channelFuncCreate 随机 channel
func channelFuncCreate() func() string {
	channelStrings := []string{"huawei", "xiaomi"}
	length := len(channelStrings)
	return func() string {
		return channelStrings[rand.Intn(length)]
	}
}

// deviceIDFuncCreate 从 createWhiteList 下16个whiteList.txt选择一个deviceID
func deviceIDFuncCreate() func() string {
	var selectSlices []string
	deviceIDFilePaths := make([]string, 0, 16)
	for i := 0; i < 16; i++ {
		deviceIDFilePaths = append(deviceIDFilePaths, fmt.Sprintf("../../createWhiteList/whiteList%d.txt", i))
	}

	for _, filePath := range deviceIDFilePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		contentBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		slices := strings.Split(string(contentBytes), " ")
		selectSlices = append(selectSlices, slices[:len(slices)-1]...)
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}

	length := len(selectSlices)
	return func() string {
		return selectSlices[rand.Intn(length)]
	}
}

// versionCodeFuncCreate 从 configModels 下versionFile.txt选取一个版本
func versionCodeFuncCreate() func() string {
	filePath := "../../configModels/versionFile.txt"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	versionCodeSlices := strings.Split(string(content), " ")
	length := len(versionCodeSlices) - 1 // 去掉最后的空白位
	return func() string {
		return versionCodeSlices[rand.Intn(length)]
	}
}

// appidFuncCreate 创建appid
func appIDFuncCreate() func() string {
	const appidErrorProbability = 0.1
	const correctAppID = "10000"
	const errorAppID = "100"
	return func() string {
		if rand.Float32() < appidErrorProbability {
			return errorAppID
		} else {
			return correctAppID
		}
	}
}

// cpuArchFuncCreate 随机 cpu_arch 字符串
func cpuArchFuncCreate() func() string {
	cpuStrings := []string{"32", "64"}
	return func() string {
		return cpuStrings[rand.Intn(2)]
	}
}
