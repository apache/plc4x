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

package tests

import (
	"context"
	"testing"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"

	"github.com/stretchr/testify/require"
)

func TestManualBacnet(t *testing.T) {
	testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)
	rawBytes := []byte{
		/*00000000*/ 0x00, 0x22, 0x5f, 0x07, 0xcc, 0x90, 0x00, 0x24, 0xe8, 0xf9, 0x21, 0xa0, 0x08, 0x00, 0x45, 0x00, //|."_....$..!...E.|
		/*00000010*/ 0x00, 0x53, 0x55, 0xb9, 0x00, 0x00, 0x80, 0x11, 0x60, 0xc1, 0xc0, 0xa8, 0x01, 0x68, 0xc0, 0xa8, //|.SU.....`....h..|
		/*00000020*/ 0x01, 0x67, 0xba, 0xc0, 0xba, 0xc0, 0x00, 0x3f, 0xa9, 0x16, 0x81, 0x0a, 0x00, 0x37, 0x01, 0x04, //|.g.....?.....7..|
		/*00000030*/ 0x02, 0x04, 0x47, 0x0f, 0x0c, 0x06, 0xc0, 0x00, 0x01, 0x19, 0x84, 0x3e, 0x0c, 0x00, 0x00, 0x00, //|..G........>....|
		/*00000040*/ 0x01, 0x19, 0x55, 0x3c, 0xff, 0xff, 0xff, 0xff, 0x0c, 0x00, 0x00, 0x00, 0x01, 0x19, 0x6f, 0x3c, //|..U<..........o<|
		/*00000050*/ 0xff, 0xff, 0xff, 0xff, 0x0c, 0x00, 0xc0, 0x00, 0x01, 0x19, 0x55, 0x3c, 0xff, 0xff, 0xff, 0xff, //|..........U<....|
		/*00000060*/ 0x3f, //|?|
	}
	_, err := readWriteModel.BVLCParse(context.TODO(), rawBytes[42:])
	require.NoError(t, err)
}
