package read

import (
	"encoding/json"
	"fmt"
	"os"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go"
)

func main() {
	// Get a connection to a remote PLC
	connection, err := plc4go.NewPlcDriverManager().GetConnectedConnection("s7://192.168.23.30")
	if err != nil {
		_ = fmt.Errorf("error connecting to PLC: %s", err.Error())
		os.Exit(1)
	}

	// Prepare a read-request
	rrb, err := connection.ReadRequestBuilder()
	if err != nil {
		_ = fmt.Errorf("error getting read-request builder: %s", err.Error())
		os.Exit(2)
	}
	rrb.AddField("output-field", "%Q0.0")
	rrb.AddField("input-field", "I0.0")
	readRequest := rrb.Build()

	// Execute a read-request
	rr, err := readRequest.Execute()
	if err != nil {
		_ = fmt.Errorf("error executing read-request: %s", err.Error())
		os.Exit(2)
	}

	// Wait for the response to finish
	readResponse := <-rr

	// Do something with the response
	readResponseJson, err := json.Marshal(readResponse)
	if err != nil {
		_ = fmt.Errorf("error serializing read-response: %s", err.Error())
		os.Exit(2)
	}
	fmt.Printf("Result: %s\n", string(readResponseJson))
}
