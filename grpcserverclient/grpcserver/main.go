package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	konstruktor "github.com/mhpixxio/konstruktor"
	pb "github.com/mhpixxio/pb"
)

var res_bigdata = &pb.BigData{}
var res_smalldata = &pb.BigData{}

type server_text struct {
	pb.TextServiceServer
}
type server_bigdata struct {
	pb.BigDataServiceServer
}

func main() {

	//settings
	port_address := ":8080"
	size_bigdata := 354 //in megabytes (size when data gets encrpyted in grpc protobuf)

	//start server
	log.Printf("starting server at port" + port_address + "\n")
	listener, err := net.Listen("tcp", port_address)
	if err != nil {
		log.Fatalf("unable to listen: %v", err)
	}
	//create small data
	smalldata_proto := []*pb.RandomData{}
	smalldata_proto = konstruktor.CreateBigData_proto(1, 1)
	res_smalldata = &pb.BigData{Content: smalldata_proto}
	//create big data
	log.Printf("creating bigdata ...\n")
	bigdata_proto := []*pb.RandomData{}
	var length_bigdata int
	length_bigdata = (size_bigdata*1000000 - 17) / 3524 //notiz: empirisch ermittelt
	bigdata_proto = konstruktor.CreateBigData_proto(500, length_bigdata)
	res_bigdata = &pb.BigData{Content: bigdata_proto}
	log.Printf("finished creating bigdata. server is ready.\n")
	//define calloptions
	max_size := size_bigdata * 1000000 * 2 //in bytes
	calloption_recv := grpc.MaxRecvMsgSize(max_size)
	//start services
	s := grpc.NewServer(calloption_recv)
	pb.RegisterTextServiceServer(s, &server_text{})
	pb.RegisterBigDataServiceServer(s, &server_bigdata{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server_text) TextFunc(ctx context.Context, request *pb.TextRequest) (*pb.TextResponse, error) {
	//log.Printf("Received: %v %v", request.ProtoReflect().Descriptor().FullName(), request.Textreq) //print the name of the request and the string textreq included the request
	//return the response
	return &pb.TextResponse{Textres: "server: message received successfully. content: " + request.Textreq}, nil
}

func (s *server_bigdata) BigDataFunc(ctx context.Context, request *pb.BigDataRequest) (*pb.BigDataResponse, error) {
	//log.Printf("Received: %v %v", request.ProtoReflect().Descriptor().FullName(), request.Info) //print the name of the request
	//log.Println(request.Bigdatareq) //print the bigdata
	res_bigdata_var := &pb.BigData{}
	if request.Returnbigdata == true {
		res_bigdata_var = res_bigdata
	} else {
		res_bigdata_var = res_smalldata
	}
	//return the response
	return &pb.BigDataResponse{Bigdatares: res_bigdata_var}, nil
}
