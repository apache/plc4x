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

package profinet

import (
	"sort"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/profinet/gsdml"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/profinet/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// minCyclicDataLength is the smallest C_SDU an RT frame carries; shorter data is padded to it.
const minCyclicDataLength = 40

// realSubmodule is a submodule the device reported in its RealIdentificationData.
type realSubmodule struct {
	Slot, Subslot               uint16
	ModuleIdent, SubmoduleIdent uint32
}

// plannedSubmodule is one submodule of the application relation and where its data sits in the two
// cyclic frames. Offsets are into the C_SDU; -1 means the submodule has nothing there.
type plannedSubmodule struct {
	realSubmodule
	Resolved gsdml.Resolved
	Layout   gsdml.Layout

	// In the input frame (device to controller): the submodule's input data and its IOPS, and the
	// device's consumer status (IOCS) for the submodule's outputs.
	InputData, InputIOPS, InputIOCS int
	// In the output frame (controller to device): the output data and its IOPS, and our consumer
	// status for the submodule's inputs.
	OutputData, OutputIOPS, OutputIOCS int
}

// HasInput is a submodule that provides data in the input frame; one without any I/O still sends
// its IOPS there, as a zero length data object.
func (s *plannedSubmodule) HasInput() bool {
	return s.Layout.InputLength > 0 || s.Layout.OutputLength == 0
}

func (s *plannedSubmodule) HasOutput() bool {
	return s.Layout.OutputLength > 0
}

// plan is the application relation: every submodule the device has, laid out in the frames.
type plan struct {
	Submodules   []*plannedSubmodule
	InputLength  int
	OutputLength int
}

// Find is the planned submodule in a slot and subslot.
func (p *plan) Find(slot, subslot uint16) (*plannedSubmodule, bool) {
	for _, s := range p.Submodules {
		if s.Slot == slot && s.Subslot == subslot {
			return s, true
		}
	}
	return nil, false
}

// realSubmodulesOf flattens a RealIdentificationData block.
func realSubmodulesOf(block readWriteModel.PnIoCm_Block_RealIdentificationData) []realSubmodule {
	var result []realSubmodule
	for _, api := range block.GetApis() {
		for _, slot := range api.GetSlots() {
			for _, subslot := range slot.GetSubslots() {
				result = append(result, realSubmodule{
					Slot: slot.GetSlotNumber(), Subslot: subslot.GetSubslotNumber(),
					ModuleIdent: slot.GetModuleIdentNumber(), SubmoduleIdent: subslot.GetSubmoduleIdentNumber(),
				})
			}
		}
	}
	return result
}

// makePlan resolves every reported submodule in the GSDML and lays out the frames the way the
// captured connect of a real controller does.
func makePlan(index *gsdml.Index, reported []realSubmodule) (*plan, error) {
	return planFrom(reported, func(real realSubmodule) (gsdml.Resolved, gsdml.Layout, error) {
		resolved, ok := index.Resolve(real.ModuleIdent, real.SubmoduleIdent)
		if !ok {
			return gsdml.Resolved{}, gsdml.Layout{}, errors.Errorf("the GSDML has no module 0x%08X with submodule 0x%08X (slot %d, subslot 0x%04X)",
				real.ModuleIdent, real.SubmoduleIdent, real.Slot, real.Subslot)
		}
		layout, err := resolved.Submodule.IOData.Layout()
		if err != nil {
			return gsdml.Resolved{}, gsdml.Layout{}, errors.Wrapf(err, "slot %d subslot 0x%04X", real.Slot, real.Subslot)
		}
		return resolved, layout, nil
	})
}

// planFrom lays out the frames for submodules whose descriptions resolve gives.
func planFrom(reported []realSubmodule, resolve func(realSubmodule) (gsdml.Resolved, gsdml.Layout, error)) (*plan, error) {
	if len(reported) == 0 {
		return nil, errors.New("the device reported no submodules")
	}
	sorted := append([]realSubmodule(nil), reported...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Slot != sorted[j].Slot {
			return sorted[i].Slot < sorted[j].Slot
		}
		return sorted[i].Subslot < sorted[j].Subslot
	})
	p := &plan{}
	for _, real := range sorted {
		resolved, layout, err := resolve(real)
		if err != nil {
			return nil, err
		}
		p.Submodules = append(p.Submodules, &plannedSubmodule{
			realSubmodule: real, Resolved: resolved, Layout: layout,
			InputData: -1, InputIOPS: -1, InputIOCS: -1, OutputData: -1, OutputIOPS: -1, OutputIOCS: -1,
		})
	}

	// Per submodule, in slot and subslot order: in the input frame the data object (data, then
	// IOPS) before the IOCS for the submodule's outputs; in the output frame the IOCS for the
	// submodule's inputs before its data object. That is the layout the captured controller chose.
	in, out := 0, 0
	for _, s := range p.Submodules {
		if s.HasInput() {
			s.InputData = in
			in += s.Layout.InputLength
			s.InputIOPS = in
			in += s.Layout.IOPSLength
		}
		if s.HasOutput() {
			s.InputIOCS = in
			in += s.Layout.IOCSLength
		}
		if s.HasInput() {
			s.OutputIOCS = out
			out += s.Layout.IOCSLength
		}
		if s.HasOutput() {
			s.OutputData = out
			out += s.Layout.OutputLength
			s.OutputIOPS = out
			out += s.Layout.IOPSLength
		}
	}
	p.InputLength = max(in, minCyclicDataLength)
	p.OutputLength = max(out, minCyclicDataLength)
	if p.InputLength > 1440 || p.OutputLength > 1440 {
		return nil, errors.Errorf("the cyclic data needs %d input and %d output bytes, more than one RT frame carries (1440)", in, out)
	}
	return p, nil
}

