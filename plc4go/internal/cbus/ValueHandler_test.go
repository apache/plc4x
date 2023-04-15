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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/stretchr/testify/assert"
	"testing"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func TestNewValueHandler(t *testing.T) {
	tests := []struct {
		name string
		want ValueHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewValueHandler(), "NewValueHandler()")
		})
	}
}

func TestValueHandler_NewPlcValue(t *testing.T) {
	type fields struct {
		DefaultValueHandler spiValues.DefaultValueHandler
	}
	type args struct {
		tag   model.PlcTag
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiValues.PlcValue
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ValueHandler{
				DefaultValueHandler: tt.fields.DefaultValueHandler,
			}
			got, err := m.NewPlcValue(tt.args.tag, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("NewPlcValue(%v, %v)", tt.args.tag, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewPlcValue(%v, %v)", tt.args.tag, tt.args.value)
		})
	}
}
