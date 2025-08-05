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

package netservice

import (
	"maps"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
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

//go:generate go tool plc4xGenerator -type=Router -prefix=netservice_
type Router struct {
	addresses map[string]*RouterInfo // TODO: this is a address but using pointer as key is bad
}

//go:generate go tool plc4xGenerator -type=RouterInfo -prefix=netservice_
type RouterInfo struct {
	snet    netKey   `stringer:"true"`
	address *Address `stringer:"true"`
	dnets   map[netKey]RouterStatus
	status  RouterStatus `stringer:"true"`
}

type snetDnetTuple struct {
	snet netKey
	dnet netKey
}

//go:generate go tool plc4xGenerator -type=RouterInfoCache -prefix=netservice_
type RouterInfoCache struct {
	routers  map[netKey]*Router
	pathInfo map[snetDnetTuple]*RouterInfo

	log zerolog.Logger
}

func NewRouterInfoCache(localLog zerolog.Logger) *RouterInfoCache {
	localLog.Trace().Msg("NewRouterInfoCache")
	if _debug != nil {
		_debug("__init__")
	}
	return &RouterInfoCache{
		routers:  map[netKey]*Router{},
		pathInfo: map[snetDnetTuple]*RouterInfo{},

		log: localLog,
	}
}

func (n *RouterInfoCache) GetRouterInfo(snet, dnet netKey) *RouterInfo {
	n.log.Debug().Stringer("snet", snet).Stringer("dnet", dnet).Msg("GetRouterInfo")
	if _debug != nil {
		_debug("get_routerInfo %r %r", snet, dnet)
	}

	// return the network and address
	routerInfo, _ := n.pathInfo[snetDnetTuple{snet, dnet}]
	if _debug != nil {
		_debug("   - routerInfo: %r", routerInfo)
	}
	return routerInfo
}

func (n *RouterInfoCache) UpdateRouterInfo(snet netKey, address *Address, dnets []uint16, status *RouterStatus) error {
	if status == nil {
		status = ToPtr(ROUTER_AVAILABLE)
	}
	n.log.Debug().Stringer("snet", snet).Stringer("dnet", address).Uints16("dnets", dnets).Msg("UpdateRouterInfo")
	if _debug != nil {
		_debug("update_routerInfo %r %r %r", snet, address, dnets)
	}

	var existingRouterInfo *RouterInfo
	if r, ok := n.routers[snet]; ok {
		existingRouterInfo = r.addresses[address.String()]
	}

	var otherRouters map[string]*RouterInfo
	for _, dnet := range dnets {
		otherRouter, _ := n.pathInfo[snetDnetTuple{snet, nk(&dnet)}]
		if otherRouter != nil && otherRouter != existingRouterInfo {
			otherRouters[otherRouter.String()] = otherRouter
		}
	}

	// remove the dnets from other routers and paths
	for _, routerInfo := range otherRouters {
		for _, dnet := range dnets {
			if _, ok := routerInfo.dnets[nk(&dnet)]; ok {
				delete(routerInfo.dnets, nk(&dnet))
				delete(n.pathInfo, snetDnetTuple{snet, nk(&dnet)})
				n.log.Debug().
					Stringer("snet", snet).
					Uint16("dnet", dnet).
					Stringer("routerInfoAddress", routerInfo.address).
					Msg("del path: snet -> dnet via routerInfoAddress")
				if _debug != nil {
					_debug("    - del path: %r -> %r via %r", snet, dnet, routerInfo.address)
				}
			}
		}
		if len(routerInfo.dnets) == 0 {
			delete(n.routers[snet].addresses, routerInfo.address.String())
			n.log.Debug().
				Stringer("snet", snet).
				Stringer("routerInfoAddress", routerInfo.address).
				Msg("no dnets: snet via routerInfoAddress")
			if _debug != nil {
				_debug("    - no dnets: %r via %r", snet, routerInfo.address)
			}
		}
	}

	// update current router info if there is one
	if existingRouterInfo == nil {
		routerInfo := &RouterInfo{snet: snet, address: address, dnets: make(map[netKey]RouterStatus)}
		if _, ok := n.routers[snet]; !ok {
			n.routers[snet] = &Router{addresses: map[string]*RouterInfo{
				address.String(): routerInfo,
			}}
		} else {
			n.routers[snet].addresses[address.String()] = routerInfo
		}

		for _, dnet := range dnets {
			n.pathInfo[snetDnetTuple{snet, nk(&dnet)}] = routerInfo
			n.log.Debug().
				Stringer("snet", snet).
				Uint16("dnet", dnet).
				Stringer("routerInfoAddress", routerInfo.address).
				Msg("add path: snet -> dnet via routerInfoAddress")
			if _debug != nil {
				_debug("    - add path: %r -> %r via %r", snet, dnet, routerInfo.address)
			}
			routerInfo.dnets[nk(&dnet)] = *status
		}
	} else {
		for _, dnet := range dnets {
			if _, ok := existingRouterInfo.dnets[nk(&dnet)]; !ok {
				n.pathInfo[snetDnetTuple{snet, nk(&dnet)}] = existingRouterInfo
				n.log.Info().
					Stringer("snet", snet).
					Uint16("dnet", dnet).
					Msg("add path: snet -> dnet")
				if _debug != nil {
					_debug("    - add path: %r -> %r", snet, dnet)
				}
			}
			existingRouterInfo.dnets[nk(&dnet)] = *status
		}
	}
	return nil
}

func (n *RouterInfoCache) UpdateRouterStatus(snet netKey, address *Address, status RouterStatus) error {
	n.log.Debug().Stringer("snet", snet).Stringer("dnet", address).Stringer("status", status).Msg("UpdateRouterStatus")
	if _debug != nil {
		_debug("update_router_status %r %r %r", snet, address, status)
	}

	var existingRouterInfo *RouterInfo
	if r, ok := n.routers[snet]; ok {
		existingRouterInfo = r.addresses[address.String()]
	}

	if existingRouterInfo == nil {
		n.log.Trace().Msg("not a router we know about")
		if _debug != nil {
			_debug("    - not a router we know about")
		}
		return nil
	}

	existingRouterInfo.status = status
	n.log.Trace().Msg("status updated")
	if _debug != nil {
		_debug("    - status updated")
	}
	return nil
}

func (n *RouterInfoCache) DeleteRouterInfo(snet netKey, address *Address, dnets []uint16) error {
	n.log.Debug().Stringer("snet", snet).Stringer("dnet", address).Uints16("dnets", dnets).Msg("DeleteRouterInfo")
	if _debug != nil {
		_debug("delete_routerInfo %r %r %r", dnets)
	}
	if address == nil && len(dnets) == 0 {
		return errors.New("inconsistent parameters")
	}

	var routerInfo *RouterInfo
	// remove the dnets from a router of the whole router
	if address != nil {
		if r, ok := n.routers[snet]; ok {
			routerInfo = r.addresses[address.String()]
		}
		if routerInfo == nil {
			n.log.Trace().Msg("no router info")
			if _debug != nil {
				_debug("    - no route info")
			}
		} else {
			for dnet := range routerInfo.dnets {
				if !dnet.IsNil {
					dnets = append(dnets, dnet.Value)
				}
			}
			for _, dnet := range dnets {
				n.log.Debug().
					Stringer("snet", snet).
					Uint16("dnet", dnet).
					Stringer("routerInfoAddress", routerInfo.address).
					Msg("del path: snet -> dnet via routerInfoAddress")
				if _debug != nil {
					_debug("    - del path: %r -> %r via %r", snet, dnet, routerInfo.address)
				}
			}
		}
		return nil
	}

	// look for routers to the dnets
	var otherRouters map[string]*RouterInfo
	for _, dnet := range dnets {
		otherRouter, _ := n.pathInfo[snetDnetTuple{snet, nk(&dnet)}]
		if otherRouter != nil {
			otherRouters[otherRouter.String()] = otherRouter
		}
	}

	// remove the dnets from other routers and paths
	for _, routerInfo := range otherRouters {
		for _, dnet := range dnets {
			if _, ok := routerInfo.dnets[nk(&dnet)]; ok {
				delete(routerInfo.dnets, nk(&dnet))
				delete(n.pathInfo, snetDnetTuple{snet, nk(&dnet)})
				n.log.Debug().
					Stringer("snet", snet).
					Uint16("dnet", dnet).
					Stringer("routerInfoAddress", routerInfo.address).
					Msg("del path: snet -> dnet via routerInfoAddress")
				if _debug != nil {
					_debug("    - del path: %r -> %r via %r", snet, dnet, routerInfo.address)
				}
			}
		}
		if len(routerInfo.dnets) == 0 {
			delete(n.routers[snet].addresses, routerInfo.address.String())
			n.log.Debug().
				Stringer("snet", snet).
				Stringer("routerInfoAddress", routerInfo.address).
				Msg("no dnets: snet via routerInfoAddress")
			if _debug != nil {
				_debug("    - no dnets: %r via %r", snet, routerInfo.address)
			}
		}
	}
	return nil
}

func (n *RouterInfoCache) UpdateSourceNetwork(oldSnet netKey, newSnet netKey) error {
	n.log.Debug().Stringer("oldSnet", oldSnet).Stringer("newSnet", newSnet).Msg("UpdateSourceNetwork")
	if _debug != nil {
		_debug("update_source_network %r %r", oldSnet, newSnet)
	}

	if _, ok := n.routers[oldSnet]; !ok {
		n.log.Debug().Interface("routers", n.routers).Msg("No router references")
		if _debug != nil {
			_debug("    - no router references: %r", maps.Keys(n.routers))
		}
		return nil
	}

	// move the route info to the new set
	n.routers[newSnet] = n.routers[oldSnet]
	delete(n.routers, oldSnet)
	snetRouters := n.routers[newSnet]

	// update the paths
	for _, routerInfo := range snetRouters.addresses {
		for dnet := range routerInfo.dnets {
			n.pathInfo[snetDnetTuple{newSnet, dnet}] = n.pathInfo[snetDnetTuple{oldSnet, newSnet}]
			delete(n.pathInfo, snetDnetTuple{oldSnet, newSnet})
		}
	}
	return nil
}
