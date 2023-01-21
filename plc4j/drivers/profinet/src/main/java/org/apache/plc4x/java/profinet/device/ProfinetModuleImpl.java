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

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.profinet.gsdml.*;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.spi.messages.DefaultPlcBrowseItem;
import org.apache.plc4x.java.spi.values.PlcSTRING;

import java.util.*;

public class ProfinetModuleImpl implements ProfinetModule {

    private final ProfinetDeviceItem module;
    private List<PnIoCm_Block_ExpectedSubmoduleReq> expectedSubmoduleReq = new ArrayList<>();
    private final int ioCsOffset;
    private final int ioPsOffset;
    private final int slot;
    private List<PnIoCm_IoDataObject> inputIoPsApiBlocks = new ArrayList<>();
    private List<PnIoCm_IoCs> outputIoCsApiBlocks = new ArrayList<>();
    private List<PnIoCm_Submodule> expectedSubModuleApiBlocks = new ArrayList<>();

    private List<PnIoCm_IoCs> inputIoCsApiBlocks = new ArrayList<>();
    private List<PnIoCm_IoDataObject> outputIoPsApiBlocks = new ArrayList<>();
    private int ioPsSize;
    private int ioCsSize;

    public ProfinetModuleImpl(ProfinetDeviceItem module, int ioCsOffset, int ioPsOffset, int slot) {
        this.module = module;
        this.ioCsOffset = ioCsOffset;
        this.ioPsOffset = ioPsOffset;
        this.slot = slot;
        populateNode();
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
            expectedSubModuleApiBlocks.addAll(populateExpectedSubModuleApiBlocks());
            inputIoPsOffset += 1;
            outputIoCsOffset += 1;
        }

