package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	konstruktor "github.com/mhpixxio/konstruktor"
	pb "github.com/mhpixxio/pb"
)

var res_bigdata = &pb.BigData{}
var res_smalldata = &pb.BigData{}
var size_bigdata int
var filename_clientsidestreaming string

type server_text struct {
	pb.TextServiceServer
}
type server_bigdata struct {
	pb.BigDataServiceServer
}
type server_upload struct {
	pb.UploadServiceServer
}
type server_download struct {
	pb.DownloadServiceServer
}
type server_serversidestreaming struct {
	pb.ServerSideStreamingServiceServer
}
type server_clientsidestreaming struct {
	pb.ClientSideStreamingServiceServer
}

func main() {

	//flags
	port_address_flag := flag.String("port_address", ":8080", "the port_address")
	size_bigdata_flag := flag.Int("size_bigdata", 100, "size of big data responses in megabytes (size when data gets encoded in grpc protobuf)")
	flag.Parse()
	port_address := *port_address_flag
	size_bigdata = *size_bigdata_flag
	log.Printf("port_address: %v, size_bigdata: %v", port_address, size_bigdata)

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
	var length_bigdata_float, slope float64
	var length_bigdata int
	slope = 291.8782939
	length_bigdata_float = math.Round((float64(size_bigdata)*1000000 - 4) / slope) //note: determined empirically
	length_bigdata = int(length_bigdata_float)
	bigdata_proto = konstruktor.CreateBigData_proto(7, length_bigdata)
	res_bigdata = &pb.BigData{Content: bigdata_proto}
	log.Printf("finished creating bigdata. server is ready.\n")
	//define calloptions
	max_size := size_bigdata * 1000000 * 2 * 100 //in bytes
	calloption_recv := grpc.MaxRecvMsgSize(max_size)
	calloption_send := grpc.MaxSendMsgSize(max_size)
	//start services
	s := grpc.NewServer(calloption_recv, calloption_send)
	pb.RegisterTextServiceServer(s, &server_text{})
	pb.RegisterBigDataServiceServer(s, &server_bigdata{})
	pb.RegisterUploadServiceServer(s, &server_upload{})
	pb.RegisterDownloadServiceServer(s, &server_download{})
	pb.RegisterServerSideStreamingServiceServer(s, &server_serversidestreaming{})
	pb.RegisterClientSideStreamingServiceServer(s, &server_clientsidestreaming{})
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

func (s *server_upload) UploadFunc(ctx context.Context, request *pb.UploadRequest) (*pb.UploadResponse, error) {
	//log.Printf("Received: %v %v", request.ProtoReflect().Descriptor().FullName(), request.Info) //print the name of the request
	err := ioutil.WriteFile("../grpcserver/uploadedfiles/"+request.Filename, request.Filebytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	res_bigdata_var := &pb.BigData{}
	res_bigdata_var = res_smalldata
	//return the response
	return &pb.UploadResponse{Bigdatares: res_bigdata_var}, nil
}

func (s *server_download) DownloadFunc(ctx context.Context, request *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	//log.Printf("Received: %v %v", request.ProtoReflect().Descriptor().FullName(), request.Info) //print the name of the request
	data, err := ioutil.ReadFile("../grpcserver/uploadedfiles/" + request.Filename)
	if err != nil {
		log.Fatal(err)
	}
	//return the response
	return &pb.DownloadResponse{Filebytes: data}, nil
}

func (s *server_serversidestreaming) ServerSideStreamingFunc(request *pb.StreamingRequestServerSide, stream pb.ServerSideStreamingService_ServerSideStreamingFuncServer) error {
	filename := request.Filename
	buffersize := request.Buffersize
	file, err := os.Open("../grpcserver/uploadedfiles/" + filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	buffer := make([]byte, buffersize)
	for {
		counter_buffer, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return err
			}
			break
		}
		if counter_buffer > 0 {
			stream.Send(&pb.Bytesmessage{Bytesmes: buffer})
		}
	}
	return err
}

func (s *server_clientsidestreaming) ClientSideStreamingFilenameFunc(ctx context.Context, request *pb.StreamingRequestClientSide) (*pb.Successmessage, error) {
	filename_clientsidestreaming = request.Filename
	return &pb.Successmessage{Successmes: true}, nil
}

func (s *server_clientsidestreaming) ClientSideStreamingFunc(stream pb.ClientSideStreamingService_ClientSideStreamingFuncServer) error {
	streaming_file_address := "../grpcserver/uploadedfiles/" + filename_clientsidestreaming
	f, err := os.Create(streaming_file_address)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer f.Close()
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error: %v", err)
			return err
		}
		// If the file doesn't exist, create it, or append to the file
		f_2, err := os.OpenFile(streaming_file_address, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
		if _, err := f_2.Write(data.Bytesmes); err != nil {
			log.Fatal(err)
			return err
		}
		if err := f_2.Close(); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return stream.SendAndClose(&pb.Successmessage{Successmes: true})
}
