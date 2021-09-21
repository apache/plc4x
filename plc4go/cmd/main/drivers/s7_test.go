/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package drivers

import (
	"encoding/hex"
	_ "github.com/apache/plc4x/plc4go/cmd/main/initializetest"
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"testing"
)

func TestS7(t *testing.T) {
	t.Skip()
	request, err := hex.DecodeString("000a00000006010300000004")
	if err != nil {
		// Output an error ...
	}
	rb := utils.NewReadBufferByteBased(request)
	adu, err := model.TPKTPacketParse(rb)
	if err != nil {
		t.Errorf("Error parsing: %s", err)
	}
	if adu != nil {
		// Output success ...
	}
}
