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

package model

import (
	"context"
	"strconv"
)

func Utf8LengthToPascalLength(_ context.Context, stringValue string) int32 {
	if stringValue == "" {
		return -1
	}
	return int32(len(stringValue))
}

func PascalLengthToUtf8Length(_ context.Context, slength int32) int32 {
	return max(slength, 0)
}

func ExtensionId(_ context.Context, expandedNodeId ExpandedNodeId) int32 {
	nodeId, err := strconv.ParseUint(expandedNodeId.GetNodeId().GetIdentifier(), 10, 16)
	if err != nil {
		return -1
	}
	return int32(nodeId)
}
