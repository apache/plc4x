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

package s7

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"reflect"
	"strings"
	"testing"
)

func TestS7MessageBytes(t *testing.T) {
	type debuggable interface {
		utils.Serializable
		fmt.Stringer
		utils.AsciiBoxer
	}
	type args struct {
		debuggable debuggable
	}
	tests := []struct {
		name       string
		args       args
		wantString string
		wantDump   string
	}{
		{
			name: "TPKT Packet with Read var response data",
			args: args{
				debuggable: model.NewTPKTPacket(
					model.NewCOTPPacketData(
						false,
						13,
						[]*model.COTPParameter{model.NewCOTPParameterTpduSize(model.COTPTpduSize_SIZE_4096)},
						model.NewS7MessageResponseData(
							0,
							0,
							11,
							model.NewS7ParameterReadVarResponse(4),
							model.NewS7PayloadReadVarResponse(
								[]*model.S7VarPayloadDataItem{
									model.NewS7VarPayloadDataItem(
										model.DataTransportErrorCode_OK,
										model.DataTransportSize_BIT,
										[]int8{1},
									),
								},
							),
						),
					),
				),
			},
			wantString: `
╔═TPKTPacket════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═Payload═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗║
║║╔═Eot═╗╔═TpduRef╗╔═Parameters═════════════╗╔═Payload════════════════════════════════════════════════════════════════════════════════════════╗║║
║║║false║║   13   ║║╔═COTPParameterTpduSize╗║║╔═TpduReference╗╔═ErrorClass╗╔═ErrorCode╗╔═Parameter═╗╔═Payload════════════════════════════════╗║║║
║║╚═════╝╚════════╝║║     ╔═TpduSize╗      ║║║║      11      ║║     0     ║║    0     ║║╔═NumItems╗║║╔═Items════════════════════════════════╗║║║║
║║                 ║║     ║SIZE_4096║      ║║║╚══════════════╝╚═══════════╝╚══════════╝║║    4    ║║║║╔═S7VarPayloadDataItem═══════════════╗║║║║║
║║                 ║║     ╚═════════╝      ║║║                                         ║╚═════════╝║║║║╔═ReturnCode╗╔═TransportSize╗╔═Data╗║║║║║║
║║                 ║╚══════════════════════╝║║                                         ╚═══════════╝║║║║    OK     ║║     BIT      ║║ ╔═╗ ║║║║║║║
║║                 ╚════════════════════════╝║                                                      ║║║╚═══════════╝╚══════════════╝║ ║1║ ║║║║║║║
║║                                           ║                                                      ║║║                             ║ ╚═╝ ║║║║║║║
║║                                           ║                                                      ║║║                             ╚═════╝║║║║║║
║║                                           ║                                                      ║║╚════════════════════════════════════╝║║║║║
║║                                           ║                                                      ║╚══════════════════════════════════════╝║║║║
║║                                           ║                                                      ╚════════════════════════════════════════╝║║║
║║                                           ╚════════════════════════════════════════════════════════════════════════════════════════════════╝║║
║╚═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║
╚═══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
			wantDump: `
00|03 00 00 1d 05 f0 0d c0 01 0c '..........'
10|32 03 00 00 00 0b 00 02 00 05 '2.........'
20|00 00 04 04 ff 03 00 01 01    '......... '
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantString = strings.Trim(tt.wantString, "\n")
			if got := tt.args.debuggable.String(); got != tt.wantString {
				t.Errorf("String() = '\n%v\n', want '\n%v\n'", got, tt.wantString)
			}
			buffer := utils.NewWriteBuffer()
			err := tt.args.debuggable.Serialize(*buffer)
			if err != nil {
				t.Error(err)
			}
			tt.wantDump = strings.Trim(tt.wantDump, "\n")
			if got := utils.Dump(buffer.GetBytes()); !reflect.DeepEqual(got, tt.wantDump) {
				t.Errorf("Serialize() = '\n%v\n', want '\n%v\n'", got, tt.wantDump)
			}
		})
	}
}
