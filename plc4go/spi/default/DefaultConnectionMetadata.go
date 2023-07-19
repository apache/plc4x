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

package _default

// DefaultConnectionMetadata implements the model.PlcConnectionMetadata interface
type DefaultConnectionMetadata struct {
	ConnectionAttributes map[string]string
	ProvidesReading      bool
	ProvidesWriting      bool
	ProvidesSubscribing  bool
	ProvidesBrowsing     bool
}

func (m DefaultConnectionMetadata) GetConnectionAttributes() map[string]string {
	return m.ConnectionAttributes
}

func (m DefaultConnectionMetadata) CanRead() bool {
	return m.ProvidesReading
}

func (m DefaultConnectionMetadata) CanWrite() bool {
	return m.ProvidesWriting
}

func (m DefaultConnectionMetadata) CanSubscribe() bool {
	return m.ProvidesSubscribing
}

func (m DefaultConnectionMetadata) CanBrowse() bool {
	return m.ProvidesBrowsing
}
