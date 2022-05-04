package main

import (
	"context"
	"log"
	"sync"

	"google.golang.org/grpc"

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
	//create small data
	smalldata_proto := []*pb.RandomData{}
	smalldata_proto = konstruktor.CreateBigData_proto(1, 1)
	req_smalldata := &pb.BigData{Content: smalldata_proto}
	//create big data
	log.Printf("creating bigdata ...\n")
	bigdata_proto := []*pb.RandomData{}
	bigdata_proto = konstruktor.CreateBigData_proto(5, 100)
	req_bigdata := &pb.BigData{Content: bigdata_proto}
	log.Printf("finished creating bigdata. client is ready.\n")
	//define calloptions
	max_size := 704800012 //in bytes
	calloption_recv := grpc.MaxCallRecvMsgSize(max_size)
	//define loopsizes
	loops := 10
	amountSmalldata := 100

	//Connection Test
	//define message
	var req string
	req = "Connection Test"
	//call service
	TextResponse, err := client_text.TextFunc(context.Background(), &pb.TextRequest{Textreq: req})
	if err != nil {
		log.Fatalf("could not use service: %v", err)
	}
	//print response
	log.Println(TextResponse.Textres)

	//Sending Big Data to Server
	for i := 0; i < loops; i++ {
		//call service
		BigDataResponse, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "bigdata", Returnbigdata: false, Bigdatareq: req_bigdata}, calloption_recv)
		if err != nil {
			log.Fatalf("could not use service: %v", err)
		}
		//print response
		log.Println(BigDataResponse.Info)
		//log.Println(BigDataResponse.Bigdatares) //print the bigdata
	}

	//Receiving Big Data from Server
	for i := 0; i < loops; i++ {
		//call service
		BigDataResponse_2, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "bigdata", Returnbigdata: true, Bigdatareq: req_smalldata}, calloption_recv)
		if err != nil {
			log.Fatalf("could not use service: %v", err)
		}
		//print response
		log.Println(BigDataResponse_2.Info)
		//log.Println(BigDataResponse_2.Bigdatares) //print the bigdata
	}

	//Sending Small Data to Server and Recieving Small Data
	for i := 0; i < loops; i++ {
		//call service
		SmallDataResponse, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata}, calloption_recv)
		if err != nil {
			log.Fatalf("could not use service: %v", err)
		}
		//print response
		log.Println(SmallDataResponse.Info)
		//log.Println(SmallDataResponse.Bigdatares) //print the bigdata
	}

	//Sending a lot of Small Data to Server simultaniously
	var wg sync.WaitGroup
	for i := 0; i < loops; i++ {
		ch := make(chan *pb.BigDataResponse)
		wg.Add(amountSmalldata)
		for j := 0; j < amountSmalldata; j++ {
			go func() {
				//call service
				BigDataResponse, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata}, calloption_recv)
				if err != nil {
					log.Fatalf("could not use service: %v", err)
				}
				//print response
				//fmt.Println(BigDataResponse.Info)
				ch <- BigDataResponse
				//log.Println(BigDataResponse.Bigdatares) //print the bigdata
			}()
			log.Println((<-ch).Info)
		}
		wg.Wait()
	}

	//Sending a lot of Small Data to Server after one another
	for i := 0; i < loops; i++ {
		for j := 0; j < amountSmalldata; j++ {
			//call service
			SmallDataResponse, err := client_bigdata.BigDataFunc(context.Background(), &pb.BigDataRequest{Info: "smalldata", Returnbigdata: false, Bigdatareq: req_smalldata}, calloption_recv)
			if err != nil {
				log.Fatalf("could not use service: %v", err)
			}
			//print response
			log.Println(SmallDataResponse.Info)
			//log.Println(SmallDataResponse.Bigdatares) //print the bigdata
		}
	}
}
