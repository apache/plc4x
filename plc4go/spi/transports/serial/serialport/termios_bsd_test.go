//go:build darwin || freebsd

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

package serialport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeTermiosMarkSpaceParityUnsupported(t *testing.T) {
	_, err := makeTermios(Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityMark})
	require.ErrorContains(t, err, "mark/space parity")
	_, err = makeTermios(Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParitySpace})
	require.ErrorContains(t, err, "mark/space parity")
}
