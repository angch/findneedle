package main

// Not my (angch) code. Motionman's

import (
	"io/ioutil"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("TestData.txt")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(data); i += 8 {
		if int(data[i+0])+
			int(data[i+1])+
			int(data[i+2])+
			int(data[i+3])+
			int(data[i+4])+
			int(data[i+5])+
			int(data[i+6])+
			int(data[i+7]) > 384 {
			log.Println("found!")
			break
		}
	}
}
