package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//MLService is an interface for an implementation of a machine learning algorithm
type MLService interface {
	Fit(input []interface{}, labels []interface{}) ([]byte, error)
	Predict(input []interface{}, trainedModel []byte) ([]byte, error)
}

type pythonSample struct {
	Client *rpc.Client
}

type FitRequest struct {
	Input  []interface{}
	Labels []interface{}
}

type PredictRequest struct {
	Input        interface{}
	TrainedModel string
}

type FitResponse struct {
	TrainedModel string
}

type PredictResponse struct {
	PredictedLabel interface{}
}

//Start server and wait
func startServer() (int, *os.Process, error) {
	port := 8000 + rand.Intn(1000)
	cmd := exec.Command("python", "model.py", strconv.Itoa(port))
	err := cmd.Start()
	go func() {
		err := cmd.Wait()
		fmt.Println(err)
	}()

	return port, cmd.Process, err
}

func killServer(proc *os.Process) error {
	return proc.Kill()
}

func (p *pythonSample) Fit(input []interface{}, labels []interface{}) ([]byte, error) {

	args := FitRequest{
		Input:  input,
		Labels: labels,
	}

	var reply FitResponse

	err := p.Client.Call("fit", args, &reply)

	return []byte(reply.TrainedModel), err

}

func (p *pythonSample) Predict(input interface{}, trainedModel []byte) ([]byte, error) {

	args := PredictRequest{
		Input:        input,
		TrainedModel: string(trainedModel),
	}

	var reply FitResponse

	err := p.Client.Call("predict", args, &reply)

	return []byte(reply.TrainedModel), err
}

func main() {

	port, proc, err := startServer()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting for server to start...")

	var client *rpc.Client

	for {
		client, err = jsonrpc.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))

		if err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	fmt.Println("Server started")

	p := pythonSample{Client: client}

	fittedModel, err := p.Fit([]interface{}{1, 2, 3}, []interface{}{"a", "b", "c"})
	if err != nil {
		fmt.Println("Error doing RPC")
		log.Fatal(err)
	}

	//some code to persist goes here

	prediction, err := p.Predict(1, fittedModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Prediction: %s", prediction)

	killServer(proc)
}
