package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/mhpixxio/konstruktor"
)

func main() {

	//create small data
	smalldata := konstruktor.CreateBigData(1, 1)
	//create big data
	log.Printf("creating bigdata ...\n")
	bigdata := konstruktor.CreateBigData(5, 100)
	log.Printf("finished creating bigdata. server is ready.\n")
	//adress
	httpurl := "http://localhost:4040"
	//define loopsizes
	loops := 3
	amountSmalldata := 5

	//Sending Big Data to Server
	for i := 0; i < loops; i++ {
		jsonclient(httpurl, "/postjson", bigdata)
	}
	//Receiving Big Data from Server
	for i := 0; i < loops; i++ {
		jsonclient(httpurl, "/getjson", smalldata)
	}
	//Sending Small Data to Server and Recieving Small Data
	for i := 0; i < loops; i++ {
		jsonclient(httpurl, "/postjson", smalldata)
	}
	//Sending a lot of Small Data to Server simultaniously
	var wg sync.WaitGroup
	for i := 0; i < loops; i++ {
		wg.Add(amountSmalldata)
		for j := 0; j < amountSmalldata; j++ {
			go func() {
				jsonclient(httpurl, "/postjson", smalldata)
				defer wg.Done()
			}()
		}
		wg.Wait()
	}
	//Sending a lot of Small Data to Server after one another
	for i := 0; i < loops; i++ {
		for j := 0; j < amountSmalldata; j++ {
			jsonclient(httpurl, "/postjson", smalldata)
		}
	}
}

func jsonclient(httpurl string, endpoint string, data []konstruktor.RandomData) {
	//Serialisierung
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//connection
	httpposturl := httpurl + endpoint
	fmt.Println("HTTP JSON POST URL:", httpposturl)
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer response.Body.Close()
	//fmt.Println("sent data: ", string(jsonData)) //print the data
	fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil || body == nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Println("response Body:", string(body))
}
