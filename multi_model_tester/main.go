package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	mc "multi_model_test/composer"
	pb "multi_model_test/inference"
	"os"
	"strconv"
	"strings"
	"time"

	trunks "github.com/straightdave/trunks/lib"

	// need the old version of proto module.
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

type Flags struct {
	//ModelName string
	// ModelList      string
	ModelArray  string
	NumPerModel int
	// ModelVersion   string
	// BatchSize      int
	URL            string
	Duration       int
	QPS            int
	WorkerPool     uint64
	ConnectionPool uint64
	Debug          bool
	// InputNum       int
}

// Ported from:
// https://github.com/triton-inference-server/client/blob/main/src/python/examples/simple_http_string_infer_client.py
// And
// https://github.com/triton-inference-server/client/blob/main/src/go/grpc_simple_client.go
func parseFlags() Flags {
	var flags Flags

	// flag.StringVar(&flags.ModelName, "m", "simple-string", "Name of model being served. (Required)")
	// flag.StringVar(&flags.ModelList, "ml", "", "A file containing model names. One model per line.")
	flag.StringVar(&flags.ModelArray, "ma", "", "List of different types models separate by space. Available options: SimpleStringTF, MnistSklearn, MushroomXgboot, CifarPytorch, MushroomLightgbm, MnistOnnx")
	flag.IntVar(&flags.NumPerModel, "npm", 1, "Number of model name 1 to npm per model to generate")
	// flag.StringVar(&flags.ModelVersion, "x", "", "Version of model. Default: Latest Version.")
	// flag.IntVar(&flags.BatchSize, "b", 1, "Batch size. Default: 1.")
	flag.StringVar(&flags.URL, "u", "localhost:8033", "Inference Server URL. Default: localhost:8033")
	flag.IntVar(&flags.Duration, "dur", 1, "Test duration in seconds. Default: 1")
	flag.IntVar(&flags.QPS, "qps", 1, "Constant Queries Per Second to hold. Default: 1")
	flag.Uint64Var(&flags.WorkerPool, "wp", 1, "Number of worker pool. Default: 1")
	flag.Uint64Var(&flags.ConnectionPool, "cp", 1, "Number of connections to create. Default: 1")
	flag.BoolVar(&flags.Debug, "debug", false, "Sends 1 request and view response. Default: false")
	// flag.IntVar(&flags.InputNum, "i", 1, "Number of unique simple-string model input tensors to generate for used in the infer request. Default: 1")
	flag.Parse()
	return flags
}

func sendReqest(URL *string, req *pb.ModelInferRequest) {
	// Directly send request to the gRPC server. Needs to import pb.go with grpc client interface
	// protoc -I path/to/dependency_protos -go_out=plugins=grpc=./ path/*.proto
	conn, err := grpc.Dial(*URL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to endpoint %s: %v", *URL, err)
	}
	defer conn.Close()
	// Create client from gRPC server connection
	client := pb.NewGRPCInferenceServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	modelInferResponse, err := client.ModelInfer(ctx, req)
	if err != nil {
		log.Fatalf("Error processing InferRequest: %v", err)
	}
	fmt.Println(modelInferResponse)
}
func genMultiModelRequest(ms string, n int) []proto.Message {
	ms = strings.TrimSpace(ms)
	a := strings.Split(ms, " ")
	var requests []proto.Message
	for _, mt := range a {
		switch mt {
		case "SimpleStringTF":
			ss := mc.SimpleStringTF{}
			for i := 1; i <= n; i++ {
				name := "simple-string-tf-" + strconv.Itoa(i)
				requests = append(requests, ss.GetInferRequest(name))
			}
		case "MnistSklearn":
			mskl := mc.MnistSklearn{}
			for i := 1; i <= n; i++ {
				name := "mnist-sklearn-" + strconv.Itoa(i)
				requests = append(requests, mskl.GetInferRequest(name))
			}
		case "MushroomXgboot":
			mxgb := mc.MushroomXgboot{}
			for i := 1; i <= n; i++ {
				name := "mushroom-xgboost-" + strconv.Itoa(i)
				requests = append(requests, mxgb.GetInferRequest(name))
			}
		case "CifarPytorch":
			cifarpyt := mc.CifarPytorch{}
			for i := 1; i <= n; i++ {
				name := "cifar-pytorch-" + strconv.Itoa(i)
				requests = append(requests, cifarpyt.GetInferRequest(name))
			}
		case "MushroomLightgbm":
			mlgbm := mc.MushroomLightgbm{}
			for i := 1; i <= n; i++ {
				name := "mushroom-lightgbm-" + strconv.Itoa(i)
				requests = append(requests, mlgbm.GetInferRequest(name))
			}

		case "MnistOnnx":
			monnx := mc.MnistOnnx{}
			for i := 1; i <= n; i++ {
				name := "mnist-onnx-" + strconv.Itoa(i)
				requests = append(requests, monnx.GetInferRequest(name))
			}
		}
	}
	return requests
}

