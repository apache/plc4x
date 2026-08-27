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

package serial

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
)

// serialConfig bundles everything the serial transport configures on a
// connection: the port settings handed to serialport.Open, the modem lines
// asserted after open, and the transport-level timeouts.
type serialConfig struct {
	port serialport.Config
	dtr  bool
	rts  bool
	// readTimeout bounds reads that carry no explicit context deadline;
	// 0 means blocking reads.
	readTimeout time.Duration
	// writeTimeout bounds writes that carry no explicit context deadline;
	// 0 means no write timeout.
	writeTimeout   time.Duration
	connectTimeout uint32 // milliseconds
	// reusePort shares one physical port between transport instances with
	// identical configuration (broadcast reads, serialized paced writes).
	reusePort bool
	// interframeDelay is the minimum silence enforced before each write,
	// measured from the later of last-write-end and last-received-byte.
	// 0 disables pacing. Works on shared and dedicated ports alike.
	interframeDelay time.Duration
}

func defaultSerialConfig() serialConfig {
	return serialConfig{
		port:           serialport.Config{BaudRate: 9600}, // 8N1 via serialport defaults
		readTimeout:    time.Second,
		writeTimeout:   time.Second,
		connectTimeout: 1000,
	}
}

// parseSerialOptions translates the connection-string options map into a
// serialConfig. Enum values are case-insensitive and accept "-" or "_" as
// separator; invalid values are hard errors naming the option; unknown
// options and empty value slices are ignored.
func parseSerialOptions(options map[string][]string) (serialConfig, error) {
	cfg := defaultSerialConfig()

	if raw, ok := firstValue(options, "baud-rate"); ok {
		value, err := strconv.ParseUint(raw, 10, 32)
		if err != nil || value == 0 {
			return cfg, optionError("baud-rate", raw, "must be a positive integer")
		}
		cfg.port.BaudRate = uint(value)
	}
	if raw, ok := firstValue(options, "data-bits"); ok {
		value, err := strconv.ParseUint(raw, 10, 8)
		if err != nil || value < 5 || value > 8 {
			return cfg, optionError("data-bits", raw, "must be 5..8")
		}
		cfg.port.DataBits = uint(value)
	}
	if raw, ok := firstValue(options, "stop-bits"); ok {
		switch raw {
		case "1":
			cfg.port.StopBits = serialport.StopBitsOne
		case "2":
			cfg.port.StopBits = serialport.StopBitsTwo
		default:
			return cfg, optionError("stop-bits", raw, "must be 1 or 2")
		}
	}
	if raw, ok := firstValue(options, "parity"); ok {
		switch normalizeEnum(raw) {
		case "none":
			cfg.port.Parity = serialport.ParityNone
		case "odd":
			cfg.port.Parity = serialport.ParityOdd
		case "even":
			cfg.port.Parity = serialport.ParityEven
		case "mark":
			cfg.port.Parity = serialport.ParityMark
		case "space":
			cfg.port.Parity = serialport.ParitySpace
		default:
			return cfg, optionError("parity", raw, "must be one of none, odd, even, mark, space")
		}
	}
	if raw, ok := firstValue(options, "flow-control"); ok {
		switch normalizeEnum(raw) {
		case "none":
			// defaults already off
		case "rts-cts":
			cfg.port.RTSCTSFlowControl = true
		case "xon-xoff":
			cfg.port.XONXOFFFlowControl = true
		default:
			return cfg, optionError("flow-control", raw, "must be one of none, rts-cts, xon-xoff")
		}
	}
	if raw, ok := firstValue(options, "dtr"); ok {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return cfg, optionError("dtr", raw, "must be true or false")
		}
		cfg.dtr = value
	}
	if raw, ok := firstValue(options, "rts"); ok {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return cfg, optionError("rts", raw, "must be true or false")
		}
		cfg.rts = value
	}
	if raw, ok := firstValue(options, "read-timeout-ms"); ok {
		millis, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return cfg, optionError("read-timeout-ms", raw, "must be a non-negative integer (milliseconds)")
		}
		cfg.readTimeout = time.Duration(millis) * time.Millisecond
	}
	if raw, ok := firstValue(options, "write-timeout-ms"); ok {
		millis, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return cfg, optionError("write-timeout-ms", raw, "must be a non-negative integer (milliseconds)")
		}
		cfg.writeTimeout = time.Duration(millis) * time.Millisecond
	}
	if raw, ok := firstValue(options, "connect-timeout-ms"); ok {
		millis, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return cfg, optionError("connect-timeout-ms", raw, "must be a non-negative integer (milliseconds)")
		}
		cfg.connectTimeout = uint32(millis)
	}
	if raw, ok := firstValue(options, "reuse-port"); ok {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return cfg, optionError("reuse-port", raw, "must be true or false")
		}
		cfg.reusePort = value
	}
	if raw, ok := firstValue(options, "interframe-delay"); ok {
		millis, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return cfg, optionError("interframe-delay", raw, "must be a non-negative integer (milliseconds)")
		}
		cfg.interframeDelay = time.Duration(millis) * time.Millisecond
	}
	return cfg, nil
}

func firstValue(options map[string][]string, key string) (string, bool) {
	values, ok := options[key]
	if !ok || len(values) == 0 {
		return "", false
	}
	return values[0], true
}

func normalizeEnum(value string) string {
	return strings.ReplaceAll(strings.ToLower(value), "_", "-")
}

func optionError(option, value, requirement string) error {
	return fmt.Errorf("error parsing option %q: invalid value %q (%s)", option, value, requirement)
}
