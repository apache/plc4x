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

package serial

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPacer_NilAndZeroDelayAreNoops(t *testing.T) {
	var nilPacer *pacer
	nilPacer.noteActivity()
	require.NoError(t, nilPacer.waitTurn(time.Time{}))

	zero := newPacer(0)
	zero.noteActivity()
	start := time.Now()
	require.NoError(t, zero.waitTurn(time.Time{}))
	// Generously bounded: the claim is "did not wait", and no configured gap exists to wait out,
	// so the bound only has to be tighter than a wait - not tighter than a scheduling stall.
	assert.Less(t, time.Since(start), 500*time.Millisecond)
}

func TestPacer_EnforcesGapAfterActivity(t *testing.T) {
	p := newPacer(60 * time.Millisecond)
	// Sampled before noteActivity, so it is never later than the pacer's own stamp: the gap it
	// then has to wait out is at least this long, whatever the machine does in between. Measuring
	// from after the call would let a stall eat into the margin.
	activity := time.Now()
	p.noteActivity()
	require.NoError(t, p.waitTurn(time.Time{}))
	assert.GreaterOrEqual(t, time.Since(activity), 60*time.Millisecond, "must wait out the gap")
}

func TestPacer_NoWaitWhenGapAlreadyElapsed(t *testing.T) {
	p := newPacer(200 * time.Millisecond)
	p.noteActivity()
	time.Sleep(250 * time.Millisecond)
	start := time.Now()
	require.NoError(t, p.waitTurn(time.Time{}))
	// The gap is 200ms, so returning within 100ms proves it was not waited out. The wide margin
	// between the two is what keeps a scheduling stall from failing the test.
	assert.Less(t, time.Since(start), 100*time.Millisecond, "gap already elapsed")
}

func TestPacer_DeadlineShorterThanGapFailsFast(t *testing.T) {
	p := newPacer(500 * time.Millisecond)
	p.noteActivity()
	start := time.Now()
	err := p.waitTurn(time.Now().Add(20 * time.Millisecond))
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	// Failing fast means not sleeping out the 500ms gap, so any bound below that proves the
	// point - this one leaves room for a scheduling stall on a loaded machine.
	assert.Less(t, time.Since(start), 400*time.Millisecond, "must not sleep out the full gap")
}

func TestPacer_ActivityDuringWaitExtendsGap(t *testing.T) {
	const gap = 200 * time.Millisecond
	const midWait = 50 * time.Millisecond
	// The premise is that the traffic goroutine gets scheduled while the writer is still waiting
	// its turn. A stalled machine can deliver it only after the original gap has already expired -
	// such a run proves nothing about the pacer, so it is retried rather than failed.
	for attempt := 1; ; attempt++ {
		p := newPacer(gap)
		// Sampled before noteActivity, so it is never later than the pacer's own stamp.
		firstActivity := time.Now()
		p.noteActivity()
		activityCh := make(chan time.Time, 1)
		go func() {
			time.Sleep(midWait)
			// Stamped before the call, so it is never later than the pacer's own stamp.
			activityCh <- time.Now()
			p.noteActivity() // traffic arrives while a writer is waiting its turn
		}()
		require.NoError(t, p.waitTurn(time.Time{}))
		returned := time.Now()
		activity := <-activityCh

		if activity.Sub(firstActivity) >= gap {
			if attempt >= 3 {
				t.Skip("machine too stalled to deliver the mid-wait activity inside the gap")
			}
			continue
		}
		// The contract: the gap restarts from that activity. Both stamps come from the same clock
		// and bracket the pacer's own, so this holds no matter how the goroutines were scheduled.
		assert.GreaterOrEqual(t, returned.Sub(activity), gap,
			"the gap must restart from the mid-wait activity")
		return
	}
}
