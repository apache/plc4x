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

import "context"

// Count computes the 2 bit size field of an SDO initiate upload response or
// download request. Ported from plc4j StaticHelper.count: for an expedited,
// size-indicated transfer carrying its data inline (at most 4 bytes), size is
// the number of unused bytes, 4 - len(data); for every other payload it is 0.
func Count(_ context.Context, expedited bool, indicated bool, payload SDOInitiateUploadResponsePayload) uint8 {
	if !expedited || !indicated {
		return 0
	}
	if inline, ok := payload.(SDOInitiateExpeditedUploadResponse); ok {
		return uint8(4 - len(inline.GetData()))
	}
	return 0
}
