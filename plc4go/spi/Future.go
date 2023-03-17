/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package spi

import "time"

type CompletableFuture[T any] interface {
	// Complete If not already completed, sets the value returned by get()
	// and related methods to the given value.
	Complete(value T)
	// CompleteWithError If not already completed, causes invocations of get()
	// and related methods to return the given error.
	CompleteWithError(err error)
	// Cancel If not already completed, complete this CompletableFuture with
	// a CancellationError
	Cancel()

	// IsDone Returns true if completed in any fashion: normally, exceptionally,
	// or via cancellation.
	IsDone() bool
	// IsCancelled Returns true if this CompletableFuture was cancelled before
	// it completed normally.
	IsCancelled() bool
	// IsCompletedWithError Returns true if this CompletableFuture completed
	// with an error, in any way (Cancelled or completed with error).
	IsCompletedWithError() bool

	// Get Waits if necessary for this future to complete, and then returns its
	// result.
	Get() T
	// GetWithTimeout Waits if necessary for at most the given time for this
	// future to complete, and then returns its result, if available.
	GetWithTimeout(timeout time.Duration) T
	// GetNow Returns the result value (or throws any encountered exception)
	// if completed, else returns the given valueIfAbsent.
	GetNow(valueIfAbsent T) T
	// WhenComplete When this CompletableFuture completes either normally or
	// exceptionally, it passes this CompletableFutures result and error as
	// arguments to the supplied function.
	WhenComplete(func(value T, err error))
	// HandleAsync When this CompletableFuture completes either normally or
	// exceptionally, it passes this CompletableFutures result and error as
	// arguments to the supplied function.
	HandleAsync(func(value T, err error))
	// ThenApply When this CompletableFuture completes normally, it passes its
	// result as the argument to the supplied function.
	ThenApply(func(value T))
}
