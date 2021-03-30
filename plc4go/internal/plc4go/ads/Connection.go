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
package ads

import (
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type ConnectionMetadata struct {
}

func (m ConnectionMetadata) GetConnectionAttributes() map[string]string {
	return map[string]string{}
}

func (m ConnectionMetadata) CanRead() bool {
	return true
}

func (m ConnectionMetadata) CanWrite() bool {
	return true
}

func (m ConnectionMetadata) CanSubscribe() bool {
	return true
}

func (m ConnectionMetadata) CanBrowse() bool {
	return false
}

// TODO: maybe we can use a DefaultConnection struct here with delegates
type Connection struct {
	messageCodec       spi.MessageCodec
	options            map[string][]string
	fieldHandler       spi.PlcFieldHandler
	valueHandler       spi.PlcValueHandler
	requestInterceptor internalModel.RequestInterceptor
	// TODO: check if this is the right place here (it is kinda connection bound)
	sourceAmsNetId readWriteModel.AmsNetId
	sourceAmsPort  uint16
	targetAmsNetId readWriteModel.AmsNetId
	targetAmsPort  uint16
}

func NewConnection(messageCodec spi.MessageCodec, options map[string][]string, fieldHandler spi.PlcFieldHandler) (*Connection, error) {
	if err := checkForRequiredParameters(options, []string{"sourceAmsNetId", "sourceAmsPort", "targetAmsNetId", "targetAmsPort"}); err != nil {
		return nil, err
	}
	// TODO: check array
	split := strings.Split(options["sourceAmsNetId"][0], ".")
	octet1, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet2, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet3, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet4, err := strconv.Atoi(split[3])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet5, err := strconv.Atoi(split[4])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	octet6, err := strconv.Atoi(split[5])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsNetId")
	}
	sourceAmsNetId := readWriteModel.AmsNetId{
		Octet1: uint8(octet1),
		Octet2: uint8(octet2),
		Octet3: uint8(octet3),
		Octet4: uint8(octet4),
		Octet5: uint8(octet5),
		Octet6: uint8(octet6),
	}
	// TODO: check array
	sourceAmsPort, err := strconv.Atoi(options["sourceAmsPort"][0])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing sourceAmsPort")
	}
	// TODO: check array
	split = strings.Split(options["targetAmsNetId"][0], ".")
	octet1, err = strconv.Atoi(split[0])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet2, err = strconv.Atoi(split[1])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet3, err = strconv.Atoi(split[2])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet4, err = strconv.Atoi(split[3])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet5, err = strconv.Atoi(split[4])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	octet6, err = strconv.Atoi(split[5])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing targetAmsNetId")
	}
	targetAmsNetId := readWriteModel.AmsNetId{
		Octet1: uint8(octet1),
		Octet2: uint8(octet2),
		Octet3: uint8(octet3),
		Octet4: uint8(octet4),
		Octet5: uint8(octet5),
		Octet6: uint8(octet6),
	}
	// TODO: check array
	targetAmsPort, err := strconv.Atoi(options["targetAmsPort"][0])
	if err != nil {
		return nil, errors.Wrap(err, "error prasing targetAmsPort")
	}
	return &Connection{
		messageCodec:       messageCodec,
		options:            options,
		fieldHandler:       fieldHandler,
		valueHandler:       NewValueHandler(),
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(),
		sourceAmsNetId:     sourceAmsNetId,
		sourceAmsPort:      uint16(sourceAmsPort),
		targetAmsNetId:     targetAmsNetId,
		targetAmsPort:      uint16(targetAmsPort),
	}, nil
}

// TODO: move to a common utils place
func checkForRequiredParameters(options map[string][]string, requiredParameters []string) error {
	for _, parameter := range requiredParameters {
		if options[parameter] == nil {
			return errors.Errorf("required parameter %s missing", parameter)
		}
	}
	return nil
}

func (m Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		err := m.messageCodec.Connect()
		ch <- plc4go.NewPlcConnectionConnectResult(m, err)
	}()
	return ch
}

func (m Connection) BlockingClose() {
	log.Trace().Msg("Closing blocked")
	closeResults := m.Close()
	select {
	case <-closeResults:
		return
	case <-time.After(time.Second * 5):
		return
	}
}

func (m Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	log.Trace().Msg("Close")
	// TODO: Implement ...
	ch := make(chan plc4go.PlcConnectionCloseResult)
	go func() {
		ch <- plc4go.NewPlcConnectionCloseResult(m, nil)
	}()
	return ch
}

func (m Connection) IsConnected() bool {
	panic("implement me")
}

func (m Connection) Ping() <-chan plc4go.PlcConnectionPingResult {
	panic("implement me")
}

func (m Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return ConnectionMetadata{}
}

func (m Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilderWithInterceptor(m.fieldHandler,
		NewReader(m.messageCodec, m.targetAmsNetId, m.targetAmsPort, m.sourceAmsNetId, m.sourceAmsPort), m.requestInterceptor)
}

func (m Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(
		m.fieldHandler, m.valueHandler, NewWriter(m.messageCodec, m.targetAmsNetId, m.targetAmsPort, m.sourceAmsNetId, m.sourceAmsPort))
}

func (m Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	panic("implement me")
}

func (m Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	panic("implement me")
}

func (m Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	panic("implement me")
}

func (m Connection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m Connection) GetPlcFieldHandler() spi.PlcFieldHandler {
	return m.fieldHandler
}

func (m Connection) GetPlcValueHandler() spi.PlcValueHandler {
	return m.valueHandler
}

func (m Connection) String() string {
	return fmt.Sprintf("ads.Connection{}")
}
