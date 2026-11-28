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

package socketcan

import (
	"net/url"
	"testing"
)

func TestErrorMaskOption(t *testing.T) {
	transport := NewTransport()
	for _, test := range []struct {
		values []string
		want   uint32
		fails  bool
	}{
		{values: nil, want: 0},
		{values: []string{"0x1FFFFFFF"}, want: 0x1FFFFFFF},
		{values: []string{"4"}, want: 4},
		{values: []string{"0x20000000"}, fails: true},
		{values: []string{"bus-off"}, fails: true},
		{values: []string{"1", "2"}, fails: true},
	} {
		options := map[string][]string{}
		if test.values != nil {
			options[ErrorMaskOption] = test.values
		}
		instance, err := transport.CreateTransportInstance(url.URL{Scheme: TransportCode, Host: "can0"}, options)
		if test.fails {
			if err == nil {
				t.Errorf("%v: accepted, want an error", test.values)
			}
			continue
		}
		if err != nil {
			t.Errorf("%v: %v", test.values, err)
			continue
		}
		if got := instance.(*TransportInstance).ErrorMask; got != test.want {
			t.Errorf("%v: mask 0x%X, want 0x%X", test.values, got, test.want)
		}
	}
}
