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

package utils

import (
	"context"
)

type key int

const keyArrayInfo key = iota

type arrayInfo struct {
	numItems int
	curItem  int
}

// CreateArrayContext creates a new context, which contains information on the size and the current position in the array.
func CreateArrayContext(parent context.Context, numItems int, curItem int) context.Context {
	return context.WithValue(parent, keyArrayInfo, arrayInfo{
		numItems: numItems,
		curItem:  curItem,
	})
}

// GetNumItemsFromContext helper to get access to the numItems value stored in the context.
func GetNumItemsFromContext(ctx context.Context) int {
	return ctx.Value(keyArrayInfo).(arrayInfo).numItems
}

// GetCurItemFromContext helper to get access to the curItem value stored in the context.
func GetCurItemFromContext(ctx context.Context) int {
	return ctx.Value(keyArrayInfo).(arrayInfo).curItem
}

// GetLastItemFromContext helper to get access to the lastItem value which is calculated from data stored in the context.
func GetLastItemFromContext(ctx context.Context) bool {
	info := ctx.Value(keyArrayInfo).(arrayInfo)
	return info.curItem == (info.numItems - 1)
}
