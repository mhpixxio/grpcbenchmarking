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

	konstruktor "github.com/mhpixxio/konstruktor"
)

func main() {

	//settings
	http_url := "http://localhost:4040"
	size_bigdata := 354 //in megabytes (size when same data gets encrpyted in grpc protobuf)
	runs := 10
	loops := 10 //amount of messages for one time measurement
	amountSmalldata := 100
	only_size_measurement := false

	//define variables to save benchmark results
	benchmark_time := make([][]int, runs)
	for i := range benchmark_time {
		benchmark_time[i] = make([]int, 5) // 5 = the number of tests
	}
	benchmark_size := make([][]int, runs)
	for i := range benchmark_size {
		benchmark_size[i] = make([]int, 8) // 8 = the number of different measurements per run
	}

	for k := 0; k < runs; k++ {

		//create small data
		smalldata := konstruktor.CreateBigData(1, 1)
		//create big data
		log.Printf("creating bigdata ...\n")
		var length_bigdata int
		length_bigdata = (size_bigdata*1000000 - 17) / 3524 //notiz: empirisch ermittelt
		bigdata := konstruktor.CreateBigData(500, length_bigdata)
		log.Printf("finished creating bigdata. server is ready.\n")

		log.Printf("starting benchmark run %v...\n", k)

		//Measuring the Size of Small and Big Requests und Responses
		_, requestsize_small, requestheadersize_small, responsesize_small, responseheadersize_small := jsonclient(http_url, "/postjson", smalldata)
		_, requestsize_big, requestheadersize_big, responsesize_big, responseheadersize_big := jsonclient(http_url, "/getjson", bigdata)
		benchmark_size[k][0] = requestsize_small
		benchmark_size[k][1] = requestheadersize_small
		benchmark_size[k][2] = responsesize_small
		benchmark_size[k][3] = responseheadersize_small
		benchmark_size[k][4] = requestsize_big
		benchmark_size[k][5] = requestheadersize_big
		benchmark_size[k][6] = responsesize_big
		benchmark_size[k][7] = responseheadersize_big
		log.Printf("done with size measurement")

		if only_size_measurement == false {
			//Sending Big Data to Server
			start := time.Now()
			for i := 0; i < loops; i++ {
				jsonclient(http_url, "/postjson", bigdata)
			}
			elapsed := int(time.Since(start)) / loops
			benchmark_time[k][0] = int(elapsed)
			log.Printf("done with test 0")

			//Receiving Big Data from Server
			start = time.Now()
			for i := 0; i < loops; i++ {
				jsonclient(http_url, "/getjson", smalldata)
			}
			elapsed = int(time.Since(start)) / loops
			benchmark_time[k][1] = int(elapsed)
			log.Printf("done with test 1")

			//Sending Small Data to Server and Recieving Small Data
			start = time.Now()
			for i := 0; i < loops; i++ {
				jsonclient(http_url, "/postjson", smalldata)
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
						jsonclient(http_url, "/postjson", smalldata)
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
					jsonclient(http_url, "/postjson", smalldata)
				}
			}
			elapsed = int(time.Since(start)) / loops
			benchmark_time[k][4] = int(elapsed)
			log.Printf("done with test 4")
		}
		log.Printf("done with benchmark run %v...\n", k)
	}

	if only_size_measurement == false {
		//print benchmark_time to a file
		file, err := os.OpenFile("benchmarking_time_http_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter := bufio.NewWriter(file)
		for k := 0; k < runs; k++ {
			for t := 0; t < 5; t++ {
				_, _ = datawriter.WriteString(strconv.Itoa(k) + "\t" + strconv.Itoa(t) + "\t" + strconv.Itoa(benchmark_time[k][t]) + "\t")
			}
			_, _ = datawriter.WriteString("\n")
		}
		datawriter.Flush()
		file.Close()
	}

	//print benchmark_size to a file
	file, err := os.OpenFile("benchmarking_size_http_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)
	for k := 0; k < runs; k++ {
		for t := 0; t < 8; t++ {
			_, _ = datawriter.WriteString(strconv.Itoa(k) + "\t" + strconv.Itoa(t) + "\t" + strconv.Itoa(benchmark_size[k][t]) + "\t")
		}
		_, _ = datawriter.WriteString("\n")
	}
	datawriter.Flush()
	file.Close()
}

func jsonclient(http_url string, endpoint string, data []konstruktor.RandomData) (body []byte, requestsize int, requestheadersize int, responsesize int, responseheadersize int) {
	//Serialisierung
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//connection
	httpposturl := http_url + endpoint
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
	resp_body, err := ioutil.ReadAll(response.Body)
	if err != nil || resp_body == nil {
		log.Fatalf("error: %v", err)
	}
	//log.Println("response Body:", string(resp_body))
	resp_requestsize := int(request.ContentLength) + len(request.Header)
	resp_requestheadersize := len(request.Header)
	resp_responsesize := len(resp_body) + len(response.Header)
	resp_responseheadersize := len(response.Header)
	return resp_body, resp_requestsize, resp_requestheadersize, resp_responsesize, resp_responseheadersize
}
