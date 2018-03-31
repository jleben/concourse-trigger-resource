package main

import (
    "encoding/json"
    "io"
    "os"
    "os/exec"
//    "path/filepath"
    "fmt"
    "github.com/jleben/trigger-resource/protocol"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

    destination := os.Args[1]

    var request protocol.InRequest

    var err error

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("Parsing request.", err)
	}

	println("Request parsed.")

	fmt.Printf("channel: %v\n", request.Source.Channel)
	fmt.Printf("token: %v\n", request.Source.Token)
    fmt.Printf("version: %v\n", request.Version.Request)

    target_request := protocol.TargetInRequest {
        request.Source.Target,
        request.Version.Target,
    }

    target_response := input_target(target_request, destination)

    var response protocol.InResponse

    response.Version = protocol.Version{
        request.Version.Request,
        target_response.Version,
    }

    err = json.NewEncoder(os.Stdout).Encode(&response)
    if err != nil {
        fatal("encoding response", err)
    }
}

func input_target(request protocol.TargetInRequest, destination string) (protocol.TargetInResponse) {

    var err error

    var target_request_bytes []byte

    target_request_bytes, err = json.Marshal(request)
    if err != nil {
        fatal("Serializing target JSON input.", err)
    }

    fmt.Printf("target request: %v\n", string(target_request_bytes))

    prefix := os.Getenv("PROXY_RESOURCE_PREFIX")
    if len(prefix) == 0 {
        prefix = "."
    }

    cmd := exec.Command(prefix + "/in", destination)

    var in_pipe io.WriteCloser
    var out_pipe io.ReadCloser

    in_pipe, err = cmd.StdinPipe()
    if err != nil {
        fatal("getting stdin pipe", err)
    }

    out_pipe, err = cmd.StdoutPipe()
    if err != nil {
        fatal("getting stdout pipe", err)
    }

    err = cmd.Start()
    if err != nil {
        fatal("starting target", err)
    }

    in_pipe.Write(target_request_bytes)
    in_pipe.Close()

    var response protocol.TargetInResponse
    json.NewDecoder(out_pipe).Decode(&response)

    err = cmd.Wait()
    if err != nil {
        fatal("waiting for target", err)
    }

    {
        response_bytes, err := json.Marshal(response)
        if err != nil {
            fatal("marshaling response", err)
        }
        fmt.Printf("target response: %v\n", string(response_bytes))
    }

    return response
}


func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
