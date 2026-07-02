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

package errors

import (
	stderrors "errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// panickyError panics in Error() - simulating a node whose message cannot even
// be rendered.
type panickyError struct{}

func (p *panickyError) Error() string { panic("boom in Error()") }

// selfWrapping unwraps to itself - a (broken) cyclic chain.
type selfWrapping struct{}

func (s *selfWrapping) Error() string { return "cyclic" }
func (s *selfWrapping) Unwrap() error { return s }

func poisonedChain() error {
	var nilOpErr *net.OpError
	return fmt.Errorf("read failed: %w", error(nilOpErr))
}

func TestSanitizeError_nil(t *testing.T) {
	assert.Nil(t, SanitizeError(nil))
}

func TestSanitizeError_cleanChainReturnedUnchanged(t *testing.T) {
	inner := stderrors.New("inner")
	err := fmt.Errorf("outer: %w", inner)
	sanitized := SanitizeError(err)
	assert.Equal(t, err, sanitized, "clean chains must pass through untouched")
	assert.True(t, stderrors.Is(sanitized, inner), "chain semantics must be preserved")
}

func TestSanitizeError_typedNilInChain(t *testing.T) {
	sanitized := SanitizeError(poisonedChain())
	require.NotNil(t, sanitized)

	// The sanitized error must be safe to walk...
	assert.NotPanics(t, func() {
		_ = stderrors.Is(sanitized, ErrCorruptErrorChain)
		var target *net.OpError
		_ = stderrors.As(sanitized, &target)
	})
	// ...carry the sentinel...
	assert.True(t, stderrors.Is(sanitized, ErrCorruptErrorChain))
	// ...and preserve the original message text.
	assert.Contains(t, sanitized.Error(), "read failed")
}

func TestSanitizeError_typedNilAtTopLevel(t *testing.T) {
	var nilOpErr *net.OpError
	sanitized := SanitizeError(error(nilOpErr))
	require.NotNil(t, sanitized)
	assert.True(t, stderrors.Is(sanitized, ErrCorruptErrorChain))
	assert.NotPanics(t, func() { _ = stderrors.Is(sanitized, stderrors.New("x")) })
}

func TestSanitizeError_panickyErrorMessage(t *testing.T) {
	err := fmt.Errorf("outer: %w", error(&panickyError{}))
	// Chain walk panics via Error() on the child -> corrupt. Note fmt.Errorf
	// pre-renders the message, so the flattened string comes from the wrapper
	// and stays usable; a top-level panicky Error() falls back to a placeholder.
	sanitized := SanitizeError(err)
	assert.True(t, stderrors.Is(sanitized, ErrCorruptErrorChain))

	topLevel := SanitizeError(error(&panickyError{}))
	assert.True(t, stderrors.Is(topLevel, ErrCorruptErrorChain))
	assert.Contains(t, topLevel.Error(), "error message unavailable")
}

func TestSanitizeError_cyclicChain(t *testing.T) {
	sanitized := SanitizeError(error(&selfWrapping{}))
	assert.True(t, stderrors.Is(sanitized, ErrCorruptErrorChain), "cyclic chains must classify as corrupt instead of hanging")
	assert.Contains(t, sanitized.Error(), "cyclic")
}

func TestSanitizeError_joinWithPoisonedChild(t *testing.T) {
	joined := stderrors.Join(stderrors.New("healthy"), poisonedChain())
	sanitized := SanitizeError(joined)
	assert.True(t, stderrors.Is(sanitized, ErrCorruptErrorChain))
	assert.Contains(t, sanitized.Error(), "healthy")
	assert.NotPanics(t, func() { _ = stderrors.Is(sanitized, ErrCorruptErrorChain) })
}

func TestMarshalStack_poisonedChainReportsPanicMarker(t *testing.T) {
	var result any
	assert.NotPanics(t, func() {
		result = MarshalStack(poisonedChain())
	})
	require.NotNil(t, result, "the recovered panic must be reported, not swallowed")
	marker, ok := result.([]map[string]string)
	require.True(t, ok)
	require.Len(t, marker, 1)
	assert.Contains(t, marker[0][StackMarshalPanicKey], "panic while extracting stack")
}