        if (module.getSystemDefinedSubmoduleList() != null) {
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
                    0,
                    identNumber,
                    inputIoPsOffset));
                outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                    0,
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
        }


        ioPsSize = inputIoPsOffset - ioPsOffset;
        ioCsSize = outputIoCsOffset - ioCsOffset;

        Integer identNumber = Integer.decode(module.getModuleIdentNumber());
        if (module.getInputDataLength() != 0) {
            inputIoPsApiBlocks.add(new PnIoCm_IoDataObject(
                slot,
                0x01,
                inputIoPsOffset));
            inputIoPsOffset += 1 + module.getInputDataLength();
        }
        if (module.getInputDataLength() != 0) {
            outputIoCsApiBlocks.add(new PnIoCm_IoCs(
                slot,
                0x01,
                outputIoCsOffset));
            outputIoCsOffset += 1;
        }

        if (module.getOutputDataLength() != 0) {
            inputIoCsApiBlocks.add(new PnIoCm_IoCs(
                slot,
                0x01,
                inputIoPsOffset));
            inputIoPsOffset += module.getOutputDataLength();
        }

        if (module.getOutputDataLength() != 0) {
            outputIoPsApiBlocks.add(new PnIoCm_IoDataObject(
                slot,
                0x01,
                outputIoCsOffset));
            outputIoCsOffset += 1 + module.getOutputDataLength();
        }
    }

    private List<PnIoCm_Submodule> populateExpectedSubModuleApiBlocks() {
        List<PnIoCm_Submodule> expectedSubModuleApiBlocks = new ArrayList<>();
        if (module.getInputDataLength() != 0 && module.getOutputDataLength() != 0) {
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_InputAndOutputData(
                0x01,
                Long.decode(module.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                false,
                false,
                false,
                false,
                module.getInputDataLength(),
                (short) 0x01,
                (short) 0x01,
                module.getOutputDataLength(),
                (short) 0x01,
                (short) 0x01
            ));
        } else if (module.getInputDataLength() != 0) {
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_InputData(
                0x01,
                Long.decode(module.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                false,
                false,
                false,
                false,
                module.getInputDataLength(),
                (short) 0x01,
                (short) 0x01));


        } else if (module.getOutputDataLength() != 0) {
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_OutputData(
                0x01,
                Long.decode(module.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                false,
                false,
                false,
                false,
                module.getOutputDataLength(),
                (short) 0x01,
                (short) 0x01));
        } else if (module.getInputDataLength() == 0 && module.getOutputDataLength() == 0) {
            expectedSubModuleApiBlocks.add(new PnIoCm_Submodule_NoInputNoOutputData(
                0x01,
                Long.decode(module.getVirtualSubmoduleList().get(0).getSubmoduleIdentNumber()),
                false,
                false,
                false,
                false
            ));
        }
        return expectedSubModuleApiBlocks;
    }

    public int getInputIoPsSize() {
        return ioPsSize;
    }

    public int getOutputIoCsSize() {
        return ioCsSize;
    }

    @Override
    public List<PnIoCm_Submodule> getExpectedSubModuleApiBlocks() {
        return expectedSubModuleApiBlocks;
    }

    @Override
    public List<PnIoCm_IoDataObject> getInputIoPsApiBlocks() {
        return inputIoPsApiBlocks;
    }

    @Override
    public List<PnIoCm_IoCs> getOutputIoCsApiBlocks() {
        return outputIoCsApiBlocks;
    }

    @Override
    public List<PnIoCm_IoCs> getInputIoCsApiBlocks() {
        return inputIoCsApiBlocks;
    }

    @Override
    public List<PnIoCm_IoDataObject> getOutputIoPsApiBlocks() {
        return outputIoPsApiBlocks;
    }

    @Override
    public Integer getIdentNumber() {
        return Integer.decode(module.getModuleIdentNumber());
    }

    @Override
    public Integer getSlotNumber() {
        return slot;
    }

    @Override
    public Map<String, List<PlcBrowseItem>> browseTags(Map<String, List<PlcBrowseItem>> browseItems, String addressSpace, Map<String, PlcValue> options) {
        for (PnIoCm_IoDataObject block : inputIoPsApiBlocks) {
            int identNumber = block.getSubSlotNumber();
            for (ProfinetVirtualSubmoduleItem virtual : module.getVirtualSubmoduleList()) {
                if (identNumber == Integer.decode(virtual.getSubmoduleIdentNumber())) {
                    if (virtual.getModuleInfo().getName() != null) {
                        options.put("module_name", new PlcSTRING(virtual.getModuleInfo().getName().getTextId()));
                    }
                    if (virtual.getModuleInfo().getName() != null) {
                        options.put("module_info_text", new PlcSTRING(virtual.getModuleInfo().getInfoText().getTextId()));
                    }

                    String statusName = addressSpace + "." + this.slot + "." + block.getSubSlotNumber() + "." + virtual.getId() + ".Status";
                    browseItems.put(statusName, Collections.singletonList(new DefaultPlcBrowseItem(ProfinetTag.of(statusName + ":INT"), statusName, false, false, true, new HashMap<>(), options)));
                    if (virtual.getIoData() != null && virtual.getIoData().getInput() != null) {
                        for (ProfinetIoDataInput input : virtual.getIoData().getInput()) {
                            for (ProfinetDataItem item : input.getDataItemList()) {
                                if (item.isUseAsBits()) {
                                    for (ProfinetBitDataItem bitItem : item.getBitDataItem()) {
                                        String tagName = addressSpace + "." + this.slot + "." + block.getSubSlotNumber() + "." + item.getTextId() + "." + bitItem.getBitOffset();
                                        browseItems.put(tagName, Collections.singletonList(new DefaultPlcBrowseItem(ProfinetTag.of(tagName + ":BOOL"), tagName, false, false, true, new HashMap<>(), options)));
                                    }
                                } else {
                                    String tagName = addressSpace + "." + this.slot + "." + block.getSubSlotNumber() + "." + item.getTextId();
                                    String datatype = ProfinetDataType.firstEnumForFieldConversion(item.getDataType().toUpperCase()).toString();
                                    browseItems.put(tagName, Collections.singletonList(new DefaultPlcBrowseItem(ProfinetTag.of(tagName + ":" + datatype), tagName, false, false, true, new HashMap<>(), options)));
                                }
                            }
                        }
                    }
                }
            }
            if (module.getSystemDefinedSubmoduleList() != null) {
                for (ProfinetInterfaceSubmoduleItem systemInterface : module.getSystemDefinedSubmoduleList().getInterfaceSubmodules()) {
                    if (identNumber == Integer.decode(systemInterface.getSubmoduleIdentNumber())) {
                        String statusName = addressSpace + "." + this.slot + "." + block.getSubSlotNumber() + "." + systemInterface.getId() + ".Status";
                        browseItems.put(statusName, Collections.singletonList(new DefaultPlcBrowseItem(ProfinetTag.of(statusName + ":INT"), statusName, false, false, true, new HashMap<>(), options)));
                    }
                }
                for (ProfinetPortSubmoduleItem systemPort : module.getSystemDefinedSubmoduleList().getPortSubmodules()) {
                    if (identNumber == Integer.decode(systemPort.getSubmoduleIdentNumber())) {
                        String statusName = addressSpace + "." + this.slot + "." + block.getSubSlotNumber() + "." + systemPort.getId() + ".Status";
                        browseItems.put(statusName, Collections.singletonList(new DefaultPlcBrowseItem(ProfinetTag.of(statusName + ":INT"), statusName, false, false, true, new HashMap<>(), options)));
                    }
                }
            }
        }

        return browseItems;
    }
}
