package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	protoV2 "google.golang.org/protobuf/proto"

	"github.com/mhpixxio/konstruktor"
	pb "github.com/mhpixxio/pb"
)

var res_bigdata = &pb.BigData{}

type server_text struct {
	pb.TextServiceServer
}
type server_bigdata struct {
	pb.BigDataServiceServer
}

func main() {
	//start server
	port := ":8080"
	log.Printf("starting server at port" + port + "\n")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("unable to listen: %v", err)
	}
	//create big data
	log.Printf("creating bigdata ...\n")
	bigdata_proto := []*pb.RandomData{}
	bigdata_proto = konstruktor.CreateBigData_proto(500, 100000)
	res_bigdata = &pb.BigData{Content: bigdata_proto}
	log.Printf("finished creating bigdata. server is ready.\n")
	//define calloptions
	max_size := 704800012 //in bytes
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
	//print the name of the request and the string textreq included the request
	//log.Printf("Received: %v %v", request.ProtoReflect().Descriptor().FullName(), request.Textreq)
	//return the response
	return &pb.TextResponse{Textres: "server: message received successfully. content: " + request.Textreq}, nil
}

func (s *server_bigdata) BigDataFunc(ctx context.Context, request *pb.BigDataRequest) (*pb.BigDataResponse, error) {
	//print the name of the request
	//log.Printf("Received: %v %v", request.ProtoReflect().Descriptor().FullName(), request.Info)
	//log.Println(request.Bigdatareq) //print the bigdata
	//check if big data should be responded

	res_bigdata_var := &pb.BigData{}
	res_bigdata_var = nil
	if request.Returnbigdata == true {
		res_bigdata_var = res_bigdata
	}
	var sizeofrequestres_var int64
	if request.Returnsizeofrequest == true {
		sizeofrequestres_var = int64(protoV2.Size(request))
	} else {
		sizeofrequestres_var = 0
	}

	//return the response
	return &pb.BigDataResponse{Info: "server: successfully received request.", Bigdatares: res_bigdata_var, Sizeofrequestres: sizeofrequestres_var}, nil
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	//can be implemented later
}
