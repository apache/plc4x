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

package gsdml

// Index answers the questions a controller asks of a profile: what a text id says, which device
// access point a device has, and what a submodule the device reports in a slot looks like.
type Index struct {
	Profile *Profile

	texts      map[string]string
	submodules map[string]*SubmoduleItem
	values     map[string]*ValueItem
}

func NewIndex(profile *Profile) *Index {
	process := &profile.ProfileBody.ApplicationProcess
	index := &Index{
		Profile:    profile,
		texts:      make(map[string]string, len(process.ExternalTextList.PrimaryLanguage)),
		submodules: make(map[string]*SubmoduleItem, len(process.Submodules)),
		values:     make(map[string]*ValueItem, len(process.ValueItems)),
	}
	for _, text := range process.ExternalTextList.PrimaryLanguage {
		index.texts[text.TextID] = text.Value
	}
	for i := range process.Submodules {
		index.submodules[process.Submodules[i].ID] = &process.Submodules[i]
	}
	for i := range process.ValueItems {
		index.values[process.ValueItems[i].ID] = &process.ValueItems[i]
	}
	return index
}

// Text is the primary language text of a text id, or the id itself when the list has none.
func (x *Index) Text(textID string) string {
	if text, ok := x.texts[textID]; ok {
		return text
	}
	return textID
}

// ValueItem is a value list entry by id.
func (x *Index) ValueItem(id string) (*ValueItem, bool) {
	item, ok := x.values[id]
	return item, ok
}

// DeviceAccessPoint is a DAP by its ID attribute.
func (x *Index) DeviceAccessPoint(id string) (*DeviceAccessPointItem, bool) {
	daps := x.Profile.ProfileBody.ApplicationProcess.DeviceAccessPoints
	for i := range daps {
		if daps[i].ID == id {
			return &daps[i], true
		}
	}
	return nil, false
}

// DeviceAccessPointsByIdent are the DAPs with a module ident number; a device reports its DAP's
// ident in slot 0.
func (x *Index) DeviceAccessPointsByIdent(moduleIdent uint32) []*DeviceAccessPointItem {
	var matches []*DeviceAccessPointItem
	daps := x.Profile.ProfileBody.ApplicationProcess.DeviceAccessPoints
	for i := range daps {
		if identEquals(daps[i].ModuleIdentNumber, moduleIdent) {
			matches = append(matches, &daps[i])
		}
	}
	return matches
}

// Module is a module item by its ID attribute.
func (x *Index) Module(id string) (*ModuleItem, bool) {
	modules := x.Profile.ProfileBody.ApplicationProcess.Modules
	for i := range modules {
		if modules[i].ID == id {
			return &modules[i], true
		}
	}
	return nil, false
}

// SubmodulesOf are the submodules a module can carry: its virtual submodules, then the ones it
// references from the SubmoduleList.
func (x *Index) SubmodulesOf(module *ModuleItem) []*SubmoduleItem {
	var result []*SubmoduleItem
	for i := range module.VirtualSubmodules {
		result = append(result, &module.VirtualSubmodules[i])
	}
	for _, ref := range module.UseableSubmodules {
		if submodule, ok := x.submodules[ref.SubmoduleItemTarget]; ok {
			result = append(result, submodule)
		}
	}
	return result
}

// Resolved is what the GSDML says about a submodule a device reported: which module it belongs to
// and its description.
type Resolved struct {
	// DeviceAccessPoint is set when the submodule is part of a DAP (slot 0), Module otherwise.
	DeviceAccessPoint *DeviceAccessPointItem
	Module            *ModuleItem
	Submodule         *SubmoduleItem
}

// ModuleID is the ID attribute of the DAP or module.
func (r Resolved) ModuleID() string {
	if r.DeviceAccessPoint != nil {
		return r.DeviceAccessPoint.ID
	}
	return r.Module.ID
}

// Resolve finds the description of the submodule a device reports by its module and submodule
// ident numbers. Both are needed: modular stations reuse one module ident for many modules (WAGO's
// 750-606 and 750-610 are both 0x00008200) and tell them apart by submodule ident. Interface and
// port submodules of a DAP resolve to a submodule without I/O data. The first match wins.
func (x *Index) Resolve(moduleIdent, submoduleIdent uint32) (Resolved, bool) {
	process := &x.Profile.ProfileBody.ApplicationProcess
	for i := range process.DeviceAccessPoints {
		dap := &process.DeviceAccessPoints[i]
		if !identEquals(dap.ModuleIdentNumber, moduleIdent) {
			continue
		}
		for j := range dap.VirtualSubmodules {
			if identEquals(dap.VirtualSubmodules[j].SubmoduleIdentNumber, submoduleIdent) {
				return Resolved{DeviceAccessPoint: dap, Submodule: &dap.VirtualSubmodules[j]}, true
			}
		}
		for _, item := range dap.SystemDefined.Interfaces {
			if identEquals(item.SubmoduleIdentNumber, submoduleIdent) {
				return Resolved{DeviceAccessPoint: dap, Submodule: &SubmoduleItem{ID: item.ID, SubmoduleIdentNumber: item.SubmoduleIdentNumber}}, true
			}
		}
		for _, item := range dap.SystemDefined.Ports {
			if identEquals(item.SubmoduleIdentNumber, submoduleIdent) {
				return Resolved{DeviceAccessPoint: dap, Submodule: &SubmoduleItem{ID: item.ID, SubmoduleIdentNumber: item.SubmoduleIdentNumber}}, true
			}
		}
	}
	for i := range process.Modules {
		module := &process.Modules[i]
		if !identEquals(module.ModuleIdentNumber, moduleIdent) {
			continue
		}
		for _, submodule := range x.SubmodulesOf(module) {
			if identEquals(submodule.SubmoduleIdentNumber, submoduleIdent) {
				return Resolved{Module: module, Submodule: submodule}, true
			}
		}
	}
	return Resolved{}, false
}

func identEquals(text string, ident uint32) bool {
	value, err := ParseNumber(text)
	return err == nil && value == uint64(ident)
}
