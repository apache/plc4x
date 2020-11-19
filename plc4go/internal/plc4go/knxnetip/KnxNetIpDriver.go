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
package knxnetip

import (
    "errors"
    "fmt"
    "net/url"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
    "github.com/apache/plc4x/plc4go/internal/plc4go/transports"
    "github.com/apache/plc4x/plc4go/pkg/plc4go"
)

type KnxNetIpDriver struct {
    fieldHandler spi.PlcFieldHandler
    plc4go.PlcDriver
}

func NewKnxNetIpDriver() *KnxNetIpDriver {
	return &KnxNetIpDriver{
	    fieldHandler: NewFieldHandler(),
    }
}

func (m KnxNetIpDriver) GetProtocolCode() string {
    return "knxnet-ip"
}

func (m KnxNetIpDriver) GetProtocolName() string {
    return "KNXNet/IP"
}

func (m KnxNetIpDriver) GetDefaultTransport() string {
    return "udp"
}

func (m KnxNetIpDriver) CheckQuery(query string) error {
    _, err := m.fieldHandler.ParseQuery(query)
    return err
}

func (m KnxNetIpDriver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
    // If the host is set to "discover", use the KNX discovery mechanism
    if transportUrl.Host == "-discover-" {
        // Multicast address every KNX gateway is required to respond to.
        transportUrl.Host =  "224.0.23.12"
    }

    // Get an the transport specified in the url
    transport, ok := transports[transportUrl.Scheme]
    if !ok {
        ch := make(chan plc4go.PlcConnectionConnectResult)
        ch <- plc4go.NewPlcConnectionConnectResult(nil, errors.New("couldn't find transport for given transport url "+transportUrl.String()))
        return ch
    }
    // Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
    options["defaultUdpPort"] = []string{"3671"}
    // Have the transport create a new transport-instance.
    transportInstance, err := transport.CreateTransportInstance(transportUrl, options)
    if err != nil {
        ch := make(chan plc4go.PlcConnectionConnectResult)
        ch <- plc4go.NewPlcConnectionConnectResult(nil, errors.New("couldn't initialize transport configuration for given transport url "+transportUrl.String()))
        return ch
    }

    // Create a new codec for taking care of encoding/decoding of messages
    defaultChanel := make(chan interface{})
    go func() {
        for {
            _ = <-defaultChanel
            fmt.Printf("Hurz")
/*            adu := model.CastModbusTcpADU(msg)
            serialized, err := json.Marshal(adu)
            if err != nil {
                fmt.Printf("got error serializing adu: %s\n", err.Error())
            } else {
                fmt.Printf("got message in the default handler %s\n", serialized)
            }*/
//            06100201000e0801c0a82a32c70e
//                          xxxx
//            06100201000e08  c0a82a32c95f
        }
    }()
    codec := NewKnxNetIpMessageCodec(transportInstance, nil)

    // Create the new connection
    connection := NewKnxNetIpConnection(codec, options, m.fieldHandler)

    return connection.Connect()
}

