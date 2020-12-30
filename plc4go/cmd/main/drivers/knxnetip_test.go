//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package drivers

import (
    "encoding/hex"
    "fmt"
    "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip"
    "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports/udp"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "github.com/apache/plc4x/plc4go/pkg/plc4go"
    apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "testing"
    "time"
)

func KnxNetIp(t *testing.T) {
    t.Skip()
    request, err := hex.DecodeString("000a00000006010300000004")
    if err != nil {
        t.Errorf("Error decoding test input")
    }
    rb := utils.NewReadBuffer(request)
    adu, err := model.KnxNetIpMessageParse(rb)
    if err != nil {
        t.Errorf("Error parsing: %s", err)
    }
    if adu != nil {
        // Output success ...
    }

}

func TestKnxNetIpPlc4goDiscovery(t *testing.T) {
    driverManager := plc4go.NewPlcDriverManager()
    driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
    driverManager.RegisterTransport(udp.NewUdpTransport())

    found := make(chan bool)
    err := driverManager.Discover(func(event apiModel.PlcDiscoveryEvent) {
        fmt.Printf("Found device: %s:%s://%s\n - Name: %s\n", event.ProtocolCode, event.TransportCode,
            event.TransportUrl.Host, event.Name)
        found <- true
    })
    if err != nil {
        fmt.Printf("got error %s", err.Error())
        return
    }
    for {
        select {
        case _ = <-found:
            time.Sleep(time.Second * 2)
            fmt.Print("Found devices")
            return
        case <-time.After(time.Second * 10):
            t.Error("Couldn't find device")
            t.Fail()
            return
        }
    }
}

func KnxNetIpPlc4goDriver(t *testing.T) {
    driverManager := plc4go.NewPlcDriverManager()
    driverManager.RegisterDriver(knxnetip.NewKnxNetIpDriver())
    driverManager.RegisterTransport(udp.NewUdpTransport())

    // Get a connection to a remote PLC
    crc := driverManager.GetConnection("knxnet-ip://192.168.42.11")
    //crc := driverManager.GetConnection("knxnet-ip://-discover-")

    // Wait for the driver to connect (or not)
    connectionResult := <-crc
    if connectionResult.Err != nil {
        t.Errorf("error connecting to PLC: %s", connectionResult.Err.Error())
        t.Fail()
        return
    }
    connection := connectionResult.Connection

    attributes := connection.GetMetadata().GetConnectionAttributes()
    fmt.Printf("Successfully connected to KNXnet/IP Gateway '%s' with KNX address '%s' got assigned client KNX address '%s'\n",
        attributes["GatewayName"],
        attributes["GatewayKnxAddress"],
        attributes["ClientKnxAddress"])

    // Try to ping the remote device
    pingResultChannel := connection.Ping()
    pingResult := <-pingResultChannel
    if pingResult.Err != nil {
        t.Errorf("couldn't ping device: %s", pingResult.Err.Error())
        t.Fail()
        return
    }

    // Make sure the connection is closed at the end
    defer connection.Close()

    // Prepare a read-request
    /*pollingInterval, err := time.ParseDuration("5s")
      if err != nil {
          t.Errorf("invalid format")
          t.Fail()
          return
      }*/
    srb := connection.SubscriptionRequestBuilder()
    srb.AddChangeOfStateItem("heating-actual-temperature", "*/*/10:DPT_Value_Temp")
    srb.AddChangeOfStateItem("heating-target-temperature", "*/*/11:DPT_Value_Temp")
    srb.AddChangeOfStateItem("heating-valve-open", "*/*/12:DPT_OpenClose")
    srb.AddItemHandler(knxEventHandler)
    subscriptionRequest, err := srb.Build()
    if err != nil {
        t.Errorf("error preparing subscription-request: %s", connectionResult.Err.Error())
        t.Fail()
        return
    }

    // Execute a subscription-request
    rrc := subscriptionRequest.Execute()

    // Wait for the response to finish
    rrr := <-rrc
    if rrr.Err != nil {
        t.Errorf("error executing read-request: %s", rrr.Err.Error())
        t.Fail()
        return
    }

    // Wait 2 minutes
    time.Sleep(120 * time.Second)

    // Execute a read request
    rrb := connection.ReadRequestBuilder()
    rrb.AddItem("energy-consumption", "1/1/211:DPT_Value_Power")
    rrb.AddItem("actual-temperatures", "*/*/10:DPT_Value_Temp")
    rrb.AddItem("set-temperatures", "*/*/11:DPT_Value_Temp")
    rrb.AddItem("window-status", "*/*/[60,64]:DPT_Value_Temp")
    rrb.AddItem("power-consumption", "*/*/[111,121,131,141]:DPT_Value_Temp")
    readRequest, err := rrb.Build()
    if err == nil {
        rrr := readRequest.Execute()
        readRequestResult := <-rrr
        if readRequestResult.Err == nil {
            for _, fieldName := range readRequestResult.Response.GetFieldNames() {
                if readRequestResult.Response.GetResponseCode(fieldName) == apiModel.PlcResponseCode_OK {
                    fmt.Printf(" - Field %s Value %s\n", fieldName, readRequestResult.Response.GetValue(fieldName).GetString())
                }
            }
        }
    }
    /*value1 := rrr.Response.GetValue("field1")
      value2 := rrr.Response.GetValue("field2")
      fmt.Printf("\n\nResult field1: %f\n", value1.GetFloat32())
      fmt.Printf("\n\nResult field1: %f\n", value2.GetFloat32())

      // Prepare a write-request
      wrb := connection.WriteRequestBuilder()
      wrb.AddItem("field1", "holding-register:1:REAL", 1.2345)
      wrb.AddItem("field2", "holding-register:3:REAL", 2.3456)
      writeRequest, err := rrb.Build()
      if err != nil {
          t.Errorf("error preparing read-request: %s", connectionResult.Err.Error())
          t.Fail()
          return
      }

      // Execute a write-request
      wrc := writeRequest.Execute()

      // Wait for the response to finish
      wrr := <-wrc
      if wrr.Err != nil {
          t.Errorf("error executing read-request: %s", rrr.Err.Error())
          t.Fail()
          return
      }

      fmt.Printf("\n\nResult field1: %d\n", wrr.Response.GetResponseCode("field1"))
      fmt.Printf("\n\nResult field2: %d\n", wrr.Response.GetResponseCode("field2"))*/
}

func knxEventHandler(event apiModel.PlcSubscriptionEvent) {
    for _, fieldName := range event.GetFieldNames() {
        if event.GetResponseCode(fieldName) == apiModel.PlcResponseCode_OK {
            groupAddress := event.GetAddress(fieldName)
            fmt.Printf("Got update for field %s with address %s. Value changed to: %s\n",
                fieldName, groupAddress, event.GetValue(fieldName).GetString())
        }
    }
}
