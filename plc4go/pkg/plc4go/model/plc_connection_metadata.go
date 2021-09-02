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

package model

// PlcConnectionMetadata Information about connection capabilities.
// This includes connection and driver specific metadata.
type PlcConnectionMetadata interface {

	// GetConnectionAttributes Gives access to a map of additional information the driver might be able to provide.
	GetConnectionAttributes() map[string]string

	// CanRead Indicates that the connection supports reading.
	CanRead() bool
	// CanWrite Indicates that the connection supports writing.
	CanWrite() bool
	// CanSubscribe Indicates that the connection supports subscription.
	CanSubscribe() bool
	// CanBrowse Indicates that the connection supports browsing.
	CanBrowse() bool
}
