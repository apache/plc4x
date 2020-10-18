package write

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
	wrb := connection.WriteRequestBuilder()
	wrb.AddField("output-field", "%Q0.0:BOOL", true)
	wrb.AddField("input-field", "%I0.0:USINT", 42)
	writeRequest, err := wrb.Build()
	if err != nil {
		_ = fmt.Errorf("error preparing read-request: %s", connectionResult.Err.Error())
		return 2
	}

	// Execute a read-request
	wrc := writeRequest.Execute()

	// Wait for the response to finish
	wrr := <-wrc
	if wrr.Err != nil {
		_ = fmt.Errorf("error executing write-request: %s", wrr.Err.Error())
		return 3
	}

	// Do something with the response
	writeResponseJson, err := json.Marshal(wrr.Response)
	if err != nil {
		_ = fmt.Errorf("error serializing write-response: %s", err.Error())
		return 4
	}
	fmt.Printf("Result: %s\n", string(writeResponseJson))

	return 0
}
