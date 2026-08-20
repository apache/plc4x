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

package iec608705104

import (
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

const (
	// defaultRequestTimeout is plc4j's Iec608705014Configuration request-timeout default of four
	// seconds. It bounds the two handshake round trips (test frame and start-data-transfer) which
	// are the only request/response interactions the protocol has - everything after them is
	// unsolicited.
	defaultRequestTimeout = 4000 * time.Millisecond
	// defaultAckThreshold is the IEC 60870-5-104 'w' parameter: after this many unacknowledged
	// I-format frames the controlling station has to send an S-format acknowledgement. The standard
	// recommends 8 (two thirds of the default k of 12) and plc4j's Iec60870Connection hard-codes
	// exactly that.
	defaultAckThreshold = 8
	// maxAckThreshold is the largest window the sequence number can express: the receive sequence
	// number is 15 bits wide.
	maxAckThreshold = 32767
)

// Configuration is what an IEC 60870-5-104 connection string can say about the connection as a
// whole. plc4j's Iec608705014Configuration knows only request-timeout; the acknowledgement window
// is added here because it is a real protocol parameter ('w') that plc4j buried in a constant, and
// a station which insists on a smaller window otherwise drops the connection.
type Configuration struct {
	// requestTimeout bounds how long each of the two handshake round trips waits for its
	// confirmation.
	requestTimeout time.Duration
	// ackThreshold is how many I-format frames may go unacknowledged before an S-format
	// acknowledgement is sent ('w' in the standard).
	ackThreshold uint16
}

// DefaultConfiguration is a connection without any options set.
func DefaultConfiguration() Configuration {
	return Configuration{
		requestTimeout: defaultRequestTimeout,
		ackThreshold:   defaultAckThreshold,
	}
}

// ParseFromOptions reads the connection options out of a parsed connection string. The timeout is
// spelled in milliseconds, the way plc4j's @IntDefaultValue(4000) does.
func ParseFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string) (Configuration, error) {
	configuration := DefaultConfiguration()

	if requestTimeoutString := getFromOptions(localLog, connectionOptions, "request-timeout"); requestTimeoutString != "" {
		parsedInt, err := strconv.ParseUint(requestTimeoutString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing request-timeout %s", requestTimeoutString)
		}
		if parsedInt == 0 {
			return Configuration{}, errors.New("request-timeout must be greater than zero")
		}
		configuration.requestTimeout = time.Duration(parsedInt) * time.Millisecond
	}

	if ackThresholdString := getFromOptions(localLog, connectionOptions, "ack-threshold"); ackThresholdString != "" {
		parsedInt, err := strconv.ParseUint(ackThresholdString, 10, 32)
		if err != nil {
			return Configuration{}, errors.Wrapf(err, "Error parsing ack-threshold %s", ackThresholdString)
		}
		if parsedInt == 0 {
			return Configuration{}, errors.New("ack-threshold must be greater than zero")
		}
		if parsedInt > maxAckThreshold {
			return Configuration{}, errors.Errorf("ack-threshold may not be larger than %d. Was %d", maxAckThreshold, parsedInt)
		}
		configuration.ackThreshold = uint16(parsedInt)
	}

	return configuration, nil
}

// getFromOptions plucks a single-valued option out of the parsed connection string.
func getFromOptions(localLog zerolog.Logger, connectionOptions map[string][]string, key string) string {
	if optionValues, ok := connectionOptions[key]; ok {
		if len(optionValues) <= 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLog.Warn().Str("key", key).Msg("Option must be unique")
		}
		return optionValues[0]
	}
	return ""
}
