package main

import (
	"bufio"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	proto "google.golang.org/protobuf/proto"

	konstruktor "github.com/mhpixxio/konstruktor"
	pb "github.com/mhpixxio/pb"
)

func main() {

	//flags
	address_flag := flag.String("address", "localhost:8080", "the address")
	filename_flag := flag.String("filename", "Star_Wars_Style_A_poster_1977.webp", "the name of the file for uploading and downloading")
	size_bigdata_flag := flag.Int("size_bigdata", 100, "in megabytes (size when data gets encrpyted in grpc protobuf)")
	runs_flag := flag.Int("runs", 50, "number of runs")
	loops_flag := flag.Int("loops", 10, "number of repeated messages for small data before time measurement and taking average. Gives a more accurate result")
	amount_smalldata_flag := flag.Int("amount_smalldata", 100, "amount of small-data-messages for sending a lot of small messages simultaniously or after one another")
	only_size_measurement_flag := flag.Bool("only_size_measurement", false, "if true, skips the time measurements")
	random_data_measurement_flag := flag.Bool("random_data_measurement", true, "if false, skips the random data measurements")
	file_measurement_flag := flag.Bool("file_measurement", false, "if false, skips the file measurements")
	stream_measurement_flag := flag.Bool("stream_measurement", false, "if false, skips the stream measurements")
	flag.Parse()
	address := *address_flag
	filename := *filename_flag
	size_bigdata := *size_bigdata_flag
	runs := *runs_flag
	loops := *loops_flag
	amount_smalldata := *amount_smalldata_flag
	only_size_measurement := *only_size_measurement_flag
	random_data_measurement := *random_data_measurement_flag
	file_measurement := *file_measurement_flag
	stream_measurement := *stream_measurement_flag

	log.Printf("address: %v, size_bigdata: %v, runs: %v, loops: %v, amount_smalldata: %v, only_size_measurement: %v, random_data_measurement: %v, file_measurement: %v, stream_measurement: %v", address, size_bigdata, runs, loops, amount_smalldata, only_size_measurement, random_data_measurement, file_measurement, stream_measurement)

	//set dial
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//create the client_stubs
	client_text := pb.NewTextServiceClient(conn)
	client_bigdata := pb.NewBigDataServiceClient(conn)
	client_upload := pb.NewUploadServiceClient(conn)
	client_download := pb.NewDownloadServiceClient(conn)
	//define calloptions
	max_size := size_bigdata * 1000000 * 2 //in bytes
	calloption_recv := grpc.MaxCallRecvMsgSize(max_size)

	//define variables to save benchmark results
	benchmark_time_entries := 7
	benchmark_time := make([][]int, runs)
	for i := range benchmark_time {
		benchmark_time[i] = make([]int, benchmark_time_entries)
	}
	benchmark_size_entries := 8
	benchmark_size := make([][]int, runs)
	for i := range benchmark_size {
		benchmark_size[i] = make([]int, benchmark_size_entries)
	}

	//start the runs
	for k := 0; k < runs; k++ {

		//create small data
		smalldata_proto := []*pb.RandomData{}
		smalldata_proto = konstruktor.CreateBigData_proto(1, 1)
		req_smalldata := &pb.BigData{Content: smalldata_proto}
		//create big data
		log.Printf("creating bigdata ...\n")
		bigdata_proto := []*pb.RandomData{}
		var length_bigdata_float, slope float64
		var length_bigdata int
		slope = 291.8782939
		length_bigdata_float = math.Round((float64(size_bigdata)*1000000 - 4) / slope) //note: determined empirically
		length_bigdata = int(length_bigdata_float)
		bigdata_proto = konstruktor.CreateBigData_proto(7, length_bigdata)
		req_bigdata := &pb.BigData{Content: bigdata_proto}
		log.Printf("finished creating bigdata. client is ready.\n")

		//Connection Test
		var req string
		req = "Connection Test"
		responseTextFunc, err := client_text.TextFunc(context.Background(), &pb.TextRequest{Textreq: req})
		if err != nil || responseTextFunc == nil {
			log.Fatalf("could not use service: %v", err)
		}
		log.Println(responseTextFunc.Textres)

		log.Printf("starting benchmark run %v...\n", k)

		//Measuring the Sizes of transfered data
		if random_data_measurement == true {
			//call service with small data
			responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_smalldata, Returnbigdata: false}, calloption_recv)
			if err != nil || responseBigDataFunc == nil {
				log.Fatalf("could not use service: %v", err)
			}
			requestsize_small := proto.Size(&pb.BigDataRequest{Bigdatareq: req_smalldata, Returnbigdata: false})
			responsesize_small := proto.Size(responseBigDataFunc)
			//call service with big data
			responseBigDataFunc, err = client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_bigdata, Returnbigdata: true}, calloption_recv)
			if err != nil || responseBigDataFunc == nil {
				log.Fatalf("could not use service: %v", err)
			}
			requestsize_big := proto.Size(&pb.BigDataRequest{Bigdatareq: req_bigdata, Returnbigdata: true})
			responsesize_big := proto.Size(responseBigDataFunc)
			benchmark_size[k][0] = requestsize_small
			benchmark_size[k][1] = responsesize_small
			benchmark_size[k][2] = requestsize_big
			benchmark_size[k][3] = responsesize_big
		} else {
			benchmark_size[k][0] = 0
			benchmark_size[k][1] = 0
			benchmark_size[k][2] = 0
			benchmark_size[k][3] = 0
		}
		if file_measurement == true {
			//call upload
			data, err := ioutil.ReadFile("../grpcclient/foruploadfiles/" + filename)
			responseUpload, err := client_upload.UploadFunc(context.Background(), &pb.UploadRequest{Filebytes: data, Filename: filename}, calloption_recv)
			if err != nil || responseUpload == nil {
				log.Fatalf("could not use service: %v", err)
			}
			requestsize_upload := proto.Size(&pb.UploadRequest{Filebytes: data, Filename: filename})
			responsesize_upload := proto.Size(responseUpload)
			//call download
			responseDownload, err := client_download.DownloadFunc(context.Background(), &pb.DownloadRequest{Filename: filename}, calloption_recv)
			if err != nil || responseDownload == nil {
				log.Fatalf("could not use service: %v", err)
			}
			err = ioutil.WriteFile("../grpcclient/downloadedfiles/"+filename, responseDownload.Filebytes, 0644)
			if err != nil {
				log.Fatal(err)
			}
			requestsize_download := proto.Size(&pb.DownloadRequest{Filename: filename})
			responsesize_download := proto.Size(responseDownload)
			benchmark_size[k][4] = requestsize_upload
			benchmark_size[k][5] = responsesize_upload
			benchmark_size[k][6] = requestsize_download
			benchmark_size[k][7] = responsesize_download
		} else {
			benchmark_size[k][4] = 0
			benchmark_size[k][5] = 0
			benchmark_size[k][6] = 0
			benchmark_size[k][7] = 0
		}

		log.Printf("done with size measurement")

		//Time Measurements
		if only_size_measurement == false {
			//Random Data Measurements
			if random_data_measurement == true {
				//Sending Big Data to Server
				start := time.Now()
				responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_bigdata, Returnbigdata: false}, calloption_recv)
				if err != nil || responseBigDataFunc == nil {
					log.Fatalf("could not use service: %v", err)
				}
				elapsed := int(time.Since(start))
				benchmark_time[k][0] = int(elapsed)
				log.Printf("done with test 0")
				//Receiving Big Data from Server
				start = time.Now()
				responseBigDataFunc, err = client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_smalldata, Returnbigdata: true}, calloption_recv)
				if err != nil || responseBigDataFunc == nil {
					log.Fatalf("could not use service: %v", err)
				}
				elapsed = int(time.Since(start))
				benchmark_time[k][1] = int(elapsed)
				log.Printf("done with test 1")
				//Sending Small Data to Server and Receiving Small Data
				start = time.Now()
				for i := 0; i < loops; i++ {
					//call service
					responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_smalldata, Returnbigdata: false}, calloption_recv)
					if err != nil || responseBigDataFunc == nil {
						log.Fatalf("could not use service: %v", err)
					}
				}
				elapsed = int(time.Since(start)) / loops
				benchmark_time[k][2] = int(elapsed)
				log.Printf("done with test 2")
				//Sending a lot of Small Data to Server simultaniously
				var wg sync.WaitGroup
				start = time.Now()
				for i := 0; i < loops; i++ {
					ch := make(chan *pb.BigDataResponse)
					wg.Add(amount_smalldata)
					for j := 0; j < amount_smalldata; j++ {
						go func() {
							//call service
							responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_smalldata, Returnbigdata: false}, calloption_recv)
							if err != nil || responseBigDataFunc == nil {
								log.Fatalf("could not use service: %v", err)
							}
							ch <- responseBigDataFunc
							defer wg.Done()
						}()
						if (<-ch) == nil {
							log.Fatalf("did not receive response: %v", (<-ch))
						}
					}
					wg.Wait()
				}
				elapsed = int(time.Since(start)) / loops
				benchmark_time[k][3] = int(elapsed)
				log.Printf("done with test 3")
				//Sending a lot of Small Data to Server after one another
				start = time.Now()
				for i := 0; i < loops; i++ {
					for j := 0; j < amount_smalldata; j++ {
						//call service
						responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Bigdatareq: req_smalldata, Returnbigdata: false}, calloption_recv)
						if err != nil || responseBigDataFunc == nil {
							log.Fatalf("could not use service: %v", err)
						}
					}
				}
				elapsed = int(time.Since(start)) / loops
				benchmark_time[k][4] = int(elapsed)
				log.Printf("done with test 4")
			} else {
				benchmark_time[k][0] = 0
				benchmark_time[k][1] = 0
				benchmark_time[k][2] = 0
				benchmark_time[k][3] = 0
				benchmark_time[k][4] = 0
			}

			//Time Measurement of Files
			if file_measurement == true {
				//Upload a file to the server
				start := time.Now()
				data, err := ioutil.ReadFile("../grpcclient/foruploadfiles/" + filename)
				responseUpload, err := client_upload.UploadFunc(context.Background(), &pb.UploadRequest{Filebytes: data, Filename: filename}, calloption_recv)
				if err != nil || responseUpload == nil {
					log.Fatalf("could not use service: %v", err)
				}
				elapsed := int(time.Since(start))
				benchmark_time[k][5] = int(elapsed)
				log.Printf("done with test 5")
				//Download a file from the server
				start = time.Now()
				responseDownload, err := client_download.DownloadFunc(context.Background(), &pb.DownloadRequest{Filename: filename}, calloption_recv)
				if err != nil || responseDownload == nil {
					log.Fatalf("could not use service: %v", err)
				}
				err = ioutil.WriteFile("../grpcclient/downloadedfiles/"+filename, responseDownload.Filebytes, 0644)
				if err != nil {
					log.Fatal(err)
				}
				elapsed = int(time.Since(start))
				benchmark_time[k][6] = int(elapsed)
				log.Printf("done with test 6")
			} else {
				benchmark_time[k][5] = 0
				benchmark_time[k][6] = 0
			}
		}
		log.Printf("done with benchmark run %v...\n", k)
	}

	//print benchmark_time to a file
	if only_size_measurement == false {
		file, err := os.OpenFile("../../results/benchmarking_grpc_time_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	//print benchmark_size to a file
	file, err := os.OpenFile("../../results/benchmarking_grpc_size_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
