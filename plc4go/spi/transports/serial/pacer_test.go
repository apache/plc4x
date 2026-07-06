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
	assert.Less(t, time.Since(start), 20*time.Millisecond)
}

func TestPacer_EnforcesGapAfterActivity(t *testing.T) {
	p := newPacer(60 * time.Millisecond)
	p.noteActivity()
	start := time.Now()
	require.NoError(t, p.waitTurn(time.Time{}))
	assert.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond, "must wait out the gap")
}

func TestPacer_NoWaitWhenGapAlreadyElapsed(t *testing.T) {
	p := newPacer(30 * time.Millisecond)
	p.noteActivity()
	time.Sleep(40 * time.Millisecond)
	start := time.Now()
	require.NoError(t, p.waitTurn(time.Time{}))
	assert.Less(t, time.Since(start), 20*time.Millisecond)
}

func TestPacer_DeadlineShorterThanGapFailsFast(t *testing.T) {
	p := newPacer(500 * time.Millisecond)
	p.noteActivity()
	start := time.Now()
	err := p.waitTurn(time.Now().Add(20 * time.Millisecond))
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	assert.Less(t, time.Since(start), 200*time.Millisecond, "must not sleep out the full gap")
}

func TestPacer_ActivityDuringWaitExtendsGap(t *testing.T) {
	p := newPacer(80 * time.Millisecond)
	p.noteActivity()
	go func() {
		time.Sleep(40 * time.Millisecond)
		p.noteActivity() // traffic arrives while a writer is waiting its turn
	}()
	start := time.Now()
	require.NoError(t, p.waitTurn(time.Time{}))
	assert.GreaterOrEqual(t, time.Since(start), 110*time.Millisecond,
		"the gap must restart from the mid-wait activity (40ms + 80ms)")
}
