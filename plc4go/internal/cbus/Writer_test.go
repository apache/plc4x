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
	"github.com/stretchr/testify/assert"
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
)

func TestNewWriter(t *testing.T) {
	type args struct {
		tpduGenerator *AlphaGenerator
		messageCodec  *MessageCodec
		tm            spi.RequestTransactionManager
	}
	tests := []struct {
		name string
		args args
		want Writer
	}{
		// TODO: Add test cases.
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
		tm             spi.RequestTransactionManager
	}
	type args struct {
		ctx          context.Context
		writeRequest apiModel.PlcWriteRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan apiModel.PlcWriteRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Writer{
				alphaGenerator: tt.fields.alphaGenerator,
				messageCodec:   tt.fields.messageCodec,
				tm:             tt.fields.tm,
			}
			assert.Equalf(t, tt.want, m.Write(tt.args.ctx, tt.args.writeRequest), "Write(%v, %v)", tt.args.ctx, tt.args.writeRequest)
		})
	}
}
