+++grpc benchmarking+++

Comparison between a grpc server & client using protobuf protocols for data transfer and a http server & client using
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

The first four tests use big and small data packages (see chapter "the data" for futher explanation). The sizes 
and amounts can be changed (see chapter "flags"). The flag "size_bigdata" refers to the size of the data when encoded as
protobuf data. If both servers/clients are given size_bigdata=100, then the grpc server/client will send 100MB of encoded
data, and the http server/client will use the same data but encoded as json (which is always larger as the tests have shown).
With default flags the client runs all tests 50 times.

1. Size measurements
  1.1 Size of a request with small data
  1.2 Size of a response with small data
  1.3 Size of a request with big data
  1.4 Size of a response with big data
  1.5 Size of file upload to server
  1.6 Size of file download from server
  (Additionally for the http tranfers, the size of the headers always gets measured)
2. Time measurements
  2.1 Sending big data to server
  2.2 Receiving big data from server
  2.3 Sending small data to server and receiving small data
  2.4 Sending a lot of small data to server simultaniously
  2.5 Sending a lot of small data to server after one another
  2.6 Uploading a file to the server
  2.7 Downloading a file from the server


+++flags+++

grpc server flags:
 -port_address string
        the port_address (default ":8080")
  -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 100)
gprc client flags:
  -address string
        the address (default "localhost:8080")
  -amount_smalldata int
        amount of small-data-messages for sending a lot of small messages simultaniously or after one another (default 100)
  -buffersize_streaming int
        buffersize for streaming (default 100)
  -filename_filetransfer string
        the name of the file for uploading and downloading (default "filetransfer_Star_Wars_Style_A_poster_1977.webp")
  -filename_streaming string
        the name of the file for streaming (default "streaming_Star_Wars_Style_A_poster_1977.webp")
  -filetransfer_measurement
        if false, skips the file measurements (default false)
  -loops int
        number of repeated messages for small data before time measurement and taking average. Gives a more 
accurate result (default 10)
  -only_size_measurement
        if true, skips the time measurements (default false)
  -random_data_measurement
        if false, skips the random data measurements (default true)
  -runs int
        number of runs (default 50)
  -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 100)
  -stream_measurement
        if false, skips the stream measurements (default false)
http server flags:
  -port_address string
        the port_address (default ":4040")
  -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 100)
http client flags:
  -amount_smalldata int
        amount of small-data-messages for sending a lot of small messages simultaniously or after one another (default 100)
  -file_measurement
        if false, skips the file measurements (default false)
  -filename string
        the name of the file for uploading and downloading (default "Star_Wars_Style_A_poster_1977.webp")  
  -http_url string
        the address (default "http://localhost:4040")
  -loops int
        number of repeated messages for small data before time measurement and taking average. Gives a more accurate result (default 10)
  -only_size_measurement
        if true, skips the time measurments (default false)
  -random_data_measurement
        if false, skips the random data measurements (default true)
  -runs int
        number of runs (default 50)
  -size_bigdata int
        in megabytes (size when data gets encrpyted in grpc protobuf) (default 100)
  -stream_measurement
        if false, skips the stream measurements (default false)


+++the data+++

Each server&client uses small and big data packages in the same data format to test the various scenarios.
The data format is an array holding X entries of the type RandomData struct.

type RandomData struct {
	Aaabcdef int32  `json:"aaabcdef,omitempty"`
	Ababcdef bool   `json:"ababcdef,omitempty"`
	Acabcdef bool   `json:"acabcdef,omitempty"`
	Adabcdef bool   `json:"adabcdef,omitempty"`
	Aeabcdef bool   `json:"aeabcdef,omitempty"`
	Afabcdef bool   `json:"afabcdef,omitempty"`
	Agabcdef bool   `json:"agabcdef,omitempty"`
	Ahabcdef bool   `json:"ahabcdef,omitempty"`
	Aiabcdef bool   `json:"aiabcdef,omitempty"`
	Ajabcdef bool   `json:"ajabcdef,omitempty"`
	Akabcdef bool   `json:"akabcdef,omitempty"`
	Alabcdef string `json:"alabcdef,omitempty"`
	Amabcdef string `json:"amabcdef,omitempty"`
	Anabcdef string `json:"anabcdef,omitempty"`
	Aoabcdef string `json:"aoabcdef,omitempty"`
	Apabcdef string `json:"apabcdef,omitempty"`
	Aqabcdef string `json:"aqabcdef,omitempty"`
	Arabcdef string `json:"arabcdef,omitempty"`
	Asabcdef string `json:"asabcdef,omitempty"`
	Atabcdef string `json:"atabcdef,omitempty"`
	Auabcdef string `json:"auabcdef,omitempty"`
	Avabcdef string `json:"avabcdef,omitempty"`
	Awabcdef string `json:"awabcdef,omitempty"`
	Axabcdef string `json:"axabcdef,omitempty"`
	Ayabcdef string `json:"ayabcdef,omitempty"`
	Azabcdef string `json:"azabcdef,omitempty"`
	Baabcdef string `json:"baabcdef,omitempty"`
	Bbabcdef string `json:"bbabcdef,omitempty"`
	Bcabcdef string `json:"bcabcdef,omitempty"`
	Bdabcdef string `json:"bdabcdef,omitempty"`
	Beabcdef string `json:"beabcdef,omitempty"`
	Bfabcdef string `json:"bfabcdef,omitempty"`
	Bgabcdef string `json:"bgabcdef,omitempty"`
	Bhabcdef string `json:"bhabcdef,omitempty"`
	Biabcdef string `json:"biabcdef,omitempty"`
	Bjabcdef string `json:"bjabcdef,omitempty"`
	Bkabcdef string `json:"bkabcdef,omitempty"`
	Blabcdef string `json:"blabcdef,omitempty"`
	Bmabcdef string `json:"bmabcdef,omitempty"`
	Bnabcdef string `json:"bnabcdef,omitempty"`
	Boabcdef string `json:"boabcdef,omitempty"`
	Bpabcdef string `json:"bpabcdef,omitempty"`
	Bqabcdef string `json:"bqabcdef,omitempty"`
	Brabcdef string `json:"brabcdef,omitempty"`
	Bsabcdef string `json:"bsabcdef,omitempty"`
	Btabcdef string `json:"btabcdef,omitempty"`
	Buabcdef string `json:"buabcdef,omitempty"`
	Bvabcdef string `json:"bvabcdef,omitempty"`
	Bwabcdef string `json:"bwabcdef,omitempty"`
	Bxabcdef string `json:"bxabcdef,omitempty"`
	}

Real data jsons were inspected to derive a statistically correct distribution of data types and filling of the entries.
2.56% int (filled randomly)
4.48% true bools
15.47% false bools
21.83% empty strings
55.65% strings with an average of 7 characters (filled randomly)

The package konstruktor "github.com/mhpixxio/konstruktor" provides
the method konstruktor.CreateBigData(lengthString, X) and konstruktor.CreateBigData_proto(lengthString, X) which create the data used
in the tests. lengthString is always set to 7 for big data. The data gets newly created for every run.

The empirically determined function "X=(size_bigdata)*1000000-4)/291.8782939)" gives X for size_bigdata
given in megabytes and lengthString set to 7. size_bigdata refers to the size the wired protobuf data has when fully encoded.

The small data used in the tests is created with X=1 and lengthString=1.

