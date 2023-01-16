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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.profinet.gsdml.ProfinetDeviceItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetInterfaceSubmoduleItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetPortSubmoduleItem;
import org.apache.plc4x.java.profinet.gsdml.ProfinetVirtualSubmoduleItem;
import org.apache.plc4x.java.profinet.readwrite.*;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class ProfinetModuleImpl implements ProfinetModule {

    private final ProfinetDeviceItem module;
    private final int ioCsOffset;
    private final int ioPsOffset;
    private final int slot;
    private List<PnIoCm_IoDataObject> inputIoPsApiBlocks = new ArrayList<>();
    private List<PnIoCm_IoCs> outputIoCsApiBlocks = new ArrayList<>();
    private List<PnIoCm_Submodule> expectedSubModuleApiBlocks = new ArrayList<>();

    private List<PnIoCm_IoCs> inputIoCsApiBlocks = new ArrayList<>();
    private List<PnIoCm_IoDataObject> outputIoDataApiBlocks = new ArrayList<>();
    private int ioPsSize;
    private int ioCsSize;

    public ProfinetModuleImpl(ProfinetDeviceItem module, int ioCsOffset, int ioPsOffset, int slot) {
        this.module = module;
        this.ioCsOffset = ioCsOffset;
        this.ioPsOffset = ioPsOffset;
        this.slot = slot;
    }

    private void populateNode() {
        int inputIoPsOffset = ioPsOffset;
        int outputIoCsOffset = ioCsOffset;

        for (ProfinetVirtualSubmoduleItem virtualItem : module.getVirtualSubmoduleList()) {
            Integer identNumber = Integer.decode(virtualItem.getSubmoduleIdentNumber());
            inputIoPsApiBlocks.add(new PnIoCm_IoDataObject(
                slot,
                identNumber,
                inputIoPsOffset));
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                slot,
                identNumber,
                outputIoCsOffset));
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                identNumber,
                identNumber,
                false,
                false,
                false,
                false));
            inputIoPsOffset += 1;
            outputIoCsOffset += 1;
        }

        for (ProfinetInterfaceSubmoduleItem interfaceItem : module.getSystemDefinedSubmoduleList().getInterfaceSubmodules()) {
            Integer identNumber = Integer.decode(interfaceItem.getSubmoduleIdentNumber());
            inputIoPsApiBlocks.add(new PnIoCm_IoDataObject(
                slot,
                identNumber,
                inputIoPsOffset));
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                slot,
                identNumber,
                outputIoCsOffset));
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                identNumber,
                identNumber,
                false,
                false,
                false,
                false));
            inputIoPsOffset += 1;
            outputIoCsOffset += 1;
        }

        for (
            ProfinetPortSubmoduleItem portItem : module.getSystemDefinedSubmoduleList().getPortSubmodules()) {
            Integer identNumber = Integer.decode(portItem.getSubmoduleIdentNumber());
            inputIoPsApiBlocks.add(new PnIoCm_IoDataObject(
                slot,
                identNumber,
                inputIoPsOffset));
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                slot,
                identNumber,
                outputIoCsOffset));
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                identNumber,
                identNumber,
                false,
                false,
                false,
                false));
            inputIoPsOffset += 1;
            outputIoCsOffset += 1;
        }




        ioPsSize = inputIoPsOffset - ioPsOffset;
        ioCsSize = outputIoCsOffset - ioCsOffset;
    }

    public int getInputIoPsSize() {
        return ioPsSize;
    }

    public int getOutputIoCsSize() {
        return ioCsSize;
    }
}
