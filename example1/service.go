package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

//MLService is an interface for an implementation of a machine learning algorithm
type MLService interface {
	Fit(data []interface{}, labels []interface{}) ([]byte, error)
	Predict(data []interface{}, trainedModel []byte) ([]byte, error)
}

type pythonSample struct{}

//Run, Pipe, Wait
func runPipeWait(args ...string) ([]byte, error) {
	var (
		output    []byte
		errOutput []byte
		err       error
	)
	cmd := exec.Command("python", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if output, err = ioutil.ReadAll(stdout); err != nil {
		return nil, err
	}

	if errOutput, err = ioutil.ReadAll(stderr); err != nil || len(errOutput) > 0 {
		return nil, fmt.Errorf("Error running model: %s", string(errOutput))
	}

	return output, nil
}

func (p *pythonSample) Fit(data []interface{}, labels []interface{}) ([]byte, error) {
	var (
		paramData   []byte
		paramLabels []byte
		err         error
	)

	if paramData, err = json.Marshal(data); err != nil {
		return nil, err
	}
	if paramLabels, err = json.Marshal(labels); err != nil {
		return nil, err
	}

	return runPipeWait("model.py", "fit", string(paramData), string(paramLabels))
}

func (p *pythonSample) Predict(data interface{}, trainedModel []byte) ([]byte, error) {
	var (
		paramData []byte
		err       error
	)

	if paramData, err = json.Marshal(data); err != nil {
		return nil, err
	}

	return runPipeWait("model.py", "predict", string(paramData), string(trainedModel))
}

func main() {
	p := new(pythonSample)
	fittedModel, err := p.Fit([]interface{}{1, 2, 3}, []interface{}{"a", "b", "c"})
	if err != nil {
		log.Fatal(err)
	}
	prediction, err := p.Predict(1, fittedModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Prediction: %s", prediction)
}
