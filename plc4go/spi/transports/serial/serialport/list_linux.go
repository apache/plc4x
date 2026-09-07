//go:build linux

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
)

// listPorts enumerates /sys/class/tty: entries with a device/ link are
// backed by real hardware (or USB adapters); bare line disciplines and
// virtual consoles have none.
func listPorts() ([]string, error) {
	matches, err := filepath.Glob("/sys/class/tty/*/device")
	if err != nil {
		return nil, err
	}
	ports := make([]string, 0, len(matches))
	for _, m := range matches {
		ports = append(ports, "/dev/"+filepath.Base(filepath.Dir(m)))
	}
	sort.Strings(ports)
	return ports, nil
}
