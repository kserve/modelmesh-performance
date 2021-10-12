package composer

import (
	"fmt"
	"log"
	"math/rand"
	"multi_model_test/inference"
	pb "multi_model_test/inference"
	"os"
	"time"

	"github.com/dereklstinson/cifar"
	"github.com/moverest/mnist"
)

const (
	inputSize  = 16
	outputSize = 16
)

type ModelRequestComposer interface {
	getInferRequest() pb.ModelInferRequest
}

type SimpleStringTF struct {
	// tensorflow model
}

type MnistSklearn struct {
	//sklearn model
}

type MushroomXgboot struct {
	//xgboost model
}

type CifarPytorch struct {
	//pytorch model
	Image []float32
}

type MushroomLightgbm struct {
	//lightgbm model
}

type MnistOnnx struct {
	//MnistOnxx model
	Image []float32
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
func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func LoadCifarImage(index int) []float32 {
	file, err := os.Open("testdata/cifar_test_images.bin")
	exitOnError(err)
	images, err := cifar.Decode10(file)
	exitOnError(err)

	imageBytes := images[index].RawData()
	var imageFloat [3 * 32 * 32]float32
	for i, v := range imageBytes {
		// the test PyTorch CIFAR model was trained based on:
		// - https://github.com/kubeflow/kfserving/tree/master/docs/samples/v1alpha2/pytorch
		// - https://pytorch.org/tutorials/beginner/blitz/cifar10_tutorial.html
		// These models are trained on images with pixels normalized to the range
		// [-1 1]. The testdata contains images with pixels in bytes [0 255] that
		// must be transformed
		imageFloat[i] = (float32(v) / 127.5) - 1
	}

	return imageFloat[:]
}

func LoadMnistImage(index int) []float32 {
	images, err := mnist.LoadImageFile("testdata/t10k-images-idx3-ubyte.gz")
	exitOnError(err)

	imageBytes := [mnist.Width * mnist.Height]byte(*images[index])
	var imageFloat [mnist.Width * mnist.Height]float32
	for i, v := range imageBytes {
		imageFloat[i] = float32(v)
	}
	return imageFloat[:]
}

func (m *SimpleStringTF) GetInferRequest(modelName string) *pb.ModelInferRequest {
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

	inferOutputs := []*pb.ModelInferRequest_InferRequestedOutputTensor{}

	inferReq := pb.ModelInferRequest{
		ModelName: modelName,
		Inputs:    inferInputs,
		Outputs:   inferOutputs,
	}
	return &inferReq
}

func (m *MnistSklearn) GetInferRequest(modelName string) *pb.ModelInferRequest {
	// { "inputs": [{ "name": "predict", "shape": [1, 64], "datatype": "FP32", "data": [0.0, 0.0, 1.0, 11.0, 14.0, 15.0, 3.0, 0.0, 0.0, 1.0, 13.0, 16.0, 12.0, 16.0, 8.0, 0.0, 0.0, 8.0, 16.0, 4.0, 6.0, 16.0, 5.0, 0.0, 0.0, 5.0, 15.0, 11.0, 13.0, 14.0, 0.0, 0.0, 0.0, 0.0, 2.0, 12.0, 16.0, 13.0, 0.0, 0.0, 0.0, 0.0, 0.0, 13.0, 16.0, 16.0, 6.0, 0.0, 0.0, 0.0, 0.0, 16.0, 16.0, 16.0, 7.0, 0.0, 0.0, 0.0, 0.0, 11.0, 13.0, 12.0, 1.0, 0.0]}]}
	inputData := []float32{0.0, 0.0, 1.0, 11.0, 14.0, 15.0, 3.0, 0.0, 0.0, 1.0, 13.0, 16.0, 12.0, 16.0, 8.0, 0.0, 0.0, 8.0, 16.0, 4.0, 6.0, 16.0, 5.0, 0.0, 0.0, 5.0, 15.0, 11.0, 13.0, 14.0, 0.0, 0.0, 0.0, 0.0, 2.0, 12.0, 16.0, 13.0, 0.0, 0.0, 0.0, 0.0, 0.0, 13.0, 16.0, 16.0, 6.0, 0.0, 0.0, 0.0, 0.0, 16.0, 16.0, 16.0, 7.0, 0.0, 0.0, 0.0, 0.0, 11.0, 13.0, 12.0, 1.0, 0.0}
	inferInputs := []*pb.ModelInferRequest_InferInputTensor{
		{
			Name:     "predict",
			Datatype: "FP32",
			Shape:    []int64{1, 64},
			Contents: &pb.InferTensorContents{
				Fp32Contents: inputData,
			},
		},
	}

	inferOutputs := []*pb.ModelInferRequest_InferRequestedOutputTensor{}

	inferReq := pb.ModelInferRequest{
		ModelName: modelName,
		Inputs:    inferInputs,
		Outputs:   inferOutputs,
	}
	return &inferReq

}

func (m *MushroomXgboot) GetInferRequest(modelName string) *pb.ModelInferRequest {
	var xgBoostInputData []float32 = []float32{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0}
	inferInputs := []*pb.ModelInferRequest_InferInputTensor{
		{
			Name:     "predict",
			Shape:    []int64{1, 126},
			Datatype: "FP32",
			Contents: &pb.InferTensorContents{Fp32Contents: xgBoostInputData},
		},
	}
	inferOutputs := []*pb.ModelInferRequest_InferRequestedOutputTensor{}
	inferReq := pb.ModelInferRequest{
		ModelName: modelName,
		Inputs:    inferInputs,
		Outputs:   inferOutputs,
	}

	return &inferReq
}

func (m *CifarPytorch) GetInferRequest(modelName string) *pb.ModelInferRequest {
	inferInputs := []*pb.ModelInferRequest_InferInputTensor{
		{
			Name:     "INPUT__0",
			Shape:    []int64{1, 3, 32, 32},
			Datatype: "FP32",
			Contents: &pb.InferTensorContents{Fp32Contents: m.Image},
		},
	}

	inferReq := pb.ModelInferRequest{
		ModelName: modelName,
		Inputs:    inferInputs,
	}
	return &inferReq
}

func (m *MushroomLightgbm) GetInferRequest(modelName string) *pb.ModelInferRequest {
	var lightGBMInputData []float32 = []float32{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0}

	inferInputs := []*pb.ModelInferRequest_InferInputTensor{
		{
			Name:     "predict",
			Shape:    []int64{1, 126},
			Datatype: "FP32",
			Contents: &inference.InferTensorContents{Fp32Contents: lightGBMInputData},
		},
	}

	inferReq := pb.ModelInferRequest{
		ModelName: modelName,
		Inputs:    inferInputs,
	}
	return &inferReq
}

func (m *MnistOnnx) GetInferRequest(modelName string) *pb.ModelInferRequest {
	inferInputs := []*pb.ModelInferRequest_InferInputTensor{
		{
			Name:     "Input3",
			Shape:    []int64{1, 1, 28, 28},
			Datatype: "FP32",
			Contents: &inference.InferTensorContents{Fp32Contents: m.Image},
		},
	}

	inferOutputs := []*pb.ModelInferRequest_InferRequestedOutputTensor{}
	inferReq := pb.ModelInferRequest{
		ModelName: modelName,
		Inputs:    inferInputs,
		Outputs:   inferOutputs,
	}
	return &inferReq
}
