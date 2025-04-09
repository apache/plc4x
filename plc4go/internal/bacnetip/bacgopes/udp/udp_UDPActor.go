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

package udp

import (
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

//go:generate go tool plc4xGenerator -type=UDPActor -prefix=udp_
type UDPActor struct {
	director *UDPDirector `stringer:"true"`
	timeout  uint32
	timer    *OneShotFunctionTask `stringer:"true"`
	peer     string

	log zerolog.Logger
}

func NewUDPActor(localLog zerolog.Logger, director *UDPDirector, peer string) *UDPActor {
	a := &UDPActor{
		log: localLog,
	}

	// keep track of the director
	a.director = director

	// associated with a peer
	a.peer = peer

	// Add a timer
	a.timeout = director.timeout
	if a.timeout > 0 {
		a.timer = FunctionTask(a.idleTimeout, NoArgs, NoKWArgs())
		a.timer.InstallTask(WithInstallTaskOptionsWhen(GetTaskManagerTime().Add(time.Duration(a.timeout) * time.Millisecond)))
	}

	// tell the director this is a new actor
	a.director.AddActor(a)
	return a
}

func (a *UDPActor) idleTimeout(_ Args, _ KWArgs) error {
	a.log.Debug().Msg("idleTimeout")

	// tell the director this is gone
	a.director.DelActor(a)
	return nil
}

func (a *UDPActor) Indication(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	pdu := GA[PDU](args, 0)

	// reschedule the timer
	if a.timer != nil {
		a.timer.InstallTask(WithInstallTaskOptionsWhen(GetTaskManagerTime().Add(time.Duration(a.timeout) * time.Millisecond)))
	}

	// put it in the outbound queue for the director
	a.director.request <- pdu
	return nil
}

func (a *UDPActor) Response(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Response")

	// reschedule the timer
	if a.timer != nil {
		a.timer.InstallTask(WithInstallTaskOptionsWhen(GetTaskManagerTime().Add(time.Duration(a.timeout) * time.Millisecond)))
	}

	// process this as a response from the director
	return a.director.Response(args, kwArgs)
}

func (a *UDPActor) HandleError(err error) {
	a.log.Debug().Err(err).Msg("HandleError")

	if err != nil {
		a.director.ActorError(a, err)
	}
}
