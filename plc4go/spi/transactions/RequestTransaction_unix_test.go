//go:build unix

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

package transactions

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/options"
)

// processCPUTime is the CPU time, user and system, the test process has used so far.
func processCPUTime(t *testing.T) time.Duration {
	t.Helper()
	var usage syscall.Rusage
	require.NoError(t, syscall.Getrusage(syscall.RUSAGE_SELF, &usage))
	return time.Duration(usage.Utime.Nano() + usage.Stime.Nano())
}

// A driver's operation typically returns once its request is on the wire, and the transaction ends only when the
// reply (or its timeout) comes in. AwaitCompletion waits for that end; it must not burn a CPU core while doing so,
// as a driver with a transaction manager per connection has one such wait per connection with a request out.
func Test_requestTransaction_AwaitCompletion_waitsForTheEndWithoutSpinning(t *testing.T) {
	rtm := NewRequestTransactionManager(1, options.WithCustomLogger(produceTestingLogger(t)))
	t.Cleanup(func() {
		assert.NoError(t, rtm.Close())
	})

	const replyAfter = 500 * time.Millisecond
	transaction := rtm.StartTransaction("request on the wire")
	transaction.Submit("request on the wire", func(_ context.Context, transaction RequestTransaction) {
		// Sent: the reply ends the transaction later.
		time.AfterFunc(replyAfter, func() {
			assert.NoError(t, transaction.EndRequest())
		})
	})

	cpuBefore := processCPUTime(t)
	started := time.Now()
	require.NoError(t, transaction.AwaitCompletion(t.Context()))
	waited := time.Since(started)
	cpu := processCPUTime(t) - cpuBefore

	assert.GreaterOrEqual(t, waited, replyAfter-50*time.Millisecond, "AwaitCompletion returned before the transaction ended")
	assert.Less(t, cpu, waited/2, "waiting %s for the transaction to end took %s of CPU", waited, cpu)
}
