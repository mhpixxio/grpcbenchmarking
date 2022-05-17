package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	konstruktor "github.com/mhpixxio/konstruktor"
)

var bigdata []konstruktor.RandomData
var smalldata []konstruktor.RandomData

func main() {

	//flags
	port_address_flag := flag.String("port_address", ":4040", "the port_address")
	size_bigdata_flag := flag.Int("size_bigdata", 354, "in megabytes (size when data gets encrpyted in grpc protobuf)")
	flag.Parse()
	port_address := *port_address_flag
	size_bigdata := *size_bigdata_flag
	log.Printf("port_address: %v, size_bigdata: %v", port_address, size_bigdata)

	//define endpoints
	http.HandleFunc("/connectiontest", connectiontestHandler)
	http.HandleFunc("/postjson", postjsonHandler)
	http.HandleFunc("/getjson", getjsonHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	//create small data
	smalldata = konstruktor.CreateBigData(1, 1)
	//create big data
	log.Printf("creating bigdata ...\n")
	var length_bigdata int
	length_bigdata = (size_bigdata*1000000 - 17) / 3524 //note: determined empirically
	bigdata = konstruktor.CreateBigData(500, length_bigdata)
	log.Printf("finished creating bigdata. server is ready.\n")

	//start server
	fmt.Printf("starting server at port" + port_address + "\n")
	if err := http.ListenAndServe(port_address, nil); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func connectiontestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! Connection successful!")
}

func postjsonHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if body != nil {
		//log.Println("received data")
		//deserialisation
		var data_req []konstruktor.RandomData
		err = json.Unmarshal(body, &data_req)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//fmt.Println(data_req) //print the data from the request
		//response
		json_res, err := json.Marshal(smalldata)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		w.Write(json_res)
	} else {
		log.Fatalf("did not receive data")
	}
}

func getjsonHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if body != nil {
		//log.Println("received data")
		//deserialisation
		var data_req []konstruktor.RandomData
		err = json.Unmarshal(body, &data_req)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//fmt.Println(data_req) //print the data from the request
		//response
		json_res, err := json.Marshal(bigdata)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		w.Write(json_res)
	} else {
		log.Fatalf("did not receive data")
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// limit size
	r.ParseMultipartForm(80 << 9)
	// retrieve file from FormFile
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "could not retrieve file")
		panic(err)
	}
	//close file again
	defer file.Close()
	// storage path
	f, err := os.OpenFile("../httpserver/uploadedfiles/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(w, "could not save file")
		panic(err)
	} else {
		fmt.Fprintf(w, "upload successful")
	}
	// copy the file to the destination path
	io.Copy(f, file)
	//response
	json_res, err := json.Marshal(smalldata)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	w.Write(json_res)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// get filename
	var filename string = r.FormValue("filename")
	// search for file
	fileBytes, err := ioutil.ReadFile("../httpserver/uploadedfiles/" + filename)
	if err != nil {
		fmt.Fprintf(w, "could not find file")
		panic(err)
	}
	// send file
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	fmt.Fprintf(w, "download successful")
}
