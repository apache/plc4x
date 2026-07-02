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
)

// ErrCorruptErrorChain marks an error whose original chain could not be walked
// safely. Detect it with errors.Is(err, ErrCorruptErrorChain) after the error
// passed through SanitizeError.
var ErrCorruptErrorChain = stderrors.New("corrupt error chain (panic during unwrap)")

// sanitizeMaxDepth bounds the chain walk so a (broken) cyclic chain classifies
// as corrupt instead of hanging the caller.
const sanitizeMaxDepth = 256

// SanitizeError returns err unchanged when its chain can be walked safely.
//
// If walking the chain panics - which happens when improperly constructed
// values (e.g. a typed-nil *net.OpError wrapped via %w) are part of the chain,
// as their Unwrap/Error methods dereference a nil receiver - it returns a
// flattened copy of the error message joined with ErrCorruptErrorChain. The
// result is safe for errors.Is/As/Unwrap, keeps the original message in the
// error string, and carries the anomaly as an Is-detectable sentinel, so it
// stays visible wherever the caller eventually logs or propagates the error.
// No logging happens at this level.
//
// Note: joining the ORIGINAL error with the sentinel would not defuse
// anything - Join's Unwrap() []error hands the poisoned chain right back to
// every future walk. That is why the message is flattened instead.
func SanitizeError(err error) error {
	if err == nil {
		return nil
	}
	if chainWalksSafely(err) {
		return err
	}
	return stderrors.Join(stderrors.New(safeErrorString(err)), ErrCorruptErrorChain)
}

// chainWalksSafely performs the same traversal errors.Is would (including
// multi-child Join nodes) and reports whether it completes without panicking.
func chainWalksSafely(err error) (safe bool) {
	defer func() {
		if r := recover(); r != nil {
			safe = false
		}
	}()
	return walkChain(err, 0)
}

func walkChain(err error, depth int) bool {
	for err != nil {
		if depth > sanitizeMaxDepth {
			return false
		}
		depth++
		_ = err.Error() // Error() itself may dereference a nil receiver
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Unwrap() []error }:
			for _, child := range x.Unwrap() {
				if !walkChain(child, depth) {
					return false
				}
			}
			return true
		default:
			return true
		}
	}
	return true
}

// safeErrorString renders err.Error() with a panic guard, falling back to a
// placeholder when even the message cannot be produced.
func safeErrorString(err error) (s string) {
	defer func() {
		if r := recover(); r != nil {
			s = fmt.Sprintf("(error message unavailable: panic in Error(): %v)", r)
		}
	}()
	return err.Error()
}
