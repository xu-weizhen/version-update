package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const WhiteMapSize = 1000
const oneWhiteLength = 16

func main() {
	//start := time.Now()
	rand.Seed(time.Now().Unix())
	selectRunes := make([]byte, 36)
	pos := 0
	for pos < 10 {
		selectRunes[pos] = byte(pos + '0')
		pos++
	}
	for pos < 36 {
		selectRunes[pos] = byte(pos - 10 + 'A')
		pos++
	}

	var f func(int, string)
	f = func(num int, fileName string) {
		var file *os.File
		var err error
		if !checkFileIsExist(fileName) {
			if file, err = os.Create(fileName); err != nil {
				log.Fatal(err)
			}
		} else {
			if file, err = os.OpenFile(fileName, os.O_WRONLY, 0777); err != nil {
				log.Fatal(err)
			}
		}
		defer file.Close()
		whiteMap := make(map[string]struct{}, WhiteMapSize)
		buffer := make([]byte, oneWhiteLength+1)
		buffer[oneWhiteLength] = ' '
		var pos int
		var ok bool
		for len(whiteMap) < num {
			for pos = 0; pos < oneWhiteLength; pos++ {
				buffer[pos] = selectRunes[rand.Intn(36)]
			}
			if _, ok = whiteMap[string(buffer)]; !ok {
				whiteMap[string(buffer)] = struct{}{}
				if _, err = file.Write(buffer); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	for i := 0; i < 16; i++ {
		f(WhiteMapSize, fmt.Sprintf("whiteList%d.txt", i))
	}
	//fmt.Println(time.Since(start))
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
