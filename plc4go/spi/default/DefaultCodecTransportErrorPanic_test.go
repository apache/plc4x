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

package _default

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// A transport error chain can contain improperly constructed values - most
// notably a typed-nil *net.OpError wrapped into the chain. stdErrors.Is
// dereferences such values while unwrapping ((*net.OpError).Unwrap reads
// e.Err on a nil receiver), which killed the receive worker with a
// nil-pointer panic in the field. handleTransportError must survive such
// chains and classify them as fatal instead of panicking.
func Test_defaultCodec_handleTransportError_typedNilOpErrorInChain(t *testing.T) {
	var nilOpErr *net.OpError
	poisoned := fmt.Errorf("read failed: %w", error(nilOpErr))

	codec := &defaultCodec{}
	assert.NotPanics(t, func() {
		keepRunning := codec.handleTransportError(testutils.ProduceTestingLogger(t), poisoned)
		assert.False(t, keepRunning, "poisoned chain must be treated as fatal and stop the worker")
	})
}
