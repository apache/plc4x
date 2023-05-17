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
	"time"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/rs/zerolog/log"
)

func main() {
	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterAdsDriver(driverManager)
	connectionChan := driverManager.GetConnection("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851")
	connection := <-connectionChan

	duration, _ := time.ParseDuration("1.234S")
	lduration, _ := time.ParseDuration("24015H23M12.034002044S")
	date, _ := time.Parse("2014-11-12", "1978-03-28")
	timeOfDay, _ := time.Parse("11:45:26.371Z", "15:36:30.123")
	dateAndTime, _ := time.Parse("2014-11-12T11:45:26.371Z", "1996-05-06T15:36:30")

	// Prepare the data structure for a custom type.
	children := map[string]apiValues.PlcValue{}
	children["hurz_BOOL"] = spiValues.NewPlcBOOL(true)
	children["hurz_BYTE"] = spiValues.NewPlcBYTE(1)
	children["hurz_WORD"] = spiValues.NewPlcWORD(2)
	children["hurz_DWORD"] = spiValues.NewPlcDWORD(3)
	children["hurz_LWORD"] = spiValues.NewPlcLWORD(4)
	children["hurz_SINT"] = spiValues.NewPlcSINT(5)
	children["hurz_USINT"] = spiValues.NewPlcUSINT(6)
	children["hurz_INT"] = spiValues.NewPlcINT(7)
	children["hurz_UINT"] = spiValues.NewPlcUINT(8)
	children["hurz_DINT"] = spiValues.NewPlcDINT(9)
	children["hurz_UDINT"] = spiValues.NewPlcUDINT(10)
	children["hurz_LINT"] = spiValues.NewPlcLINT(11)
	children["hurz_ULINT"] = spiValues.NewPlcULINT(12)
	children["hurz_REAL"] = spiValues.NewPlcREAL(13.0)
	children["hurz_LREAL"] = spiValues.NewPlcLREAL(14.0)
	children["hurz_STRING"] = spiValues.NewPlcSTRING("hurz")
	children["hurz_WSTRING"] = spiValues.NewPlcWSTRING("wolf")
	children["hurz_TIME"] = spiValues.NewPlcTIME(duration)
	children["hurz_LTIME"] = spiValues.NewPlcLTIME(lduration)
	children["hurz_DATE"] = spiValues.NewPlcDATE(date)
	children["hurz_TIME_OF_DAY"] = spiValues.NewPlcTIME_OF_DAY(timeOfDay)
	children["hurz_DATE_AND_TIME"] = spiValues.NewPlcDATE_AND_TIME(dateAndTime)
	writeRequest, err := connection.GetConnection().WriteRequestBuilder().
		AddTagAddress("value-bool", "MAIN.hurz_BOOL", spiValues.NewPlcBOOL(true)).                   // 1
		AddTagAddress("value-byte", "MAIN.hurz_BYTE", spiValues.NewPlcBYTE(42)).                     // 1
		AddTagAddress("value-word", "MAIN.hurz_WORD", spiValues.NewPlcWORD(42424)).                  // 2
		AddTagAddress("value-dword", "MAIN.hurz_DWORD", spiValues.NewPlcDWORD(4242442424)).          // 4
		AddTagAddress("value-lword", "MAIN.hurz_LWORD", spiValues.NewPlcLWORD(4242442424242424242)). // 8
		AddTagAddress("value-sint", "MAIN.hurz_SINT", spiValues.NewPlcSINT(-42)).                    // 1
		AddTagAddress("value-usint", "MAIN.hurz_USINT", spiValues.NewPlcUSINT(42)).                  // 1
		AddTagAddress("value-int", "MAIN.hurz_INT", spiValues.NewPlcINT(-2424)).                     // 2
		AddTagAddress("value-uint", "MAIN.hurz_UINT", spiValues.NewPlcUINT(42424)).                  // 2
		AddTagAddress("value-dint", "MAIN.hurz_DINT", spiValues.NewPlcDINT(-242442424)).             // 4
		AddTagAddress("value-udint", "MAIN.hurz_UDINT", spiValues.NewPlcUDINT(4242442424)).          // 4
		AddTagAddress("value-lint", "MAIN.hurz_LINT", spiValues.NewPlcLINT(-4242442424242424242)).   // 8
		AddTagAddress("value-ulint", "MAIN.hurz_ULINT", spiValues.NewPlcULINT(4242442424242424242)). // 8
		AddTagAddress("value-real", "MAIN.hurz_REAL", spiValues.NewPlcREAL(3.14159265359)).          // 4
		AddTagAddress("value-lreal", "MAIN.hurz_LREAL", spiValues.NewPlcLREAL(2.71828182846)).       // 8
		AddTagAddress("value-string", "MAIN.hurz_STRING", spiValues.NewPlcSTRING("hurz")).           // 4
		AddTagAddress("value-wstring", "MAIN.hurz_WSTRING", spiValues.NewPlcWSTRING("wolf")).        // 8
		AddTagAddress("value-time", "MAIN.hurz_TIME", spiValues.NewPlcTIME(duration)).
		AddTagAddress("value-ltime", "MAIN.hurz_LTIME", spiValues.NewPlcLTIME(lduration)).
		AddTagAddress("value.date", "MAIN.hurz_DATE", spiValues.NewPlcDATE(date)).
		AddTagAddress("value-time-of-day", "MAIN.hurz_TIME_OF_DAY", spiValues.NewPlcTIME_OF_DAY(timeOfDay)).
		AddTagAddress("value-date-and-time", "MAIN.hurz_DATE_AND_TIME", spiValues.NewPlcDATE_AND_TIME(dateAndTime)).
		AddTagAddress("value-struct", "MAIN.hurz_Struct", spiValues.NewPlcStruct(children)).
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
