/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
)

func main() {
	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterAdsDriver(driverManager)
	connectionChan := driverManager.GetConnection("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851")
	connection := <-connectionChan
	readRequest, err := connection.GetConnection().ReadRequestBuilder().
		AddTagAddress("value-bool", "MAIN.hurz_BOOL").
		AddTagAddress("value-byte", "MAIN.hurz_BYTE").
		AddTagAddress("value-word", "MAIN.hurz_WORD").
		AddTagAddress("value-dword", "MAIN.hurz_DWORD").
		AddTagAddress("value-lword", "MAIN.hurz_LWORD").
		AddTagAddress("value-sint", "MAIN.hurz_SINT").
		AddTagAddress("value-usint", "MAIN.hurz_USINT").
		AddTagAddress("value-int", "MAIN.hurz_INT").
		AddTagAddress("value-uint", "MAIN.hurz_UINT").
		AddTagAddress("value-dint", "MAIN.hurz_DINT").
		AddTagAddress("value-udint", "MAIN.hurz_UDINT").
		AddTagAddress("value-lint", "MAIN.hurz_LINT").
		AddTagAddress("value-ulint", "MAIN.hurz_ULINT").
		AddTagAddress("value-real", "MAIN.hurz_REAL").
		AddTagAddress("value-lreal", "MAIN.hurz_LREAL").
		AddTagAddress("value-string", "MAIN.hurz_STRING").
		AddTagAddress("value-wstring", "MAIN.hurz_WSTRING").
		AddTagAddress("value-time", "MAIN.hurz_TIME").
		AddTagAddress("value-ltime", "MAIN.hurz_LTIME").
		AddTagAddress("value.date", "MAIN.hurz_DATE").
		AddTagAddress("value-time-of-day", "MAIN.hurz_TIME_OF_DAY").
		AddTagAddress("value-date-and-time", "MAIN.hurz_DATE_AND_TIME").
		AddTagAddress("value-struct", "MAIN.hurz_Struct").
		Build()
	if err != nil {
		panic(err)
	}
	readResponseChannel := readRequest.Execute()
	readResult := <-readResponseChannel
	readResponse := readResult.GetResponse()
	for _, tagName := range readResponse.GetTagNames() {
		plcValue := readResponse.GetValue(tagName)
		print(plcValue)
	}
}
