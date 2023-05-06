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

package cache

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// cachedPlcConnectionState
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type cachedPlcConnectionState int32

const (
	StateInitialized cachedPlcConnectionState = iota
	StateIdle
	StateInUse
	StateInvalid
)

func (c cachedPlcConnectionState) String() string {
	switch c {
	case StateInitialized:
		return "StateInitialized"
	case StateIdle:
		return "StateIdle"
	case StateInUse:
		return "StateInUse"
	case StateInvalid:
		return "StateInvalid"
	}
	return "Unknown"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Events
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type connectionListener interface {
	// onConnectionEvent: Callback called by the connection container to signal connection events
	// that have an impact on the cache itself (Like connections being permanently closed).
	onConnectionEvent(event connectionEvent)
}

type connectionEvent interface {
	getConnectionContainer() connectionContainer
}

type connectionErrorEvent struct {
	conn connectionContainer
	err  error
}

func (c connectionErrorEvent) getConnectionContainer() connectionContainer {
	return c.conn
}

func (c connectionErrorEvent) getError() error {
	return c.err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// PlcConnectionCacheCloseResult / plcConnectionCacheCloseResult
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type PlcConnectionCacheCloseResult interface {
	GetConnectionCache() PlcConnectionCache
	GetErr() error
}

type plcConnectionCacheCloseResult struct {
	connectionCache PlcConnectionCache
	err             error
}

func newDefaultPlcConnectionCacheCloseResult(connectionCache PlcConnectionCache, err error) PlcConnectionCacheCloseResult {
	return &plcConnectionCacheCloseResult{
		connectionCache: connectionCache,
		err:             err,
	}
}

func (p plcConnectionCacheCloseResult) GetConnectionCache() PlcConnectionCache {
	return p.connectionCache
}

func (p plcConnectionCacheCloseResult) GetErr() error {
	return p.err
}
