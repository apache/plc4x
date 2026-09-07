//go:build windows

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
	"fmt"
	"sort"

	"golang.org/x/sys/windows/registry"
)

// listPorts reads HKLM\HARDWARE\DEVICEMAP\SERIALCOMM, the canonical
// registry map of active serial devices to COM names.
func listPorts() ([]string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\DEVICEMAP\SERIALCOMM`, registry.QUERY_VALUE)
	if err == registry.ErrNotExist {
		return nil, nil // no serial devices present
	}
	if err != nil {
		return nil, fmt.Errorf("serialport: opening SERIALCOMM registry key: %w", err)
	}
	defer key.Close()
	names, err := key.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("serialport: reading SERIALCOMM values: %w", err)
	}
	ports := make([]string, 0, len(names))
	for _, name := range names {
		value, _, err := key.GetStringValue(name)
		if err != nil {
			return nil, fmt.Errorf("serialport: reading SERIALCOMM value %q: %w", name, err)
		}
		ports = append(ports, value)
	}
	sort.Strings(ports)
	return ports, nil
}
