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
	"github.com/apache/plc4x/plc4go/spi/values"
	"github.com/rs/zerolog/log"
)

func main() {
	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterAdsDriver(driverManager)
	connectionChan := driverManager.GetConnection("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851")
	connection := <-connectionChan
	writeRequest, err := connection.GetConnection().WriteRequestBuilder().
		AddTagAddress("value-bool", "MAIN.hurz_BOOL", values.NewPlcBOOL(true)).                   // 1
		AddTagAddress("value-byte", "MAIN.hurz_BYTE", values.NewPlcBYTE(42)).                     // 1
		AddTagAddress("value-word", "MAIN.hurz_WORD", values.NewPlcWORD(42424)).                  // 2
		AddTagAddress("value-dword", "MAIN.hurz_DWORD", values.NewPlcDWORD(4242442424)).          // 4
		AddTagAddress("value-lword", "MAIN.hurz_LWORD", values.NewPlcLWORD(4242442424242424242)). // 8
		AddTagAddress("value-sint", "MAIN.hurz_SINT", values.NewPlcSINT(-42)).                    // 1
		AddTagAddress("value-usint", "MAIN.hurz_USINT", values.NewPlcUSINT(42)).                  // 1
		AddTagAddress("value-int", "MAIN.hurz_INT", values.NewPlcINT(-2424)).                     // 2
		AddTagAddress("value-uint", "MAIN.hurz_UINT", values.NewPlcUINT(42424)).                  // 2
		AddTagAddress("value-dint", "MAIN.hurz_DINT", values.NewPlcDINT(-242442424)).             // 4
		AddTagAddress("value-udint", "MAIN.hurz_UDINT", values.NewPlcUDINT(4242442424)).          // 4
		AddTagAddress("value-lint", "MAIN.hurz_LINT", values.NewPlcLINT(-4242442424242424242)).   // 8
		AddTagAddress("value-ulint", "MAIN.hurz_ULINT", values.NewPlcULINT(4242442424242424242)). // 8
		AddTagAddress("value-real", "MAIN.hurz_REAL", values.NewPlcREAL(3.14159265359)).          // 4
		AddTagAddress("value-lreal", "MAIN.hurz_LREAL", values.NewPlcLREAL(2.71828182846)).       // 8
		AddTagAddress("value-string", "MAIN.hurz_STRING", values.NewPlcSTRING("hurz")).           // 4
		AddTagAddress("value-wstring", "MAIN.hurz_WSTRING", values.NewPlcWSTRING("wolf")).        // 8
		/*AddTagAddress("value-time", "MAIN.hurz_TIME")).
		AddTagAddress("value-ltime", "MAIN.hurz_LTIME")).
		AddTagAddress("value.date", "MAIN.hurz_DATE")).
		AddTagAddress("value-time-of-day", "MAIN.hurz_TIME_OF_DAY")).
		AddTagAddress("value-date-and-time", "MAIN.hurz_DATE_AND_TIME")).
		AddTagAddress("value-struct", "MAIN.hurz_Struct")).*/
		Build()
	if err != nil {
		panic(err)
	}
	writeResponseChannel := writeRequest.Execute()
	writeResult := <-writeResponseChannel
	if writeResult.GetErr() != nil {
		log.Error().Err(writeResult.GetErr()).Msg("error in response")
		return
	}
	writeResponse := writeResult.GetResponse()
	for _, tagName := range writeResponse.GetTagNames() {
		responseCode := writeResponse.GetResponseCode(tagName)
		print(responseCode)
	}
}
