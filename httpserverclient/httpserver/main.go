package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mhpixxio/konstruktor"
)

var bigdata []konstruktor.RandomData
var smalldata []konstruktor.RandomData

func main() {

	//define endpoints
	http.HandleFunc("/hello", helloWorldHandler)
	http.HandleFunc("/postjson", postjsonHandler)
	http.HandleFunc("/getjson", getjsonHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	//create small data
	smalldata = konstruktor.CreateBigData(1, 1)
	//create big data
	log.Printf("creating bigdata ...\n")
	bigdata = konstruktor.CreateBigData(5, 100)
	log.Printf("finished creating bigdata. server is ready.\n")

	//start server
	fmt.Printf("starting server at port 4040\n")
	if err := http.ListenAndServe(":4040", nil); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func postjsonHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if body != nil {
		log.Println("received data")
		//Deserialisierung
		var data_req []konstruktor.RandomData
		err = json.Unmarshal(body, &data_req)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//fmt.Println(data_req) //print the data from the request
		//response
		jsons_res, err := json.Marshal(smalldata)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		w.Write(jsons_res)
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
		log.Println("received request")
		//Deserialisierung
		var data_req []konstruktor.RandomData
		err = json.Unmarshal(body, &data_req)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//fmt.Println(data_req) //print the data from the request
		//response
		jsons_res, err := json.Marshal(bigdata)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		w.Write(jsons_res)
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
		log.Fatalf("error: %v", err)
	}
	//close file again
	defer file.Close()
	// storage path
	f, err := os.OpenFile("../http-server/uploadedfiles/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		fmt.Fprintf(w, "upload successful")
	}
	// copy the file to the destination path
	io.Copy(f, file)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// get filename
	var filename string = r.FormValue("filename")
	// search for file
	fileBytes, err := ioutil.ReadFile("../http-server/uploadedfiles/" + filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// send file
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	fmt.Fprintf(w, "download successful")
}
