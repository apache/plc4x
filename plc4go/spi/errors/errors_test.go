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

package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// ── Basic constructors ─────────────────────────────────────────────────────

func TestNew(t *testing.T) {
	err := errors.New("boom")
	assert.Equal(t, "boom", err.Error())
}

func TestErrorf(t *testing.T) {
	err := errors.Errorf("count=%d name=%s", 7, "x")
	assert.Equal(t, "count=7 name=x", err.Error())
}

// ── nil propagation ────────────────────────────────────────────────────────

func TestNilPropagation(t *testing.T) {
	assert.Nil(t, errors.Wrap(nil, "msg"))
	assert.Nil(t, errors.Wrapf(nil, "msg %d", 1))
	assert.Nil(t, errors.WithStack(nil))
	assert.Nil(t, errors.WithMessage(nil, "msg"))
	assert.Nil(t, errors.WithMessagef(nil, "msg %d", 1))
}

// ── Wrap / Wrapf / WithMessage / WithMessagef messages ─────────────────────

func TestWrap_Message(t *testing.T) {
	root := errors.New("root")
	wrapped := errors.Wrap(root, "outer")
	assert.Equal(t, "outer: root", wrapped.Error())
}

func TestWrapf_Message(t *testing.T) {
	root := errors.New("root")
	wrapped := errors.Wrapf(root, "outer-%d", 42)
	assert.Equal(t, "outer-42: root", wrapped.Error())
}

func TestWithMessage_NoStack(t *testing.T) {
	root := errors.New("root")
	wrapped := errors.WithMessage(root, "outer")
	assert.Equal(t, "outer: root", wrapped.Error())
}

func TestWithMessagef(t *testing.T) {
	root := errors.New("root")
	wrapped := errors.WithMessagef(root, "outer-%d", 5)
	assert.Equal(t, "outer-5: root", wrapped.Error())
}

// ── Cause / Unwrap chains ──────────────────────────────────────────────────

func TestCause_StopsAtFundamental(t *testing.T) {
	root := errors.New("root")
	err := errors.Wrap(errors.Wrap(root, "mid"), "outer")
	assert.Equal(t, "root", errors.Cause(err).Error())
}

func TestCause_StopsAtNonCauser(t *testing.T) {
	std := errors.New("std-root")
	err := errors.Wrap(std, "outer")
	// Cause walks through withStack → withMessage → std (which has no Cause).
	assert.Same(t, std, errors.Cause(err))
}

func TestCause_NilReturnsNil(t *testing.T) {
	assert.Nil(t, errors.Cause(nil))
}

func TestUnwrap(t *testing.T) {
	root := errors.New("root")
	wrapped := errors.Wrap(root, "outer")
	// withStack.Unwrap → withMessage; withMessage.Unwrap → root.
	mid := errors.Unwrap(wrapped)
	require.NotNil(t, mid)
	assert.Equal(t, root, errors.Unwrap(mid))
}

// ── stdlib errors.Is / errors.As interop ───────────────────────────────────

var sentinel = errors.New("sentinel")

func TestIs_WrapChainMatchesSentinel(t *testing.T) {
	err := errors.Wrap(sentinel, "outer")
	assert.True(t, errors.Is(err, sentinel))
}

type customErr struct{ code int }

func (c *customErr) Error() string { return fmt.Sprintf("custom-%d", c.code) }

func TestAs_FindsConcreteType(t *testing.T) {
	root := &customErr{code: 17}
	err := errors.Wrap(root, "outer")
	var got *customErr
	require.True(t, errors.As(err, &got))
	assert.Equal(t, 17, got.code)
}

// ── Stack traces ───────────────────────────────────────────────────────────

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func TestNew_HasStackTrace(t *testing.T) {
	err := errors.New("boom")
	st, ok := err.(stackTracer)
	require.True(t, ok, "expected fundamental to implement stackTracer")
	assert.NotEmpty(t, st.StackTrace())
}

