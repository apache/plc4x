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

package bacnetip

import (
	"fmt"

	"github.com/rs/zerolog"
)

type RouterStatus uint8

const (
	ROUTER_AVAILABLE    RouterStatus = iota // normal
	ROUTER_BUSY                             // router is busy
	ROUTER_DISCONNECTED                     // could make a connection, but hasn't
	ROUTER_UNREACHABLE                      // temporarily unreachable
)

func (r RouterStatus) String() string {
	switch r {
	case ROUTER_AVAILABLE:
		return "ROUTER_AVAILABLE"
	case ROUTER_BUSY:
		return "ROUTER_BUSY"
	case ROUTER_DISCONNECTED:
		return "ROUTER_DISCONNECTED"
	case ROUTER_UNREACHABLE:
		return "ROUTER_UNREACHABLE"
	default:
		return "Unknown"
	}
}

type RouterInfo struct {
	snet    *uint16
	address *Address
	dnets   map[*uint16]RouterStatus
}

func (r RouterInfo) String() string {
	return fmt.Sprintf("%#q", r)
}

type snetDnetTuple struct {
	snet netKey
	dnet netKey
}

type RouterInfoCache struct {
	routers  map[netKey]*RouterInfo        // TODO: snet -> {Address: RouterInfo}
	pathInfo map[snetDnetTuple]*RouterInfo // TODO: (snet, dnet) -> RouterInfo

	log zerolog.Logger
}

func NewRouterInfoCache(localLog zerolog.Logger) *RouterInfoCache {
	localLog.Debug().Msg("NewRouterInfoCache")
	return &RouterInfoCache{
		routers:  map[netKey]*RouterInfo{},
		pathInfo: map[snetDnetTuple]*RouterInfo{},

		log: localLog,
	}
}

func (n *RouterInfoCache) GetRouterInfo(snet, dnet netKey) *RouterInfo {
	n.log.Debug().Stringer("snet", snet).Stringer("dnet", dnet).Msg("GetRouterInfo")

	// return the network and address
	routerInfo, _ := n.pathInfo[snetDnetTuple{snet, dnet}]
	return routerInfo
}

func (n *RouterInfoCache) UpdateRouterInfo(snet netKey, address *Address, dnets []uint16, status *RouterStatus) error {
	n.log.Debug().Stringer("snet", snet).Stringer("dnet", address).Uints16("dnets", dnets).Msg("UpdateRouterInfo")

	existingRouterInfo, _ := n.routers[snet] // TODO: what is happening here with the address

	var otherRouters []*RouterInfo
	for _, dnet := range dnets {
		otherRouter, _ := n.pathInfo[snetDnetTuple{snet, nk(&dnet)}]
		if otherRouter != nil && otherRouter != existingRouterInfo {
			otherRouters = append(otherRouters, otherRouter)
		}
	}

	// TODO: finish
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateRouterStatus(*uint16, *Address, []*uint16) {
	panic("not implemented yet")
}

func (n *RouterInfoCache) DeleteRouterInfo(*uint16, any, any) error {
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateSourceNetwork() {
	panic("not implemented yet")
}
