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

package model

import (
	"testing"
)

func TestDF1UtilsCrcCheck(t *testing.T) {
	type args struct {
		destinationAddress uint8
		sourceAddress      uint8
		command            *DF1Command
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name: "Test example crc ",
			args: args{
				0x05,
				0x07,
				func() *DF1Command {
					response := &DF1UnprotectedReadResponse{
						Data: []uint8{
							0x00, 0x00, 0x00, 0x00,
							0x00, 0x00, 0x00, 0x00,
							0x00, 0x00, 0x00, 0x00,
						},
					}
					command := &DF1Command{
						0x00,
						0xAFFE,
						response,
					}
					response.Parent = command
					return command
				}(),
			},
			want:    0xBE4D,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DF1UtilsCrcCheck(tt.args.destinationAddress, tt.args.sourceAddress, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("DF1UtilsCrcCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DF1UtilsCrcCheck() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
