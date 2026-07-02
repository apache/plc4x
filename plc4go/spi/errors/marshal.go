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

import "fmt"

// Constants compatible with github.com/rs/zerolog/pkgerrors so callers can do
// entry[errors.StackSourceFileName] and keep working when swapped over.
const (
	StackSourceFileName     = "source"
	StackSourceLineName     = "line"
	StackSourceFunctionName = "func"
)

// StackMarshalPanicKey is the key under which MarshalStack reports a panic it
// recovered while walking the error chain (instead of stack frames).
const StackMarshalPanicKey = "error"

// MarshalStack walks err's chain looking for a value that implements
//
//	interface{ StackTrace() StackTrace }
//
// and returns its frames as []map[string]string keyed by StackSource*Name.
// Returns nil if no stack trace is found in the chain. Suitable for use as
// zerolog.ErrorStackMarshaler.
//
// The walk is panic-guarded: error chains can contain improperly constructed
// values (e.g. a typed-nil *net.OpError wrapped via %w) whose Unwrap method
// dereferences a nil receiver. This marshaler runs inside zerolog on every
// logged error, so a panic here would take down whatever worker is logging.
// Instead of swallowing the anomaly, the recovered panic is reported as the
// marshaled payload - it lands under the "stack" key of the very log event
// that was logging the error, so no logger is needed at this level.
func MarshalStack(err error) (result any) {
	defer func() {
		if r := recover(); r != nil {
			result = []map[string]string{{
				StackMarshalPanicKey: fmt.Sprintf("panic while extracting stack from error chain: %v", r),
			}}
		}
	}()
	type stackTracer interface {
		StackTrace() StackTrace
	}
	var st stackTracer
	for err != nil {
		if t, ok := err.(stackTracer); ok {
			st = t
			break
		}
		err = Unwrap(err)
	}
	if st == nil {
		return nil
	}
	frames := st.StackTrace()
	out := make([]map[string]string, 0, len(frames))
	s := &fmtState{}
	for _, f := range frames {
		out = append(out, map[string]string{
			StackSourceFileName:     frameField(f, s, 's'),
			StackSourceLineName:     frameField(f, s, 'd'),
			StackSourceFunctionName: frameField(f, s, 'n'),
		})
	}
	return out
}

// fmtState is a minimal fmt.State that captures the formatted bytes for a single Frame.Format call.
type fmtState struct {
	b []byte
}

func (s *fmtState) Write(b []byte) (int, error) { s.b = b; return len(b), nil }
func (s *fmtState) Width() (int, bool)          { return 0, false }
func (s *fmtState) Precision() (int, bool)      { return 0, false }
func (s *fmtState) Flag(int) bool               { return false }

func frameField(f Frame, s *fmtState, verb rune) string {
	f.Format(s, verb)
	return string(s.b)
}
