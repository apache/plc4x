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

package org.apache.plc4x.java.profinet.gsdml;

import org.apache.plc4x.java.profinet.readwrite.ProfinetDataType;
import org.w3c.dom.Element;

import java.io.IOException;
import java.util.List;

public class ProfinetModuleItem implements ProfinetDeviceItem {

    private String id;

    private String moduleIdentNumber;

    private ProfinetModuleInfo moduleInfo;

    private List<ProfinetVirtualSubmoduleItem> virtualSubmoduleList;

    public String getId() {
        return id;
    }

    public String getModuleIdentNumber() {
        return moduleIdentNumber;
    }

    public ProfinetModuleInfo getModuleInfo() {
        return moduleInfo;
    }

    public List<ProfinetVirtualSubmoduleItem> getVirtualSubmoduleList() {
        return virtualSubmoduleList;
    }

    @Override
    public ProfinetSystemDefinedSubmoduleList getSystemDefinedSubmoduleList() {
        return null;
    }

    private static int getLengthSimpleType(String dataType) {
        ProfinetDataType dt = ProfinetDataType.firstEnumForFieldConversion(dataType.toUpperCase());
        if (dt == null) {
            throw new UnsupportedOperationException("Data type " + dataType + " not supported");
        }
        return dt.getDataTypeSize();
    }

    public Integer getInputDataLength() {
        int length = 0;
        for (ProfinetVirtualSubmoduleItem module : this.virtualSubmoduleList) {
            if (module.getIoData() != null && module.getIoData().getInput() != null) {
                for (ProfinetIoDataInput inputIoData : module.getIoData().getInput()) {
                    for (ProfinetDataItem dataItem : inputIoData.getDataItemList()) {
                        length += ProfinetModuleItem.getLengthSimpleType(dataItem.getDataType());
                    }
                }
            }
        }
        return length;
    }

    public Integer getOutputDataLength() {
        int length = 0;
        for (ProfinetVirtualSubmoduleItem module : this.virtualSubmoduleList) {
            if (module.getIoData() != null && module.getIoData().getOutput() != null) {
                for (ProfinetIoDataOutput outputIoData : module.getIoData().getOutput()) {
                    for (ProfinetDataItem dataItem : outputIoData.getDataItemList()) {
                        length += ProfinetModuleItem.getLengthSimpleType(dataItem.getDataType());
                    }
                }
            }
        }
        return length;
    }

    static ProfinetModuleItem fromXml(Element element) throws IOException {
        ProfinetModuleItem result = new ProfinetModuleItem();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "ID" -> result.id = property.string();
                case "ModuleIdentNumber" -> result.moduleIdentNumber = property.string();
                case "ModuleInfo" -> result.moduleInfo = property.object(ProfinetModuleInfo::fromXml);
                case "VirtualSubmoduleList" -> result.virtualSubmoduleList = property.wrappedList(ProfinetVirtualSubmoduleItem::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}
