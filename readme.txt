+++grpc benchmarking+++

Comparison between a grpc server & client using protobuf protocols and a http server & client using
jsons for data transfer.
The clients run a number of data transfer scenarios and give out two .txt-files with all the time and 
size measurements. The content of these .txt-files can be copied and pasted into the provided excel 
sheet into the tabs "... raw data". The benchmarking analysis in the "Benchmarking"-tab gets then 
updated automatically.
The grpc server and client both need the same packages providing the specific protobuf message and 
service types used in this application. These had been defined previously and are provided by the 
import of
"github.com/mhpixxio/pb".


+++tests+++

The tests use big and small data packages (see chapter "the data" for futher explanation). The sizes 
and amounts can be changed (see chapter "flags"). With default flags the client runs all tests 50 times.
1. Size Measurements
  1.1 Size of a request with small data
  1.2 Size of a response with small data
  1.3 Size of a request with big data
  1.4 Size of a response with big data
  (Additionally for the http tranfers, the size of the headers always gets measured)
2. Time Measurements
  2.1 Sending Big Data to Server
  2.2 Receiving Big Data from Server
  2.3 Sending Small Data to Server and Recieving Small Data
  2.4 Sending a lot of Small Data to Server simultaniously
  2.5 Sending a lot of Small Data to Server after one another


+++flags+++

grpc server flags:
    -port_address string
        the port_address (default ":8080")
    -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 354)
gprc client flags:
    -address string
        the address (default "localhost:8080")
    -amountSmalldata int
        amount of small-data-messages for sending a lot of small messages simultaniously or after one another (default 100)
    -loops int
        number of repeated messages before time measurement and taking average. Gives a more accurate result (default 10)
    -only_size_measurement
        if true, skips the time measurments
    -runs int
        number of runs (default 50)
    -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 354)
http server flags:
    -port_address string
        the port_address (default ":4040")
    -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 354)
http client flags:
    -amountSmalldata int
        amount of small-data-messages for sending a lot of small messages simultaniously or after one another (default 100)
    -http_url string
        the address (default "http://localhost:4040")
    -loops int
        number of repeated messages before time measurement and taking average. Gives a more accurate result (default 10)
    -only_size_measurement
        if true, skips the time measurments
    -runs int
        number of runs (default 50)
    -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 354)


+++the data+++

Each server&client uses small and big data packages in the same data format to test the various scenarios.
The data format is an array holding X entries of the type RandomData struct.

type RandomData struct {
		a string
		b string
		c string
		d string
		e string
		f string
		g string
	}

Each string consists of Y random characters. The package konstruktor "github.com/mhpixxio/konstruktor" provides
the method konstruktor.CreateBigData(X, Y) and konstruktor.CreateBigData_proto(X, Y) which create the data used
in the tests. The empirically determined function "Y = (size_bigdata*1000000 - 17) / 3524" gives Y for size_bigdata
given in megabytes and X set to 500.
The small data used in the tests is created with X=1 and Y=1.

