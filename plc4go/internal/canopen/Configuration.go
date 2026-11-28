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

package canopen

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
)

// Configuration is what a canopen connection string can say. Ported from plc4j's
// CANOpenConfiguration (node-id, heartbeat, request-timeout-ms), plus heartbeat-interval-ms, which
// plc4j fixes at 10 seconds.
type Configuration struct {
	// nodeID is the node id this connection uses on the bus: the producer id of its heartbeat. 0
	// means none, which is fine for a connection that only acts as an SDO client and listener.
	nodeID uint8
	// heartbeat makes the connection produce heartbeats: BOOTED_UP once, then OPERATIONAL every
	// heartbeatInterval.
	heartbeat         bool
	heartbeatInterval time.Duration
	// requestTimeout bounds each step of an SDO transfer: how long to wait for each answer.
	requestTimeout time.Duration
}

const (
	defaultRequestTimeout    = time.Second
	defaultHeartbeatInterval = 10 * time.Second
)

func DefaultConfiguration() Configuration {
	return Configuration{
		heartbeatInterval: defaultHeartbeatInterval,
		requestTimeout:    defaultRequestTimeout,
	}
}

// ParseFromOptions reads the connection options. node-id is also accepted as nodeId, the spelling
// the protocol's driver test suites use.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	reader := spiOptions.NewOptionReader(localLog, connectionOptions)
	defer reader.ReportUnknown("canopen")

	configuration := DefaultConfiguration()

	nodeIDString := reader.Get("node-id")
	if alternative := reader.Get("nodeId"); nodeIDString == "" {
		nodeIDString = alternative
	}
	if nodeIDString != "" {
		nodeID, err := strconv.ParseUint(nodeIDString, 10, 8)
		if err != nil || nodeID > maxNodeID {
			return Configuration{}, errors.Errorf("node-id %q is not a CANopen node id (0..%d)", nodeIDString, maxNodeID)
		}
		configuration.nodeID = uint8(nodeID)
	}

	if heartbeatString := reader.Get("heartbeat"); heartbeatString != "" {
		heartbeat, err := strconv.ParseBool(heartbeatString)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "error parsing heartbeat %q", heartbeatString)
		}
		configuration.heartbeat = heartbeat
	}

	var err error
	if configuration.heartbeatInterval, err = parseMillis(reader.Get("heartbeat-interval-ms"), "heartbeat-interval-ms", defaultHeartbeatInterval); err != nil {
		return Configuration{}, err
	}
	if configuration.requestTimeout, err = parseMillis(reader.Get("request-timeout-ms"), "request-timeout-ms", defaultRequestTimeout); err != nil {
		return Configuration{}, err
	}

	if configuration.heartbeat && configuration.nodeID == 0 {
		return Configuration{}, errors.New("heartbeat=true needs a node-id in 1..127 to produce the heartbeat as")
	}
	return configuration, nil
}

func parseMillis(text, name string, fallback time.Duration) (time.Duration, error) {
	if text == "" {
		return fallback, nil
	}
	millis, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "error parsing %s %q", name, text)
	}
	if millis == 0 {
		return 0, errors.Errorf("%s must be greater than zero", name)
	}
	return time.Duration(millis) * time.Millisecond, nil
}

func (c Configuration) String() string {
	return fmt.Sprintf("canopen.Configuration{nodeID: %d, heartbeat: %t, heartbeatInterval: %s, requestTimeout: %s}",
		c.nodeID, c.heartbeat, c.heartbeatInterval, c.requestTimeout)
}
