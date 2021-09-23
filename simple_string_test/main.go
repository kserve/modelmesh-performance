package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	pb "simple_string_test/inference"
	"time"

	trunks "github.com/straightdave/trunks/lib"

	// need the old version of proto module.
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

const (
	inputSize  = 16
	outputSize = 16
)

type Flags struct {
	ModelName      string
	ModelList      string
	ModelVersion   string
	BatchSize      int
	URL            string
	Duration       int
	QPS            int
	WorkerPool     uint64
	ConnectionPool uint64
	Debug          bool
	InputNum       int
}

// Ported from:
// https://github.com/triton-inference-server/client/blob/main/src/python/examples/simple_http_string_infer_client.py
// And
// https://github.com/triton-inference-server/client/blob/main/src/go/grpc_simple_client.go
func parseFlags() Flags {
	var flags Flags

	flag.StringVar(&flags.ModelName, "m", "simple-string", "Name of model being served. (Required)")
	flag.StringVar(&flags.ModelList, "ml", "", "A file containing model names. One model per line.")
	flag.StringVar(&flags.ModelVersion, "x", "", "Version of model. Default: Latest Version.")
	flag.IntVar(&flags.BatchSize, "b", 1, "Batch size. Default: 1.")
	flag.StringVar(&flags.URL, "u", "localhost:8033", "Inference Server URL. Default: localhost:8033")
	flag.IntVar(&flags.Duration, "dur", 1, "Test duration in seconds. Default: 1")
	flag.IntVar(&flags.QPS, "qps", 1, "Constant Queries Per Second to hold. Default: 1")
	flag.Uint64Var(&flags.WorkerPool, "wp", 1, "Number of worker pool. Default: 1")
	flag.Uint64Var(&flags.ConnectionPool, "cp", 1, "Number of connections to create. Default: 1")
	flag.BoolVar(&flags.Debug, "debug", false, "Sends 1 request and view response. Default: false")
	flag.IntVar(&flags.InputNum, "i", 1, "Number of unique input tensors to generate for used in the infer request. Default: 1")
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

func randIntByteArrays() [][]byte {
	// array elements must be array of integers between 49 ~ 57
	//examples:
	//[[57 50 55] [52 52 50 56] [57 56 50 50] [50 48 56 55] [50 52 53 52] [51 55 49 57] [55 53 51 53] [57 57 49 56] [49 48 52 51] [55 57 57 49] [56 52 56 54] [50 57 48 55] [57 51 57 52] [51 55 49 49] [49 53 51 57] [55 55 48 54]]
	//[[57 51 52 53] [55 51 57 50] [50] [51] [52] [53] [54] [55] [56] [57] [49 48] [49 49] [49 50] [49 51] [49 52] [49 48 48 48 48 48 48]]
	arr := make([][]byte, inputSize)
	for i := 0; i < inputSize; i++ {
		rand.Seed(time.Now().UnixNano())
		b := fmt.Sprint(rand.Int31n(10000))
		arr[i] = []byte(b)
	}
	return arr
}

func getInferRequest(FLAGS *Flags, modelName string) *pb.ModelInferRequest {
	// The 'simple_string' model expects 2 BYTES tensors where each element in those tensors is the utf-8 string representation of an integer
	inputData0 := randIntByteArrays()
	inputData1 := randIntByteArrays()
	inferInputs := []*pb.ModelInferRequest_InferInputTensor{
		{
			Name:     "INPUT0",
			Datatype: "BYTES",
			Shape:    []int64{1, 16},
			Contents: &pb.InferTensorContents{
				BytesContents: inputData0,
			},
		},
		{
			Name:     "INPUT1",
			Datatype: "BYTES",
			Shape:    []int64{1, 16},
			Contents: &pb.InferTensorContents{
				BytesContents: inputData1,
			},
		},
	}
	// fmt.Println(modelName, inputData0)
	inferOutputs := []*pb.ModelInferRequest_InferRequestedOutputTensor{
		{
			Name: "OUTPUT0",
		},
		{
			Name: "OUTPUT1",
		},
	}

	inferReq := pb.ModelInferRequest{
		ModelName:    modelName,
		ModelVersion: FLAGS.ModelVersion,
		Inputs:       inferInputs,
		Outputs:      inferOutputs,
	}
	return &inferReq
}

func newTarget(FLAGS *Flags) *trunks.Gtarget {
	var requests []proto.Message
	if FLAGS.ModelList != "" {
		scanner, file := modelFileScanner(FLAGS.ModelList)
		defer file.Close() // ensure file closes when this function returns
		var modelArr []string
		for scanner.Scan() {
			modelArr = append(modelArr, scanner.Text())
		}
		size := len(modelArr)
		if FLAGS.InputNum < size {
			FLAGS.InputNum = size
		}
		requests = make([]proto.Message, FLAGS.InputNum)
		j := 0
		for i := 0; i < FLAGS.InputNum; i++ {
			if j == (size - 1) {
				j = 0
			}
			requests[i] = getInferRequest(FLAGS, modelArr[j])
			j++
		}
	} else {
		requests = make([]proto.Message, FLAGS.InputNum)
		for i := 0; i < FLAGS.InputNum; i++ {
			requests[i] = getInferRequest(FLAGS, FLAGS.ModelName)
		}
	}

	return &trunks.Gtarget{
		MethodName: "/inference.GRPCInferenceService/ModelInfer",
		Requests:   requests[:],
		Response:   &pb.ModelInferResponse{}, // providing one response data struct
	}
}

func modelFileScanner(modelList string) (*bufio.Scanner, *os.File) {
	file, err := os.Open(modelList)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	return scanner, file
}
func main() {
	FLAGS := parseFlags()

	if FLAGS.Debug {
		sendReqest(&FLAGS.URL, getInferRequest(&FLAGS, FLAGS.ModelName))
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
		fmt.Println(res.Code, res.Latency, res.Error)
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
