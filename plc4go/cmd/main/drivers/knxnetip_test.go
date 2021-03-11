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
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func TestParser(t *testing.T) {
	request, err := hex.DecodeString("06100202004e080100000000000036010200ff00000000c501010e37e000170c00246d00ec7830302d32342d36442d30302d45432d3738000000000000000000000000000a020201030104010501")
	if err != nil {
		// Output an error ...
	}
	rb := utils.NewReadBuffer(request)
	knxMessage, err := driverModel.KnxNetIpMessageParse(rb)
	if err != nil {
		fmt.Printf("Got error parsing message: %s\n", err.Error())
		// TODO: Possibly clean up ...
		return
	}
	print(knxMessage)
}