// ioCrApi lists the data objects and consumer statuses of one frame direction.
func (p *plan) ioCrApi(input bool) readWriteModel.PnIoCm_IoCrBlockReqApi {
	var objects []readWriteModel.PnIoCm_IoDataObject
	var statuses []readWriteModel.PnIoCm_IoCs
	for _, s := range p.Submodules {
		if input {
			if s.InputData >= 0 {
				objects = append(objects, readWriteModel.NewPnIoCm_IoDataObject(s.Slot, s.Subslot, uint16(s.InputData)))
			}
			if s.InputIOCS >= 0 {
				statuses = append(statuses, readWriteModel.NewPnIoCm_IoCs(s.Slot, s.Subslot, uint16(s.InputIOCS)))
			}
		} else {
			if s.OutputData >= 0 {
				objects = append(objects, readWriteModel.NewPnIoCm_IoDataObject(s.Slot, s.Subslot, uint16(s.OutputData)))
			}
			if s.OutputIOCS >= 0 {
				statuses = append(statuses, readWriteModel.NewPnIoCm_IoCs(s.Slot, s.Subslot, uint16(s.OutputIOCS)))
			}
		}
	}
	return readWriteModel.NewPnIoCm_IoCrBlockReqApi(objects, statuses)
}

// expectedSubmoduleBlocks are the ExpectedSubmoduleBlockReqs, one per slot, telling the device what
// the controller expects in it: what the device reported, with the data lengths from the GSDML.
func (p *plan) expectedSubmoduleBlocks() []readWriteModel.PnIoCm_Block {
	var blocks []readWriteModel.PnIoCm_Block
	for i := 0; i < len(p.Submodules); {
		slot := p.Submodules[i].Slot
		var submodules []readWriteModel.PnIoCm_Submodule
		moduleIdent := p.Submodules[i].ModuleIdent
		for ; i < len(p.Submodules) && p.Submodules[i].Slot == slot; i++ {
			s := p.Submodules[i]
			switch {
			case s.Layout.InputLength > 0 && s.Layout.OutputLength > 0:
				submodules = append(submodules, readWriteModel.NewPnIoCm_Submodule_InputAndOutputData(s.Subslot, s.SubmoduleIdent, false, false, false, false, uint16(s.Layout.InputLength), uint16(s.Layout.OutputLength)))
			case s.Layout.InputLength > 0:
				submodules = append(submodules, readWriteModel.NewPnIoCm_Submodule_InputData(s.Subslot, s.SubmoduleIdent, false, false, false, false, uint16(s.Layout.InputLength)))
			case s.Layout.OutputLength > 0:
				submodules = append(submodules, readWriteModel.NewPnIoCm_Submodule_OutputData(s.Subslot, s.SubmoduleIdent, false, false, false, false, uint16(s.Layout.OutputLength)))
			default:
				submodules = append(submodules, readWriteModel.NewPnIoCm_Submodule_NoInputNoOutputData(s.Subslot, s.SubmoduleIdent, false, false, false, false))
			}
		}
		blocks = append(blocks, readWriteModel.NewPnIoCm_Block_ExpectedSubmoduleReq(blockVersionHigh, blockVersionLow,
			[]readWriteModel.PnIoCm_ExpectedSubmoduleBlockReqApi{readWriteModel.NewPnIoCm_ExpectedSubmoduleBlockReqApi(slot, moduleIdent, 0, submodules)}))
	}
	return blocks
}

// withConfiguredModules puts the configured modules in their slots, replacing whatever the device
// reported there. A module's submodules go in subslot 1 and up, or in their FixedInSubslots.
func withConfiguredModules(index *gsdml.Index, reported []realSubmodule, modules map[uint16]string) ([]realSubmodule, error) {
	if len(modules) == 0 {
		return reported, nil
	}
	var result []realSubmodule
	for _, r := range reported {
		if _, replaced := modules[r.Slot]; !replaced {
			result = append(result, r)
		}
	}
	for slot, id := range modules {
		module, ok := index.Module(id)
		if !ok {
			return nil, errors.Errorf("the GSDML has no module %q for slot %d", id, slot)
		}
		moduleIdent, err := gsdml.ParseNumber(module.ModuleIdentNumber)
		if err != nil {
			return nil, errors.Errorf("module %s has a bad ModuleIdentNumber %q", id, module.ModuleIdentNumber)
		}
		submodules := index.SubmodulesOf(module)
		if len(submodules) == 0 {
			return nil, errors.Errorf("module %s has no submodules", id)
		}
		for i, submodule := range submodules {
			subslot := uint16(i + 1)
			if fields := strings.Fields(submodule.FixedInSubslots); len(fields) > 0 {
				if fixed, err := gsdml.ParseNumber(fields[0]); err == nil {
					subslot = uint16(fixed)
				}
			}
			submoduleIdent, err := gsdml.ParseNumber(submodule.SubmoduleIdentNumber)
			if err != nil {
				return nil, errors.Errorf("submodule %s has a bad SubmoduleIdentNumber %q", submodule.ID, submodule.SubmoduleIdentNumber)
			}
			result = append(result, realSubmodule{Slot: slot, Subslot: subslot, ModuleIdent: uint32(moduleIdent), SubmoduleIdent: uint32(submoduleIdent)})
		}
	}
	return result, nil
}
