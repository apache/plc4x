package read

import (
	"encoding/json"
	"fmt"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go"
)

func main() int {
	// Get a connection to a remote PLC
	crc := plc4go.NewPlcDriverManager().GetConnectedConnection("s7://192.168.23.30")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		_ = fmt.Errorf("error connecting to PLC: %s", connectionResult.Err.Error())
		return 1
	}
	connection := connectionResult.Connection

	// Make sure the connection is closed at the end
	defer connection.Close()

	// Prepare a write-request
	rrb := connection.WriteRequestBuilder()
	rrb.AddField("output-field", "%Q0.0:BOOL", true)
	rrb.AddField("input-field", "%I0.0:USINT", 42)
	readRequest, err := rrb.Build()
	if err != nil {
		_ = fmt.Errorf("error preparing read-request: %s", connectionResult.Err.Error())
		return 2
	}

	// Execute a read-request
	rrc := readRequest.Execute()

	// Wait for the response to finish
	rrr := <-rrc
	if rrr.Err != nil {
		_ = fmt.Errorf("error executing read-request: %s", rrr.Err.Error())
		return 3
	}

	// Do something with the response
	readResponseJson, err := json.Marshal(rrr.Response)
	if err != nil {
		_ = fmt.Errorf("error serializing read-response: %s", err.Error())
		return 4
	}
	fmt.Printf("Result: %s\n", string(readResponseJson))

	return 0
}
