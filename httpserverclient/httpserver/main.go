package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	konstruktor "github.com/mhpixxio/konstruktor"
)

var bigdata []konstruktor.RandomData
var smalldata []konstruktor.RandomData

func main() {

	//flags
	port_address_flag := flag.String("port_address", ":4040", "the port_address")
	size_bigdata_flag := flag.Int("size_bigdata", 100, "in megabytes (size when data gets encrpyted in grpc protobuf)")
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
	http.HandleFunc("/clientsidestreaming", clientsidestreamingHandler)
	http.HandleFunc("/serversidestreaming", serversidestreamingHandler)

	//create small data
	smalldata = konstruktor.CreateBigData(1, 1)
	//create big data
	log.Printf("creating bigdata ...\n")
	var length_bigdata_float, slope float64
	var length_bigdata int
	slope = 291.8782939
	length_bigdata_float = math.Round((float64(size_bigdata)*1000000 - 4) / slope) //note: determined empirically
	length_bigdata = int(length_bigdata_float)
	bigdata = konstruktor.CreateBigData(7, length_bigdata)
	log.Printf("finished creating bigdata. server is ready.\n")

	//start server
	fmt.Printf("starting server at port" + port_address + "\n")
	if err := http.ListenAndServe(port_address, nil); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func connectiontestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusOK)
		w.Write(json_res)
	} else {
		log.Fatalf("did not receive data")
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
	}
	// copy the file to the destination path
	io.Copy(f, file)
	//response
	json_res, err := json.Marshal(smalldata)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json_res)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	//get filename
	var filename string = r.FormValue("filename")
	//search for file and read whole file
	fileBytes, err := ioutil.ReadFile("../httpserver/uploadedfiles/" + filename)
	if err != nil {
		fmt.Fprintf(w, "could not find file")
		panic(err)
	}
	//send file
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

func clientsidestreamingHandler(w http.ResponseWriter, r *http.Request) {
	//get filename
	var filename string = r.FormValue("filename")
	//create file
	out, err := os.Create("../httpserver/uploadedfiles/" + filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer out.Close()
	w.WriteHeader(http.StatusOK)
	//write to file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func serversidestreamingHandler(w http.ResponseWriter, r *http.Request) {
	//get filename
	var filename string = r.FormValue("filename")
	//open file
	file, err := os.Open("../httpserver/uploadedfiles/" + filename)
	if err != nil {
		fmt.Fprintf(w, "could not find file")
		log.Println(err)
	}
	defer file.Close()
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.WriteHeader(http.StatusOK)
	//send file as stream
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "could not read body", http.StatusInternalServerError)
	}
}