func debugModelRequest(ms string, n int) []*pb.ModelInferRequest {
	ms = strings.TrimSpace(ms)
	a := strings.Split(ms, " ")
	var requests []*pb.ModelInferRequest
	for _, mt := range a {
		switch mt {
		case "SimpleStringTF":
			ss := mc.SimpleStringTF{}
			for i := 1; i <= n; i++ {
				name := "simple-string-tf-" + strconv.Itoa(i)
				requests = append(requests, ss.GetInferRequest(name))
			}
		case "MnistSklearn":
			mskl := mc.MnistSklearn{}
			for i := 1; i <= n; i++ {
				name := "mnist-sklearn-" + strconv.Itoa(i)
				requests = append(requests, mskl.GetInferRequest(name))
			}
		case "MushroomXgboot":
			mxgb := mc.MushroomXgboot{}
			for i := 1; i <= n; i++ {
				name := "mushroom-xgboost-" + strconv.Itoa(i)
				requests = append(requests, mxgb.GetInferRequest(name))
			}
		case "CifarPytorch":
			cifarpyt := mc.CifarPytorch{}
			for i := 1; i <= n; i++ {
				name := "cifar-pytorch-" + strconv.Itoa(i)
				requests = append(requests, cifarpyt.GetInferRequest(name))
			}
		case "MushroomLightgbm":
			mlgbm := mc.MushroomLightgbm{}
			for i := 1; i <= n; i++ {
				name := "mushroom-lightgbm-" + strconv.Itoa(i)
				requests = append(requests, mlgbm.GetInferRequest(name))
			}

		case "MnistOnnx":
			monnx := mc.MnistOnnx{}
			for i := 1; i <= n; i++ {
				name := "mnist-onnx-" + strconv.Itoa(i)
				requests = append(requests, monnx.GetInferRequest(name))
			}
		}
	}
	return requests
}

func newTarget(FLAGS *Flags) *trunks.Gtarget {
	var requests []proto.Message
	if FLAGS.ModelArray != "" {
		requests = genMultiModelRequest(FLAGS.ModelArray, FLAGS.NumPerModel)
		// fmt.Println(requests)
	}
	return &trunks.Gtarget{
		MethodName: "/inference.GRPCInferenceService/ModelInfer",
		Requests:   requests[:],
		Response:   &pb.ModelInferResponse{}, // providing one response data struct
	}
}

// func modelFileScanner(modelList string) (*bufio.Scanner, *os.File) {
// 	file, err := os.Open(modelList)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(file)
// 	return scanner, file
// }
func main() {
	FLAGS := parseFlags()

	if FLAGS.Debug {
		for _, req := range debugModelRequest(FLAGS.ModelArray, FLAGS.NumPerModel) {
			sendReqest(&FLAGS.URL, req)
		}
		os.Exit(0)
	}

	burner, err := trunks.NewBurner(
		[]string{FLAGS.URL},                    // server address pool; simple round-robin
		trunks.WithLooping(true),               // loop requests; false by default
		trunks.WithNumWorker(FLAGS.WorkerPool), // worker goroutine pool size; 10 is default
		trunks.WithDumpFile("result.dump"),     // dump responses to file
		trunks.WithNumConnPerHost(FLAGS.ConnectionPool),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer burner.Close()

	target := newTarget(&FLAGS)
	var metrics trunks.Metrics
	startT := time.Now()
	for res := range burner.Burn(target, uint64(FLAGS.QPS), time.Duration(FLAGS.Duration)*time.Second) {
		metrics.Add(res)
		// fmt.Println(res.Code, res.Latency, res.Error)
	}
	dur := time.Since(startT)
	metrics.Close()

	fmt.Printf("dur: %v\n", dur.Seconds())
	fmt.Printf("earliest: %v\n", metrics.Earliest.Sub(startT).Nanoseconds())
	fmt.Printf("latest: %v\n", metrics.Latest.Sub(startT).Nanoseconds())
	fmt.Printf("end: %v\n", metrics.End.Sub(startT).Nanoseconds())
	fmt.Printf("reqs: %d\n", metrics.Requests)
	fmt.Printf("success: %f\n", metrics.Success)
	fmt.Printf("p50: %s\n", metrics.Latencies.P50)
	fmt.Printf("p95: %s\n", metrics.Latencies.P95)
	fmt.Printf("p99: %s\n", metrics.Latencies.P99)
	fmt.Printf("mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("max: %s\n", metrics.Latencies.Max)

}
