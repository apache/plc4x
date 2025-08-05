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

package codegen

import "context"

type key int

var lastItem key

var curItem key

// NewContextLastItem returns a new Context that carries value u.
func NewContextLastItem(ctx context.Context, u bool) context.Context {
	return context.WithValue(ctx, lastItem, u)
}

// FromContextLastItem returns the last item flag value stored in ctx, if any.
func FromContextLastItem(ctx context.Context) (bool, bool) {
	u, ok := ctx.Value(lastItem).(bool)
	return u, ok
}

// NewContextCurItem returns a new Context that carries value u.
func NewContextCurItem(ctx context.Context, u int) context.Context {
	return context.WithValue(ctx, curItem, u)
}

// FromContextCurItem returns the last item flag value stored in ctx, if any.
func FromContextCurItem(ctx context.Context) (int, bool) {
	u, ok := ctx.Value(ctx).(int)
	return u, ok
}
