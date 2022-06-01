package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	konstruktor "github.com/mhpixxio/konstruktor"
)

func main() {

	//---------------------------------- set the flags ----------------------------------
	http_url_flag := flag.String("http_url", "http://localhost:4040", "the address")
	filename_filetransfer_flag := flag.String("filename_filetransfer", "filetransfer_Star_Wars_Style_A_poster_1977.webp", "the name of the file for uploading and downloading")
	filename_streaming_flag := flag.String("filename_streaming", "chunkdata_smaller.zip", "the name of the file for streaming")
	size_bigdata_flag := flag.Int("size_bigdata", 100, "size of big data requests in megabytes (size when data gets encoded in grpc protobuf)")
	runs_flag := flag.Int("runs", 50, "number of runs")
	loops_flag := flag.Int("loops", 10, "number of repeated messages for small data before time measurement and taking average. Gives a more accurate result")
	amount_smalldata_flag := flag.Int("amount_smalldata", 100, "amount of small-data-messages for sending a lot of small messages simultaneously or after one another")
	only_size_measurement_flag := flag.Bool("only_size_measurement", false, "if true, skips the time measurements")
	random_data_measurement_flag := flag.Bool("random_data_measurement", true, "if false, skips the random data measurements")
	filetransfer_measurement_flag := flag.Bool("filetransfer_measurement", true, "if false, skips the file measurements")
	stream_measurement_flag := flag.Bool("stream_measurement", true, "if false, skips the stream measurements")
	flag.Parse()
	http_url := *http_url_flag
	filename_filetransfer := *filename_filetransfer_flag
	filename_streaming := *filename_streaming_flag
	size_bigdata := *size_bigdata_flag
	runs := *runs_flag
	loops := *loops_flag
	amount_smalldata := *amount_smalldata_flag
	only_size_measurement := *only_size_measurement_flag
	random_data_measurement := *random_data_measurement_flag
	filetransfer_measurement := *filetransfer_measurement_flag
	stream_measurement := *stream_measurement_flag
	log.Printf("http_url: %v, filename_filetransfer: %v, filename_streaming: %v, size_bigdata: %v, runs: %v, loops: %v, amount_smalldata: %v, only_size_measurement: %v, random_data_measurement: %v, filetransfer_measurement: %v, stream_measurement: %v", http_url, filename_filetransfer, filename_streaming, size_bigdata, runs, loops, amount_smalldata, only_size_measurement, random_data_measurement, filetransfer_measurement, stream_measurement)

	//---------------------------------- define variables to save benchmark results ----------------------------------
	benchmark_time_entries := 9
	benchmark_time := make([][]int, runs)
	for i := range benchmark_time {
		benchmark_time[i] = make([]int, benchmark_time_entries)
	}
	benchmark_size_entries := 14
	benchmark_size := make([][]int, runs)
	for i := range benchmark_size {
		benchmark_size[i] = make([]int, benchmark_size_entries)
	}

	//---------------------------------- start the runs ----------------------------------
	for k := 0; k < runs; k++ {

		//create small data
		smalldata := konstruktor.CreateBigData(1, 1)
		//create big data
		log.Printf("creating bigdata ...\n")
		var length_bigdata_float, slope float64
		var length_bigdata int
		slope = 291.8782939
		length_bigdata_float = math.Round((float64(size_bigdata)*1000000 - 4) / slope) //note: determined empirically
		length_bigdata = int(length_bigdata_float)
		bigdata := konstruktor.CreateBigData(7, length_bigdata)
		log.Printf("finished creating bigdata. server is ready.\n")

		log.Printf("starting benchmark run %v...\n", k)

		//---------------------------------- Size Measurements ----------------------------------
		if random_data_measurement == true {
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
		} else {
			benchmark_size[k][0] = 0
			benchmark_size[k][1] = 0
			benchmark_size[k][2] = 0
			benchmark_size[k][3] = 0
			benchmark_size[k][4] = 0
			benchmark_size[k][5] = 0
			benchmark_size[k][6] = 0
			benchmark_size[k][7] = 0
		}
		if filetransfer_measurement == true {
			//size measurements of file upload and download
			_, requestsize_upload, requestheadersize_upload, responsesize_upload, responseheadersize_upload := uploadclient(http_url, "/upload", "../httpclient/foruploadfiles/"+filename_filetransfer)
			_, responsesize_download, responseheadersize_download := downloadclient(http_url, "/download?filename="+filename_filetransfer, "../httpclient/downloadedfiles/"+filename_filetransfer)
			benchmark_size[k][8] = requestsize_upload
			benchmark_size[k][9] = requestheadersize_upload
			benchmark_size[k][10] = responsesize_upload
			benchmark_size[k][11] = responseheadersize_upload
			benchmark_size[k][12] = responsesize_download
			benchmark_size[k][13] = responseheadersize_download
		} else {
			benchmark_size[k][8] = 0
			benchmark_size[k][9] = 0
			benchmark_size[k][10] = 0
			benchmark_size[k][11] = 0
			benchmark_size[k][12] = 0
			benchmark_size[k][13] = 0
		}
		log.Printf("done with size measurement")

		//---------------------------------- Time Measurements ----------------------------------
		if only_size_measurement == false {
			if random_data_measurement == true {
				//Sending Big Data to Server
				start := time.Now()
				jsonclient(http_url, "/postjson", bigdata)
				elapsed := int(time.Since(start))
				benchmark_time[k][0] = int(elapsed)
				log.Printf("done with time measurement 0")
				//Receiving Big Data from Server
				start = time.Now()
				jsonclient(http_url, "/getjson", smalldata)
				elapsed = int(time.Since(start))
				benchmark_time[k][1] = int(elapsed)
				log.Printf("done with time measurement 1")
				//Sending Small Data to Server and Receiving Small Data
				start = time.Now()
				for i := 0; i < loops; i++ {
					jsonclient(http_url, "/postjson", smalldata)
				}
				elapsed = int(time.Since(start)) / loops
				benchmark_time[k][2] = int(elapsed)
				log.Printf("done with time measurement 2")
				//Sending a lot of Small Data to Server simultaneously
				start = time.Now()
				var wg sync.WaitGroup
				for i := 0; i < loops; i++ {
					wg.Add(amount_smalldata)
					for j := 0; j < amount_smalldata; j++ {
						go func() {
							jsonclient(http_url, "/postjson", smalldata)
							defer wg.Done()
						}()
					}
					wg.Wait()
				}
				elapsed = int(time.Since(start)) / loops
				benchmark_time[k][3] = int(elapsed)
				log.Printf("done with time measurement 3")
				//Sending a lot of Small Data to Server after one another
				start = time.Now()
				for i := 0; i < loops; i++ {
					for j := 0; j < amount_smalldata; j++ {
						jsonclient(http_url, "/postjson", smalldata)
					}
				}
				elapsed = int(time.Since(start)) / loops
				benchmark_time[k][4] = int(elapsed)
				log.Printf("done with time measurement 4")
			} else {
				benchmark_time[k][0] = 0
				benchmark_time[k][1] = 0
				benchmark_time[k][2] = 0
				benchmark_time[k][3] = 0
				benchmark_time[k][4] = 0
			}

			if filetransfer_measurement == true {
				//Upload a file to the server
				start := time.Now()
				uploadclient(http_url, "/upload", "../httpclient/foruploadfiles/"+filename_filetransfer)
				elapsed := int(time.Since(start))
				benchmark_time[k][5] = int(elapsed)
				log.Printf("done with time measurement 5")
				//Download a file from the server
				start = time.Now()
				downloadclient(http_url, "/download?filename="+filename_filetransfer, "../httpclient/downloadedfiles/"+filename_filetransfer)
				elapsed = int(time.Since(start))
				benchmark_time[k][6] = int(elapsed)
				log.Printf("done with time measurement 6")
			} else {
				benchmark_time[k][5] = 0
				benchmark_time[k][6] = 0
			}

			if stream_measurement == true {
				//Stream data to the server
				start := time.Now()
				clientsidestreamingclient(http_url, "/clientsidestreaming?filename="+filename_streaming, "../httpclient/foruploadfiles/"+filename_streaming)
				log.Printf("done with time measurement 7")
				elapsed := int(time.Since(start))
				benchmark_time[k][7] = int(elapsed)
				//Stream data from the server
				start = time.Now()
				serversidestreamingclient(http_url, "/serversidestreaming?filename="+filename_streaming, "../httpclient/downloadedfiles/"+filename_streaming)
				log.Printf("done with time measurement 8")
				elapsed = int(time.Since(start))
				benchmark_time[k][8] = int(elapsed)
			} else {
				benchmark_time[k][7] = 0
				benchmark_time[k][8] = 0
			}

		}
		log.Printf("done with benchmark run %v...\n", k)
	}

	//---------------------------------- print benchmark results to files ----------------------------------
	if only_size_measurement == false {
		file, err := os.OpenFile("../../results/benchmarking_http_time_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		datawriter := bufio.NewWriter(file)
		for k := 0; k < runs; k++ {
			for t := 0; t < benchmark_time_entries; t++ {
				_, _ = datawriter.WriteString(strconv.Itoa(k) + "\t" + strconv.Itoa(t) + "\t" + strconv.Itoa(benchmark_time[k][t]) + "\t")
			}
			_, _ = datawriter.WriteString("\n")
		}
		datawriter.Flush()
		file.Close()
	}
	file, err := os.OpenFile("../../results/benchmarking_http_size_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter := bufio.NewWriter(file)
	for k := 0; k < runs; k++ {
		for t := 0; t < benchmark_size_entries; t++ {
			_, _ = datawriter.WriteString(strconv.Itoa(k) + "\t" + strconv.Itoa(t) + "\t" + strconv.Itoa(benchmark_size[k][t]) + "\t")
		}
		_, _ = datawriter.WriteString("\n")
	}
	datawriter.Flush()
	file.Close()
}

//---------------------------------- client funcs ----------------------------------

func jsonclient(http_url string, endpoint string, data []konstruktor.RandomData) (body []byte, requestsize int, requestheadersize int, responsesize int, responseheadersize int) {
	//serialisation
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
	//save body in return
	resp_body, err := ioutil.ReadAll(response.Body)
	if err != nil || resp_body == nil {
		log.Fatalf("error: %v", err)
	}
	resp_requestsize := int(request.ContentLength) + len(request.Header)
	resp_requestheadersize := len(request.Header)
	resp_responsesize := len(resp_body) + len(response.Header)
	resp_responseheadersize := len(response.Header)
	return resp_body, resp_requestsize, resp_requestheadersize, resp_responsesize, resp_responseheadersize
}

func uploadclient(http_url string, endpoint string, filepath string) (body []byte, requestsize int, requestheadersize int, responsesize int, responseheadersize int) {
	r, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	values := map[string]io.Reader{
		"file": r,
	}
	//prepare a form to submit to that URL
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		//add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				log.Fatalf("error: %v", err)
				return
			}
		} else {
			//add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				log.Fatalf("error: %v", err)
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			log.Fatalf("error: %v", err)
			return
		}
	}
	//close the multipart writer
	w.Close()
	//submit the form to the handler
	req, err := http.NewRequest("POST", http_url+endpoint, &b)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	//set the content type, this will contain the boundary
	req.Header.Set("Content-Type", w.FormDataContentType())
	//submit the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	//check the response
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", response.Status)
	}
	//response
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer response.Body.Close()
	//save body in return
	resp_body, err := ioutil.ReadAll(response.Body)
	if err != nil || resp_body == nil {
		log.Fatalf("error: %v", err)
	}
	resp_requestsize := int(req.ContentLength) + len(req.Header)
	resp_requestheadersize := len(req.Header)
	resp_responsesize := len(resp_body) + len(response.Header)
	resp_responseheadersize := len(response.Header)
	return resp_body, resp_requestsize, resp_requestheadersize, resp_responsesize, resp_responseheadersize
}

func downloadclient(http_url string, endpoint string, filepath string) (body []byte, responsesize int, responseheadersize int) {
	//get the data
	resp, err := http.Get(http_url + endpoint)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()
	//create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer out.Close()
	//write the body to file
	resp_responsesize, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//save body in return
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp_body == nil {
		log.Fatalf("error: %v", err)
	}
	resp_responseheadersize := len(resp.Header)
	return resp_body, int(resp_responsesize), resp_responseheadersize
}

func clientsidestreamingclient(http_url string, endpoint string, filepath string) {
	//open file
	r, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//send request
	resp, err := http.Post(http_url+endpoint, "multipart/form-data", r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()
	return
}

func serversidestreamingclient(http_url string, endpoint string, filepath string) {
	//get the data
	resp, err := http.Get(http_url + endpoint)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()
	//create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer out.Close()
	//write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return
}
