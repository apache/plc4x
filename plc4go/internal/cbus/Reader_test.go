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

package cbus

import (
	"context"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	type args struct {
		tpduGenerator *AlphaGenerator
		messageCodec  *MessageCodec
		tm            spi.RequestTransactionManager
	}
	tests := []struct {
		name string
		args args
		want *Reader
	}{
		{
			name: "create a new one",
			want: &Reader{
				alphaGenerator: nil,
				messageCodec:   nil,
				tm:             nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewReader(tt.args.tpduGenerator, tt.args.messageCodec, tt.args.tm), "NewReader(%v, %v, %v)", tt.args.tpduGenerator, tt.args.messageCodec, tt.args.tm)
		})
	}
}

func TestReader_Read(t *testing.T) {
	type fields struct {
		alphaGenerator *AlphaGenerator
		messageCodec   *MessageCodec
		tm             spi.RequestTransactionManager
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool
	}{
		{
			name: "read and bail",
			args: args{
				ctx: context.Background(),
				readRequest: spiModel.NewDefaultPlcReadRequest(nil, func() []string {
					return strings.Split(strings.Repeat("asd,", 40), ",")
				}(), nil, nil),
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool {
				timer := time.NewTimer(2 * time.Second)
				defer timer.Stop()
				select {
				case <-timer.C:
					t.Fail()
				case result := <-results:
					assert.NotNil(t, result.GetErr())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Reader{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Read(tt.args.ctx, tt.args.readRequest)), "Read(%v, %v)", tt.args.ctx, tt.args.readRequest)
		})
	}
}

func TestReader_readSync(t *testing.T) {
	type fields struct {
		alphaGenerator *AlphaGenerator
		messageCodec   *MessageCodec
		tm             spi.RequestTransactionManager
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
		result      chan apiModel.PlcReadRequestResult
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		resultEvaluator func(t *testing.T, results chan apiModel.PlcReadRequestResult) bool
	}{
		{
			name: "too many tags",
			args: args{
				ctx: context.Background(),
				readRequest: spiModel.NewDefaultPlcReadRequest(nil, func() []string {
					return strings.Split(strings.Repeat("asd,", 40), ",")
				}(), nil, nil),
				result: make(chan apiModel.PlcReadRequestResult, 1),
			},
			resultEvaluator: func(t *testing.T, results chan apiModel.PlcReadRequestResult) bool {
				timer := time.NewTimer(2 * time.Second)
				defer timer.Stop()
				select {
				case <-timer.C:
					t.Fail()
				case result := <-results:
					assert.NotNil(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "read something without any tag",
			args: args{
				ctx: context.Background(),
				readRequest: spiModel.NewDefaultPlcReadRequest(
					map[string]apiModel.PlcTag{},
					[]string{},
					nil,
					nil,
				),
				result: make(chan apiModel.PlcReadRequestResult, 1),
			},
			resultEvaluator: func(t *testing.T, results chan apiModel.PlcReadRequestResult) bool {
				timer := time.NewTimer(2 * time.Second)
				defer timer.Stop()
				select {
				case <-timer.C:
					t.Fail()
				case result := <-results:
					assert.Nil(t, result.GetErr())
					assert.NotNil(t, result.GetResponse())
				}
				return true
			},
		},
		{
			name: "read something",
			fields: fields{
				alphaGenerator: &AlphaGenerator{currentAlpha: 'g'},
				messageCodec: &MessageCodec{
					requestContext: readWriteModel.NewRequestContext(false),
					cbusOptions:    readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false),
				},
				tm: spi.NewRequestTransactionManager(10),
			},
			args: args{
				ctx: func() context.Context {
					timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					t.Cleanup(cancel)
					return timeout
				}(),
				readRequest: spiModel.NewDefaultPlcReadRequest(
					map[string]apiModel.PlcTag{
						"blub": NewCALIdentifyTag(nil, nil, readWriteModel.Attribute_Manufacturer, 1),
					},
					[]string{
						"blub",
					},
					nil,
					nil,
				),
				result: make(chan apiModel.PlcReadRequestResult, 1),
			},
			resultEvaluator: func(t *testing.T, results chan apiModel.PlcReadRequestResult) bool {
				timer := time.NewTimer(2 * time.Second)
				defer timer.Stop()
				select {
				case <-timer.C:
					t.Fail()
				case result := <-results:
					assert.Nil(t, result.GetErr())
					assert.NotNil(t, result.GetResponse())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Reader{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			m.readSync(tt.args.ctx, tt.args.readRequest, tt.args.result)
			assert.True(t, tt.resultEvaluator(t, tt.args.result))
		})
	}
}
