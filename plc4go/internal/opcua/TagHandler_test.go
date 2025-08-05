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

package opcua

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

func TestNewTagHandler(t *testing.T) {
	tests := []struct {
		name string
		want TagHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTagHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTagHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOPCUAAddressPatterns(t *testing.T) {
	addressPattern := NewTagHandler().tagAddress
	t.Run("Address", func(t *testing.T) {
		addresses := []string{
			//standard integer based param
			"ns=2;i=10846",
			//string based address values
			"ns=2;s=test.variable.name.inspect",
			"ns=2;s=::AsGlobalPV:ProductionOrder",
			"ns=2;s=::AsGlobalPV:ProductionOrder;BOOL",
			"ns=2;s=key param with some spaces",
			"ns=2;s=\"aweired\".\"siemens\".\"param\".\"submodule\".\"param",
			"ns=2;s=Weee314Waannaaa\\somenice=ext=a234a*#+1455!§$%&/()tttraaaaSymbols-:.,,",
			// GUID address tests
			"ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a",
			// binary encoded addresses
			"ns=2;b=asvaewavarahreb==",
		}
		for _, address := range addresses {
			t.Run(address, func(t *testing.T) {
				assert.True(t, addressPattern.MatchString(address))
			})
		}
	})
	t.Run("AddressDataType", func(t *testing.T) {
		addresses := []string{
			//standard integer based param
			"ns=2;i=10846;BOOL",
			//string based address values
			"ns=2;s=test.variable.name.inspect;DINT",
			"ns=2;s=key param with some spaces;ULINT",
			"ns=2;s=\"aweired\".\"siemens\".\"param\".\"submodule\".\"param;LREAL",
			//REGEX Valid, additional checks need to be done later
			"ns=2;s=Weee314Waannaaa\\somenice=ext=a234a*#+1455!§$%&/()tttraaaaSymbols-.,,;JIBBERISH",
		}
		for _, address := range addresses {
			t.Run(address, func(t *testing.T) {
				assert.True(t, addressPattern.MatchString(address))
			})
		}
	})
	t.Run("AddressDataType that don't match", func(t *testing.T) {
		addresses := []string{
			// GUID address tests
			"ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a;*&#%^*$(*)",
			// binary encoded addresses
			"ns=2;b=asvae;wavarahreb==",
		}
		for _, address := range addresses {
			t.Run(address, func(t *testing.T) {
				assert.False(t, addressPattern.MatchString(address))
			})
		}
	})
}

func TestTagHandler_ParseQuery(t *testing.T) {
	type fields struct {
		tagAddress *regexp.Regexp
	}
	type args struct {
		in0 string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.PlcQuery
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				tagAddress: tt.fields.tagAddress,
			}
			got, err := m.ParseQuery(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_ParseTag(t *testing.T) {
	type fields struct {
		tagAddress *regexp.Regexp
	}
	type args struct {
		tagAddress string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				tagAddress: tt.fields.tagAddress,
			}
			got, err := m.ParseTag(tt.args.tagAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleTagAddress(t *testing.T) {
	type fields struct {
		tagAddress *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				tagAddress: tt.fields.tagAddress,
			}
			got, err := m.handleTagAddress(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleTagAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleTagAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
