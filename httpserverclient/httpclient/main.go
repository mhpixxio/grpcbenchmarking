package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mhpixxio/konstruktor"
)

func main() {

	//adress
	httpurl := "http://localhost:4040"
	//define loopsizes
	bigloops := 10
	loops := 10
	amountSmalldata := 100
	//define variable to save benchmark results
	benchmark_time := make([][]int, bigloops)
	for i := range benchmark_time {
		benchmark_time[i] = make([]int, 5) // 5 = the number of tests
	}

	for k := 0; k < bigloops; k++ {

		//create small data
		smalldata := konstruktor.CreateBigData(1, 1)
		//create big data
		log.Printf("creating bigdata ...\n")
		bigdata := konstruktor.CreateBigData(500, 100000)
		log.Printf("finished creating bigdata. server is ready.\n")

		log.Printf("starting benchmark run %v...\n", k)

		//Sending Big Data to Server
		start := time.Now()
		for i := 0; i < loops; i++ {
			jsonclient(httpurl, "/postjson", bigdata)
		}
		elapsed := int(time.Since(start)) / loops
		benchmark_time[k][0] = int(elapsed)
		log.Printf("done with test 0")

		//Receiving Big Data from Server
		start = time.Now()
		for i := 0; i < loops; i++ {
			jsonclient(httpurl, "/getjson", smalldata)
		}
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][1] = int(elapsed)
		log.Printf("done with test 1")

		//Sending Small Data to Server and Recieving Small Data
		start = time.Now()
		for i := 0; i < loops; i++ {
			jsonclient(httpurl, "/postjson", smalldata)
		}
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][2] = int(elapsed)
		log.Printf("done with test 2")

		//Sending a lot of Small Data to Server simultaniously
		start = time.Now()
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
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][3] = int(elapsed)
		log.Printf("done with test 3")

		//Sending a lot of Small Data to Server after one another
		start = time.Now()
		for i := 0; i < loops; i++ {
			for j := 0; j < amountSmalldata; j++ {
				jsonclient(httpurl, "/postjson", smalldata)
			}
		}
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][4] = int(elapsed)
		log.Printf("done with test 4")
		log.Printf("done with benchmark run %v...\n", k)
	}

	//print benchmark_time to a file
	file, err := os.OpenFile("benchmarking_http_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)
	for k := 0; k < bigloops; k++ {
		for t := 0; t < 5; t++ {
			_, _ = datawriter.WriteString(strconv.Itoa(k) + "\t" + strconv.Itoa(t) + "\t" + strconv.Itoa(benchmark_time[k][t]) + "\t")
		}
		_, _ = datawriter.WriteString("\n")
	}
	datawriter.Flush()
	file.Close()
}

func jsonclient(httpurl string, endpoint string, data []konstruktor.RandomData) (answer *http.Response) {
	//Serialisierung
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//connection
	httpposturl := httpurl + endpoint
	//log.Println("HTTP JSON POST URL:", httpposturl)
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
	//log.Println("sent data: ", string(jsonData)) //print the data
	//log.Println("response Status:", response.Status)
	//log.Println("response Headers:", response.Header)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil || body == nil {
		log.Fatalf("error: %v", err)
	}
	//log.Println("response Body:", string(body))
	return response
}