func TestErrorf_HasStackTrace(t *testing.T) {
	err := errors.Errorf("boom-%d", 1)
	st, ok := err.(stackTracer)
	require.True(t, ok)
	assert.NotEmpty(t, st.StackTrace())
}

func TestWrap_HasStackTrace(t *testing.T) {
	err := errors.Wrap(errors.New("std"), "outer")
	st, ok := err.(stackTracer)
	require.True(t, ok)
	assert.NotEmpty(t, st.StackTrace())
}

func TestWithStack_HasStackTrace(t *testing.T) {
	err := errors.WithStack(errors.New("std"))
	st, ok := err.(stackTracer)
	require.True(t, ok)
	assert.NotEmpty(t, st.StackTrace())
}

func TestWithMessage_DoesNotAddStack(t *testing.T) {
	std := errors.New("std")
	err := errors.WithMessage(std, "outer")
	_, ok := err.(stackTracer)
	assert.False(t, ok, "WithMessage must not attach a stack trace")
}

// ── Format verbs ──────────────────────────────────────────────────────────

func TestFormat_New_S(t *testing.T) {
	err := errors.New("boom")
	assert.Equal(t, "boom", fmt.Sprintf("%s", err))
	assert.Equal(t, "boom", fmt.Sprintf("%v", err))
	assert.Equal(t, `"boom"`, fmt.Sprintf("%q", err))
}

func TestFormat_Wrap_S(t *testing.T) {
	err := errors.Wrap(errors.New("root"), "outer")
	assert.Equal(t, "outer: root", fmt.Sprintf("%s", err))
}

func TestFormat_PlusV_IncludesStack(t *testing.T) {
	err := errors.New("boom")
	rendered := fmt.Sprintf("%+v", err)
	assert.Contains(t, rendered, "boom")
	// %+v should dump frames — at minimum, this test function's name appears.
	assert.Contains(t, rendered, "TestFormat_PlusV_IncludesStack")
}

// ── Frame helpers ─────────────────────────────────────────────────────────

func TestStackTrace_NonEmptyTopFrame(t *testing.T) {
	err := errors.New("boom")
	st := err.(stackTracer).StackTrace()
	require.NotEmpty(t, st)

	// Top frame should point at this test function.
	rendered := fmt.Sprintf("%+v", st[0])
	assert.Contains(t, rendered, "TestStackTrace_NonEmptyTopFrame")
}

func TestFrame_MarshalText(t *testing.T) {
	err := errors.New("boom")
	st := err.(stackTracer).StackTrace()
	require.NotEmpty(t, st)
	text, mErr := st[0].MarshalText()
	require.NoError(t, mErr)
	// Contains funcname + file path + line number.
	assert.Contains(t, string(text), "TestFrame_MarshalText")
}

// ── Format of nil-wrap callers ────────────────────────────────────────────

func TestFormat_NilSafe(t *testing.T) {
	// %+v on a nil-from-Wrap shouldn't panic.
	err := errors.Wrap(nil, "outer")
	assert.Nil(t, err)
	_, _ = fmt.Fprintf(io.Discard, "%+v", err) // smoke: nil-safe
}

// ── stdlib re-exports (so callers only need this one import) ──────────────

func TestJoin_NilEntriesIgnored(t *testing.T) {
	assert.Nil(t, errors.Join(nil, nil))
}

func TestJoin_CombinesErrors(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")
	joined := errors.Join(a, b)
	require.NotNil(t, joined)
	assert.True(t, errors.Is(joined, a))
	assert.True(t, errors.Is(joined, b))
}

func TestErrUnsupported_EqualsStdlib(t *testing.T) {
	// Same sentinel — identity-compare must hold so existing errors.Is checks
	// against stdlib's symbol keep matching.
	assert.Same(t, stderrors.ErrUnsupported, errors.ErrUnsupported)
}
