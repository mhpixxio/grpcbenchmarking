package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	protoV2 "google.golang.org/protobuf/proto"

	"github.com/mhpixxio/konstruktor"
	pb "github.com/mhpixxio/pb"
)

func main() {
	//set dial
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//create the client_stubs
	client_text := pb.NewTextServiceClient(conn)
	client_bigdata := pb.NewBigDataServiceClient(conn)
	//define calloptions
	max_size := 704800012 //in bytes
	calloption_recv := grpc.MaxCallRecvMsgSize(max_size)
	//define loopsizes
	runs := 10
	loops := 10 //amount of messages for one time measurement
	amountSmalldata := 100
	//define variables to save benchmark results
	benchmark_time := make([][]int, runs)
	for i := range benchmark_time {
		benchmark_time[i] = make([]int, 5) // 5 = the number of tests
	}
	benchmark_size := make([][]int, runs)
	for i := range benchmark_size {
		benchmark_size[i] = make([]int, 4) // 4 = the number of different measurements per run
	}

	for k := 0; k < runs; k++ {

		//create small data
		smalldata_proto := []*pb.RandomData{}
		smalldata_proto = konstruktor.CreateBigData_proto(1, 1)
		req_smalldata := &pb.BigData{Content: smalldata_proto}
		//create big data
		log.Printf("creating bigdata ...\n")
		bigdata_proto := []*pb.RandomData{}
		bigdata_proto = konstruktor.CreateBigData_proto(500, 100000)
		req_bigdata := &pb.BigData{Content: bigdata_proto}
		log.Printf("finished creating bigdata. client is ready.\n")

		//Connection Test
		//define message
		var req string
		req = "Connection Test"
		//call service
		responseTextFunc, err := client_text.TextFunc(context.Background(), &pb.TextRequest{Textreq: req})
		if err != nil || responseTextFunc == nil {
			log.Fatalf("could not use service: %v", err)
		}
		//print response
		log.Println(responseTextFunc.Textres)

		log.Printf("starting benchmark run %v...\n", k)

		//Measuring the Size of Big and Small Requests und Responses
		//call service with small data
		responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata, Returnsizeofrequest: true}, calloption_recv)
		if err != nil || responseBigDataFunc == nil {
			log.Fatalf("could not use service: %v", err)
		}
		requestsize_small := int(responseBigDataFunc.Sizeofrequestres)
		responsesize_small := protoV2.Size(responseBigDataFunc)
		//call service with big data
		responseBigDataFunc, err = client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "bigdata", Returnbigdata: true, Bigdatareq: req_bigdata, Returnsizeofrequest: true}, calloption_recv)
		if err != nil || responseBigDataFunc == nil {
			log.Fatalf("could not use service: %v", err)
		}
		requestsize_big := int(responseBigDataFunc.Sizeofrequestres)
		responsesize_big := protoV2.Size(responseBigDataFunc)
		benchmark_size[k][0] = requestsize_small
		benchmark_size[k][1] = responsesize_small
		benchmark_size[k][2] = requestsize_big
		benchmark_size[k][3] = responsesize_big
		log.Printf("done with size measurement")

		//Sending Big Data to Server
		start := time.Now()
		for i := 0; i < loops; i++ {
			//call service
			responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "bigdata", Returnbigdata: false, Bigdatareq: req_bigdata, Returnsizeofrequest: false}, calloption_recv)
			if err != nil || responseBigDataFunc == nil {
				log.Fatalf("could not use service: %v", err)
			}
			//print response
			//log.Println(responseBigDataFunc.Info)
			//log.Println(responseBigDataFunc.Bigdatares) //print the bigdata
		}
		elapsed := int(time.Since(start)) / loops
		benchmark_time[k][0] = int(elapsed)
		log.Printf("done with test 0")

		//Receiving Big Data from Server
		start = time.Now()
		for i := 0; i < loops; i++ {
			//call service
			responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "bigdata", Returnbigdata: true, Bigdatareq: req_smalldata, Returnsizeofrequest: false}, calloption_recv)
			if err != nil || responseBigDataFunc == nil {
				log.Fatalf("could not use service: %v", err)
			}
			//print response
			//log.Println(responseBigDataFunc.Info)
			//log.Println(responseBigDataFunc.Bigdatares) //print the bigdata
		}
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][1] = int(elapsed)
		log.Printf("done with test 1")

		//Sending Small Data to Server and Recieving Small Data
		start = time.Now()
		for i := 0; i < loops; i++ {
			//call service
			responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata, Returnsizeofrequest: false}, calloption_recv)
			if err != nil || responseBigDataFunc == nil {
				log.Fatalf("could not use service: %v", err)
			}
			//print response
			//log.Println(responseBigDataFunc.Info)
			//log.Println(responseBigDataFunc.Bigdatares) //print the bigdata
		}
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][2] = int(elapsed)
		log.Printf("done with test 2")

		//Sending a lot of Small Data to Server simultaniously
		var wg sync.WaitGroup
		start = time.Now()
		for i := 0; i < loops; i++ {
			ch := make(chan *pb.BigDataResponse)
			wg.Add(amountSmalldata)
			for j := 0; j < amountSmalldata; j++ {
				go func() {
					//call service
					responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata, Returnsizeofrequest: false}, calloption_recv)
					if err != nil || responseBigDataFunc == nil {
						log.Fatalf("could not use service: %v", err)
					}
					ch <- responseBigDataFunc
					//log.Println(responseBigDataFunc.Bigdatares) //print the bigdata
					defer wg.Done()
				}()
				if (<-ch) == nil {
					log.Fatalf("did not receive response: %v", (<-ch))
				}
				//print response
				//log.Println((<-ch).Info)
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
				//call service
				responseBigDataFunc, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata, Returnsizeofrequest: false}, calloption_recv)
				if err != nil || responseBigDataFunc == nil {
					log.Fatalf("could not use service: %v", err)
				}
				//print response
				//log.Println(responseBigDataFunc.Info)
				//log.Println(responseBigDataFunc.Bigdatares) //print the bigdata
			}
		}
		elapsed = int(time.Since(start)) / loops
		benchmark_time[k][4] = int(elapsed)
		log.Printf("done with test 4")
		log.Printf("done with benchmark run %v...\n", k)
	}

	//print benchmark_time to a file
	file, err := os.OpenFile("benchmarking_time_grpc_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	//print benchmark_size to a file
	file, err = os.OpenFile("benchmarking_size_grpc_"+strconv.Itoa(time.Now().Year())+time.Now().Month().String()+strconv.Itoa(time.Now().Day())+"_"+strconv.Itoa(time.Now().Hour())+"_"+strconv.Itoa(time.Now().Minute())+"_"+strconv.Itoa(time.Now().Second())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	datawriter = bufio.NewWriter(file)
	for k := 0; k < runs; k++ {
		for t := 0; t < 4; t++ {
			_, _ = datawriter.WriteString(strconv.Itoa(k) + "\t" + strconv.Itoa(t) + "\t" + strconv.Itoa(benchmark_size[k][t]) + "\t")
		}
		_, _ = datawriter.WriteString("\n")
	}
	datawriter.Flush()
	file.Close()
}
