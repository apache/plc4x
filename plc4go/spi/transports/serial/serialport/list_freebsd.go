//go:build freebsd

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

package serialport

import (
	"path/filepath"
	"sort"
	"strings"
)

// listPorts returns callout devices (cuau* for onboard UARTs, cuaU* for
// USB adapters), excluding the .init/.lock control nodes.
func listPorts() ([]string, error) {
	matches, err := filepath.Glob("/dev/cua[uU]*")
	if err != nil {
		return nil, err
	}
	ports := make([]string, 0, len(matches))
	for _, m := range matches {
		if strings.HasSuffix(m, ".init") || strings.HasSuffix(m, ".lock") {
			continue
		}
		ports = append(ports, m)
	}
	sort.Strings(ports)
	return ports, nil
}
