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
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

func TestNewWriter(t *testing.T) {
	type args struct {
		tpduGenerator *AlphaGenerator
		messageCodec  *MessageCodec
		tm            transactions.RequestTransactionManager
	}
	tests := []struct {
		name string
		args args
		want *Writer
	}{
		{
			name: "create a new one",
			want: func() *Writer {
				var writer Writer
				writer.log = log.Logger
				return &writer
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewWriter(tt.args.tpduGenerator, tt.args.messageCodec, tt.args.tm), "NewWriter(%v, %v, %v)", tt.args.tpduGenerator, tt.args.messageCodec, tt.args.tm)
		})
	}
}

func TestWriter_Write(t *testing.T) {
	type fields struct {
		alphaGenerator *AlphaGenerator
		messageCodec   *MessageCodec
		tm             transactions.RequestTransactionManager
	}
	type args struct {
		ctx          context.Context
		writeRequest apiModel.PlcWriteRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		setup        func(t *testing.T, fields *fields, args *args)
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool
	}{
		{
			name: "write something",
			args: args{
				writeRequest: spiModel.NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				args.ctx = testutils.TestContext(t)
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
					t.FailNow()
				case result := <-results:
					assert.NotNil(t, result)
					assert.Nil(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "too many tags",
			args: args{
				writeRequest: spiModel.NewDefaultPlcWriteRequest(nil, func() []string {
					return strings.Split(strings.Repeat("asd,", 30), ",")
				}(), nil, nil, nil),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				args.ctx = testutils.TestContext(t)
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Fatal("timeout")
				case result := <-results:
					assert.NotNil(t, result)
					assert.NotNil(t, result.GetErr())
					assert.Equal(t, "Only 20 tags can be handled at once", result.GetErr().Error())
				}
				return true
			},
		},
		/*
			TODO: implement once we have a writable tag
			{
				name: "one tag",
				fields: fields{
					alphaGenerator: &AlphaGenerator{
						currentAlpha: 'g',
					},
					messageCodec: NewMessageCodec(func() transports.TransportInstance {
						transport := test.NewTransport()
						instance, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
						if err != nil {
							t.Fatal(err)
						}
						return instance
					}()),
					tm: spi.NewRequestTransactionManager(10),
				},
				args: args{
					writeRequest: spiModel.NewDefaultPlcWriteRequest(
						map[string]apiModel.PlcTag{
							"asd": &statusTag{},
						},
						[]string{
							"asd",
						},
						nil,
						nil,
						nil,
					),
				},
				setup: func(t *testing.T, fields *fields, args *args){
					args.ctx = testutils.TestContext(t)
				},
				wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
					timeout := time.NewTimer(2 * time.Second)
					utils.CleanupTimer(timeout)
					select {
					case <-timeout.C:
						t.Error("timeout")
						t.FailNow()
					case result := <-results:
						assert.NotNil(t, result)
						assert.NotNil(t, result.GetErr())
						assert.Equal(t, "Only 20 tags can be handled at once", result.GetErr().Error())
					}
					return true
				},
			},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := Writer{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Write(tt.args.ctx, tt.args.writeRequest)), "Write(%v, %v)", tt.args.ctx, tt.args.writeRequest)
			if m.messageCodec != nil {
				t.Log(m.messageCodec.Disconnect())
			}
			if m.tm != nil {
				t.Log(m.tm.Close())
			}
		})
	}
}
