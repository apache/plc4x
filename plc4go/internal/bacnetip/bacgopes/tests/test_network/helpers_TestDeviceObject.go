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

package test_network

import (
	"fmt"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
)

type TestDeviceObject struct {
	device.LocalDeviceObject
	*DefaultRFormatter
}

func NewTestDeviceObject(args comp.Args, kwArgs comp.KWArgs) *TestDeviceObject {
	return &TestDeviceObject{
		LocalDeviceObject: device.NewLocalDeviceObject(args, kwArgs),
		DefaultRFormatter: NewDefaultRFormatter(),
	}
}

func (t *TestDeviceObject) Format(s fmt.State, v rune) {
	t.DefaultRFormatter.Format(s, v)
}

func (t *TestDeviceObject) String() string {
	return t.DefaultRFormatter.String()
}
