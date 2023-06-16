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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	driverManager := plc4go.NewPlcDriverManager()
	defer func() {
		if err := driverManager.Close(); err != nil {
			panic(err)
		}
	}()
	drivers.RegisterAdsDriver(driverManager)
	connectionChan := driverManager.GetConnection("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851")
	connection := <-connectionChan

	subscriptionRequest, err := connection.GetConnection().SubscriptionRequestBuilder().
		AddChangeOfStateTagAddress("value-int", "MAIN.rivianTest01.HorizontalPosition").
		AddPreRegisteredConsumer("value-int", func(event apiModel.PlcSubscriptionEvent) {
			value := event.GetValue("value-int")
			log.Info().Msgf("Got value: %d", value.GetUint16())
		}).
		Build()
	if err != nil {
		panic(err)
	}
	subscriptionResponseChannel := subscriptionRequest.Execute()
	subscriptionResult := <-subscriptionResponseChannel
	if subscriptionResult.GetErr() != nil {
		log.Error().Err(subscriptionResult.GetErr()).Msg("error in response")
		return
	}
	subscriptionResponse := subscriptionResult.GetResponse()
	for _, tagName := range subscriptionResponse.GetTagNames() {
		responseCode := subscriptionResponse.GetResponseCode(tagName)
		print(responseCode)
	}

	time.Sleep(time.Second * 200)

}
