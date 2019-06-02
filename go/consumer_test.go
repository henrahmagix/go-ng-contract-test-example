package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type contract struct {
	raw      string
	Method   string                 `json:"method"`
	Path     string                 `json:"path"`
	Response map[string]interface{} `json:"response"`
}

func loadContract(filename string) (contract, error) {
	var contract contract
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return contract, err
	}
	err = json.Unmarshal(file, &contract)
	contract.raw = string(file)
	return contract, err
}

func TestConsumerGetHello(t *testing.T) {
	contract, err := loadContract("../ng/src/app/contract/get-hello.json")
	require.NoError(t, err)

	method := contract.Method
	req, err := http.NewRequest(method, "http://localhost:8080"+contract.Path, nil)
	require.NoError(t, err)

	sigs := make(chan os.Signal, 1)
	go func() {
		run("localhost:8080", sigs)
	}()

	res, err := doUntilServerReady(req)
	require.NoError(t, err)

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	var body map[string]interface{}
	err = json.Unmarshal(b, &body)
	require.NoError(t, err)

	sigs <- syscall.SIGTERM

	fmt.Printf("contract=%s\n", contract.raw)
	fmt.Printf("response=%s\n", b)
	for key, expectedType := range contract.Response {
		assert.Equal(t, typeof(expectedType), typeof(body[key]), key)
	}
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func doUntilServerReady(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultClient.Do(req)
	if urlErr, ok := err.(*url.Error); ok {
		if opErr, ok := urlErr.Err.(*net.OpError); ok && opErr.Op == "dial" {
			return doUntilServerReady(req)
		}
	}
	return res, err
}
